package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/accounting"
	"github.com/abdelrahman146/kyora/internal/domain/analytics"
	"github.com/abdelrahman146/kyora/internal/domain/asset"
	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/domain/onboarding"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/domain/storefront"
	"github.com/abdelrahman146/kyora/internal/platform/blob"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/email"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/abdelrahman146/kyora/internal/platform/middleware"
	"github.com/abdelrahman146/kyora/internal/platform/request"
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
	// Ensure config defaults are present even when running outside Cobra (e.g., tests).
	config.Configure()

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
	stripeBaseURL := conf.StripeBaseURL
	if stripeBaseURL == "" {
		stripeBaseURL = viper.GetString(config.StripeAPIBaseURL)
	}
	if stripeBaseURL != "" {
		backend := stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{
			URL: &stripeBaseURL, // optional custom base URL for testing if not provided it will automatically use the default
		})
		stripe.SetBackend(stripe.APIBackend, backend)
		slog.Info("Stripe client initialized", "baseURL", stripeBaseURL)
	} else {
		slog.Info("Stripe client initialized")
	}

	// initialize database and cache connections
	dsn := viper.GetString(config.DatabaseDSN)
	logLevel := viper.GetString(config.DatabaseLogLevel)

	db, err := database.NewConnection(dsn, logLevel)
	if err != nil {
		return nil, err
	}
	servers := viper.GetStringSlice(config.CacheHosts)
	cacheDB := cache.NewConnection(servers)
	atomicProcessor := database.NewAtomicProcess(db)
	bus := bus.New()
	emailClient, err := email.New()
	if err != nil {
		return nil, err
	}

	// asset/blob storage (uploads)
	blobProvider, err := blob.FromConfig()
	if err != nil {
		return nil, err
	}
	assetStorage := asset.NewStorage(db, cacheDB)
	assetSvc := asset.NewService(assetStorage, atomicProcessor, blobProvider)

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
	accounting.NewBusHandler(bus, accountingSvc, businessSvc)

	customerStorage := customer.NewStorage(db, cacheDB)
	customerSvc := customer.NewService(customerStorage, atomicProcessor, bus)

	orderStorage := order.NewStorage(db, cacheDB)
	orderSvc := order.NewService(orderStorage, atomicProcessor, bus, inventorySvc, customerSvc, businessSvc)

	storefrontStorage := storefront.NewStorage(db, cacheDB)
	storefrontSvc := storefront.NewService(storefrontStorage, atomicProcessor, businessSvc, inventorySvc, customerSvc, orderSvc)

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
	r.Use(request.LimitBodySize(viper.GetInt64(config.HTTPMaxBodyBytes)))
	r.Use(gin.Recovery())

	// health endpoint
	r.GET("/healthz", func(c *gin.Context) { response.SuccessText(c, 200, "ok") })
	r.GET("/livez", func(c *gin.Context) { response.SuccessText(c, 200, "ok") })

	// CORS preflight handler
	// Browsers send OPTIONS requests (preflight) for non-simple cross-origin requests.
	// Gin does not automatically provide OPTIONS routes for POST endpoints, so we
	// add a catch-all that applies the right CORS policy and returns 204.
	r.OPTIONS("/*path", func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/v1/storefront") {
			middleware.NewPublicCORSMiddleware()(c)
		} else {
			middleware.NewCORSMiddleware()(c)
		}
	})

	// register domain routes under /api
	registerBillingRoutes(r, billing.NewHttpHandler(billingSvc, accountSvc), accountSvc)

	// Public storefront routes (no auth required)
	registerStorefrontRoutes(r, storefront.NewHttpHandler(storefrontSvc))

	// Register account routes with plan limit enforcement for team members
	registerAccountRoutes(r, account.NewHttpHandler(accountSvc), accountSvc, billingSvc)

	// Register onboarding routes
	registerOnboardingRoutes(r, onboarding.NewHttpHandler(onboardingSvc))

	accountingHandler := accounting.NewHttpHandler(accountingSvc, orderSvc)
	analyticsHandler := analytics.NewHttpHandler(analyticsSvc)
	customerHandler := customer.NewHttpHandler(customerSvc)
	inventoryHandler := inventory.NewHttpHandler(inventorySvc)
	orderHandler := order.NewHttpHandler(orderSvc)
	businessHandler := business.NewHttpHandler(businessSvc)
	assetHandler := asset.NewHttpHandler(assetSvc)

	// Public asset serving routes (no auth required)
	registerPublicAssetRoutes(r, assetHandler)

	// Register business-scoped routes
	registerBusinessScopedRoutes(r, accountSvc, billingSvc, businessSvc, businessHandler, assetHandler, accountingHandler, analyticsHandler, customerHandler, inventoryHandler, orderHandler)

	// Register business routes
	registerBusinessRoutes(r, businessHandler, accountSvc, billingSvc, businessSvc)

	// Customer, inventory, analytics, and accounting routes are registered under business-scoped routes.

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
