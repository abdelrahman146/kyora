package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/accounting"
	"github.com/abdelrahman146/kyora/internal/domain/analytics"
	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v83"
)

type Server struct {
	db      *database.Database
	cacheDB *cache.Cache
	r       *gin.Engine
	httpSrv *http.Server
}

func New() (*Server, error) {
	// Initialize Stripe with API key
	stripeAPIKey := viper.GetString(config.StripeAPIKey)
	if stripeAPIKey != "" {
		stripe.Key = stripeAPIKey
		slog.Info("Stripe client initialized")
	} else {
		slog.Warn("Stripe API key not configured - billing functionality will be limited")
	}

	db := database.NewConnection()
	cacheDB := cache.NewConnection()
	atomicProcessor := database.NewAtomicProcess(db)
	bus := bus.New()

	// Email service initialization (add proper email client initialization here in future)
	// For now, we'll create a nil email integration to avoid breaking the compilation

	// DI - create storages first
	accountStorage := account.NewStorage(db, cacheDB)
	billingStorage := billing.NewStorage(db, cacheDB)

	// Create placeholder email integrations (will be properly configured with email client later)
	accountEmailIntegration := (*account.EmailIntegration)(nil)
	billingEmailIntegration := (*billing.EmailIntegration)(nil)

	// Create services with email integrations
	accountSvc := account.NewService(accountStorage, atomicProcessor, bus, accountEmailIntegration)
	billingSvc := billing.NewService(billingStorage, atomicProcessor, bus, accountSvc, billingEmailIntegration)
	_ = billingSvc // to avoid unused variable warning

	businessStorage := business.NewStorage(db, cacheDB)
	businessSvc := business.NewService(businessStorage, atomicProcessor, bus)
	_ = businessSvc // to avoid unused variable warning

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
	_ = analyticsSvc // to avoid unused variable warning

	// server initialization logic
	r := gin.New()
	r.Use(logger.Middleware())

	// health endpoint
	r.GET("/healthz", func(c *gin.Context) { response.SuccessText(c, 200, "ok") })
	r.GET("/livez", func(c *gin.Context) { response.SuccessText(c, 200, "ok") })

	// register system routes
	api := r.Group("/api")
	_ = api // to avoid unused variable warning

	return &Server{r: r, db: db, cacheDB: cacheDB}, nil
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

	// start server in background; Stop() will gracefully shut it down
	go func() {
		if err := s.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
