package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/accounting"
	"github.com/abdelrahman146/kyora/internal/domain/analytics"
	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/domain/onboarding"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/email"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v83"
)

type Server struct {
	db         *database.Database
	cacheDB    *cache.Cache
	r          *gin.Engine
	httpSrv    *http.Server
	billingSvc *billing.Service
}

type ServerConfig struct {
	Address       string
	DSN           string
	LogLevel      string
	DBLogLevel    string
	CacheHosts    []string
	StripeKey     string
	StripeBaseURL string
}

func WithDatabaseDSN(dsn string) func(*ServerConfig) {
	return func(cfg *ServerConfig) {
		cfg.DSN = dsn
	}
}

func WithServerAddress(addr string) func(*ServerConfig) {
	return func(cfg *ServerConfig) {
		cfg.Address = addr
	}
}

func WithCacheHosts(hosts []string) func(*ServerConfig) {
	return func(cfg *ServerConfig) {
		cfg.CacheHosts = hosts
	}
}

func WithStripeKey(key string) func(*ServerConfig) {
	return func(cfg *ServerConfig) {
		cfg.StripeKey = key
	}
}

func WithStripeBaseURL(url string) func(*ServerConfig) {
	return func(cfg *ServerConfig) {
		cfg.StripeBaseURL = url
	}
}

func New(opts ...func(*ServerConfig)) (*Server, error) {
	// apply options and set config values and override env variables
	conf := &ServerConfig{}
	for _, opt := range opts {
		opt(conf)
	}
	// apply address override if provided (port only or host:port)
	if conf.Address != "" {
		// Accept forms ":8081" or "8081"; strip host if present for port key simplicity
		addr := conf.Address
		// normalize to just port without leading host/protocol
		if len(addr) > 0 && addr[0] == ':' { // already :port
			addr = addr[1:]
		}
		viper.Set(config.HTTPPort, addr)
	}
	if conf.DSN != "" {
		viper.Set(config.DatabaseDSN, conf.DSN)
	}
	if conf.CacheHosts != nil {
		viper.Set(config.CacheHosts, conf.CacheHosts)
	}
	if conf.LogLevel != "" {
		viper.Set(config.LogLevel, conf.LogLevel)
	}
	if conf.DBLogLevel != "" {
		viper.Set(config.DatabaseLogLevel, conf.DBLogLevel)
	}
	if conf.StripeKey != "" {
		viper.Set(config.StripeAPIKey, conf.StripeKey)
	}

	// initialize stripe client
	stripeAPIKey := viper.GetString(config.StripeAPIKey)
	stripe.Key = stripeAPIKey
	stripe.SetAppInfo(&stripe.AppInfo{Name: "Kyora", Version: "1.0", URL: "https://github.com/abdelrahman146/kyora"})
	if conf.StripeBaseURL != "" {
		backend := stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{
			URL: &conf.StripeBaseURL, // optional custom base URL for testing if not provided it will automatically use the default
		})
		stripe.SetBackend(stripe.APIBackend, backend)
	}
	slog.Info("Stripe client initialized")

	// initialize database and cache connections
	dsn := viper.GetString(config.DatabaseDSN)
	logLevel := viper.GetString(config.DatabaseLogLevel)

	db := database.NewConnection(dsn, logLevel)
	servers := viper.GetStringSlice(config.CacheHosts)
	cacheDB := cache.NewConnection(servers)
	atomicProcessor := database.NewAtomicProcess(db)
	bus := bus.New()
	emailClient, err := email.New()
	if err != nil {
		return nil, err
	}

	// DI - create storages first
	accountStorage := account.NewStorage(db, cacheDB)
	billingStorage := billing.NewStorage(db, cacheDB)

	// Create services with email integrations
	accountSvc := account.NewService(accountStorage, atomicProcessor, bus, emailClient)

	billingSvc := billing.NewService(billingStorage, atomicProcessor, bus, accountSvc, emailClient)

	// Note: Plan auto-sync is now handled in the server command (cmd/server.go)
	// This keeps server initialization clean and allows sync to run asynchronously

	businessStorage := business.NewStorage(db, cacheDB)
	businessSvc := business.NewService(businessStorage, atomicProcessor, bus)

	inventoryStorage := inventory.NewStorage(db, cacheDB)
	inventorySvc := inventory.NewService(inventoryStorage, atomicProcessor, bus)

	accountingStorage := accounting.NewStorage(db, cacheDB)
	accountingSvc := accounting.NewService(accountingStorage, atomicProcessor, bus)

	customerStorage := customer.NewStorage(db, cacheDB)
	customerSvc := customer.NewService(customerStorage, atomicProcessor, bus)

	orderStorage := order.NewStorage(db, cacheDB)
	orderSvc := order.NewService(orderStorage, atomicProcessor, bus, inventorySvc)

	analyticsSvc := analytics.NewService(&analytics.ServiceParams{
		Inventory:  inventorySvc,
		Orders:     orderSvc,
		Accounting: accountingSvc,
		Customer:   customerSvc,
	})

	// onboarding routes
	onboardingStorage := onboarding.NewStorage(db, cacheDB)
	onboardingSvc := onboarding.NewService(onboardingStorage, atomicProcessor, accountSvc, billingSvc, businessSvc, emailClient)

	// server initialization logic
	r := gin.New()
	r.Use(logger.Middleware())
	r.Use(gin.Recovery())

	// health endpoint
	r.GET("/healthz", func(c *gin.Context) { response.SuccessText(c, 200, "ok") })
	r.GET("/livez", func(c *gin.Context) { response.SuccessText(c, 200, "ok") })

	// register domain routes under /api
	registerBillingRoutes(r, billing.NewHttpHandler(billingSvc, accountSvc), accountSvc)

	// Register account routes with plan limit enforcement for team members
	registerAccountRoutes(r, account.NewHttpHandler(accountSvc), accountSvc, billingSvc)

	// Register onboarding routes
	registerOnboardingRoutes(r, onboarding.NewHttpHandler(onboardingSvc))

	// Register accounting routes
	registerAccountingRoutes(r, accounting.NewHttpHandler(accountingSvc, businessSvc, orderSvc), accountSvc)

	// Register analytics routes
	registerAnalyticsRoutes(r, analytics.NewHttpHandler(analyticsSvc, businessSvc), accountSvc)

	// Register business routes
	registerBusinessRoutes(r, business.NewHttpHandler(businessSvc), accountSvc, billingSvc, businessSvc)

	// Register customer routes
	registerCustomerRoutes(r, customer.NewHttpHandler(customerSvc, businessSvc), accountSvc)

	return &Server{r: r, db: db, cacheDB: cacheDB, billingSvc: billingSvc}, nil
}

func (s *Server) Start() error {
	// resolve address
	addr := viper.GetString(config.HTTPPort)
	if addr == "" {
		addr = ":8080"
	} else if addr[0] != ':' {
		addr = ":" + addr
	}

	baseURL := viper.GetString(config.HTTPBaseURL)
	if baseURL == "" {
		baseURL = "http://localhost" + addr
	}

	// create http server to enable graceful shutdown
	s.httpSrv = &http.Server{
		Addr:         addr,
		Handler:      s.r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Bind synchronously so callers can reliably detect startup failures
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// start server in background; Stop() will gracefully shut it down
	go func() {
		if err := s.httpSrv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error", "error", err)
		}
	}()

	slog.Info(fmt.Sprintf("Server started successfully at %s", baseURL))
	return nil
}

func (s *Server) Stop() error {
	var retErr error

	// gracefully stop HTTP server, waiting for in-flight requests
	if s.httpSrv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.httpSrv.Shutdown(ctx); err != nil {
			slog.Error("HTTP server shutdown error", "error", err)
			retErr = errors.Join(retErr, err)
		} else {
			slog.Info("HTTP server shut down gracefully")
		}
	}

	// close database connection
	if s.db != nil {
		if err := s.db.CloseConnection(); err != nil {
			slog.Error("Database close error", "error", err)
			retErr = errors.Join(retErr, err)
		} else {
			slog.Info("Database connection closed")
		}
	}

	// cache client (gomemcache) doesn't require explicit close; log for visibility
	if s.cacheDB != nil {
		slog.Info("Cache client ready for shutdown (no close required)")
	}

	return retErr
}

// SyncPlansComplete syncs billing plans to both database and Stripe
// This method is exposed for use in the server command's auto-sync goroutine
func (s *Server) SyncPlansComplete(ctx context.Context) error {
	if s.billingSvc == nil {
		return fmt.Errorf("billing service not initialized")
	}
	return s.billingSvc.SyncPlansComplete(ctx)
}
