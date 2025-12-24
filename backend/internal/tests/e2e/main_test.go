package e2e_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/server"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/spf13/viper"
)

var (
	testEnv    *testutils.Environment
	testServer *server.Server
)

const e2eBaseURL = "http://localhost:18080"

func TestMain(m *testing.M) {
	fmt.Println("Setting up e2e test environment...")
	ctx := context.Background()

	// Override email provider to use mock for tests
	viper.Set(config.EmailProvider, "mock")
	// JWT is required for authenticated endpoints in tests.
	viper.Set(config.JWTSecret, "test_jwt_secret")

	// Keep tests quiet and faster; E2E suites can generate lots of DB queries.
	// (Set to "info" locally when debugging query issues.)
	viper.Set(config.DatabaseLogLevel, "silent")

	// Configure Stripe for stripe-mock.
	// stripe-mock accepts a limited set of test keys; use the canonical one.
	viper.Set(config.StripeAPIKey, "sk_test_123")
	// Webhook secret isn't required for most E2E flows, but keep it non-empty.
	viper.Set(config.StripeWebhookSecret, "whsec_test")

	// Disable automatic plan sync for test isolation
	// Tests will create their own plans as needed
	viper.Set(config.BillingAutoSyncPlans, false)

	// Auto-migrate is very expensive when it runs once per repository construction.
	// In E2E, the server boot will touch all storages once (migrating tables).
	// After that, tests and helpers may create additional storages for fixtures,
	// so we disable auto-migrate to keep `go test ./...` bounded.
	viper.Set(config.DatabaseAutoMigrate, true)

	env, cleanup, err := testutils.InitEnvironment(ctx)
	if err != nil {
		log.Fatalf("environment init failed: %v", err)
	}
	testEnv = env

	// Create server with container-provided dependencies (DB, cache, stripe-mock)
	testServer, err = server.New(
		server.WithDatabaseDSN(env.DatabaseDSN),
		server.WithCacheHosts([]string{env.CacheAddr}),
		server.WithStripeBaseURL(env.StripeMockBase),
		server.WithServerAddress(":18080"), // isolate test port
	)
	if err != nil {
		cleanup()
		log.Fatalf("failed to create server: %v", err)
	}
	if err := testServer.Start(); err != nil {
		cleanup()
		log.Fatalf("failed to start server: %v", err)
	}

	// Prevent repeated migrations from test helpers/suites.
	viper.Set(config.DatabaseAutoMigrate, false)

	// Run tests
	exitCode := m.Run()

	fmt.Println("Tearing down e2e test environment...")
	// Graceful server stop then cleanup containers
	if testServer != nil {
		if err := testServer.Stop(); err != nil {
			log.Printf("server stop error: %v", err)
		}
	}
	cleanup()

	os.Exit(exitCode)
}
