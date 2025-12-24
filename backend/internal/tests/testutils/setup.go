package testutils

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	PostgresImage   = "postgres:16-alpine"
	MemcachedImage  = "memcached:alpine"
	StripeMockImage = "stripe/stripe-mock:latest"
)

// CreateDatabase is a context-based variant that returns the database,
// its DSN and a cleanup function without requiring *testing.T. Used in TestMain
// and for any setup that must occur before tests start.
func CreateDatabase(ctx context.Context) (*database.Database, string, func(), error) {
	dbName := "test_db"
	dbUser := "test_user"
	dbPassword := "password"
	container, err := postgres.Run(ctx, PostgresImage,
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Second)),
	)
	if err != nil {
		return nil, "", nil, fmt.Errorf("postgres container start failed: %w", err)
	}
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, "", nil, fmt.Errorf("postgres dsn retrieval failed: %w", err)
	}
	logLevel := "silent"
	db, err := database.NewConnection(dsn, logLevel)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, "", nil, fmt.Errorf("database connection failed: %w", err)
	}
	cleanup := func() {
		_ = db.CloseConnection() // errors logged via slog internally if any
		_ = container.Terminate(ctx)
	}
	return db, dsn, cleanup, nil
}

// CreateCache returns a cache client, its host:port and cleanup without *testing.T.
func CreateCache(ctx context.Context) (*cache.Cache, string, func(), error) {
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        MemcachedImage,
			ExposedPorts: []string{"11211/tcp"},
			WaitingFor:   wait.ForListeningPort("11211/tcp").WithStartupTimeout(10 * time.Second),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, "", nil, fmt.Errorf("memcached container start failed: %w", err)
	}
	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, "", nil, fmt.Errorf("memcached host retrieval failed: %w", err)
	}
	mappedPort, err := container.MappedPort(ctx, "11211/tcp")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, "", nil, fmt.Errorf("memcached port retrieval failed: %w", err)
	}
	addr := fmt.Sprintf("%s:%s", host, mappedPort.Port())
	cacheDB := cache.NewConnection([]string{addr})
	cleanup := func() { _ = container.Terminate(ctx) }
	return cacheDB, addr, cleanup, nil
}

// CreateStripeMock context-based variant without *testing.T.
func CreateStripeMock(ctx context.Context) (string, func(), error) {
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        StripeMockImage,
			ExposedPorts: []string{"12111/tcp"},
			WaitingFor:   wait.ForListeningPort("12111/tcp").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return "", nil, fmt.Errorf("stripe-mock container start failed: %w", err)
	}
	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return "", nil, fmt.Errorf("stripe-mock host retrieval failed: %w", err)
	}
	mappedPort, err := container.MappedPort(ctx, "12111/tcp")
	if err != nil {
		_ = container.Terminate(ctx)
		return "", nil, fmt.Errorf("stripe-mock port retrieval failed: %w", err)
	}
	baseURL := fmt.Sprintf("http://%s:%s", host, mappedPort.Port())

	// Extra readiness: stripe-mock may have an open port slightly before it's ready to serve.
	// Consider the service ready once it responds to any HTTP request.
	client := &http.Client{Timeout: 750 * time.Millisecond}
	deadline := time.Now().Add(10 * time.Second)
	for {
		if time.Now().After(deadline) {
			_ = container.Terminate(ctx)
			return "", nil, fmt.Errorf("stripe-mock did not become ready in time (baseURL=%s)", baseURL)
		}
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/", nil)
		resp, reqErr := client.Do(req)
		if reqErr == nil {
			_ = resp.Body.Close()
			break
		}
		time.Sleep(150 * time.Millisecond)
	}

	cleanup := func() { _ = container.Terminate(ctx) }
	return baseURL, cleanup, nil
}

// Environment aggregates shared test resources for e2e/integration tests.
type Environment struct {
	Database       *database.Database
	DatabaseDSN    string
	Cache          *cache.Cache
	CacheAddr      string
	StripeMockBase string
}

// InitEnvironment spins up all required containers. Use from TestMain.
func InitEnvironment(ctx context.Context) (*Environment, func(), error) {
	db, dsn, dbCleanup, err := CreateDatabase(ctx)
	if err != nil {
		return nil, nil, err
	}
	cacheDB, cacheAddr, cacheCleanup, err := CreateCache(ctx)
	if err != nil {
		dbCleanup()
		return nil, nil, err
	}
	stripeURL, stripeCleanup, err := CreateStripeMock(ctx)
	if err != nil {
		dbCleanup()
		cacheCleanup()
		return nil, nil, err
	}
	env := &Environment{Database: db, DatabaseDSN: dsn, Cache: cacheDB, CacheAddr: cacheAddr, StripeMockBase: stripeURL}
	cleanup := func() {
		stripeCleanup()
		cacheCleanup()
		dbCleanup()
	}
	return env, cleanup, nil
}
