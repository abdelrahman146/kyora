package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/email"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	stripelib "github.com/stripe/stripe-go/v83"
)

var syncPlansCmd = &cobra.Command{
	Use:   "sync-plans",
	Short: "Sync billing plans to database and Stripe",
	Long: `Manually sync all defined billing plans to the database and Stripe.
This command is useful for:
- Initial setup of billing plans
- Updating plan definitions after code changes
- Recovering from plan sync failures
- Production deployment updates`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Initialize Stripe
		stripeAPIKey := viper.GetString(config.StripeAPIKey)
		if stripeAPIKey == "" {
			slog.Error("Stripe API key not configured")
			fmt.Println("❌ Stripe API key not configured. Set", config.StripeAPIKey, "in config or environment")
			return
		}
		stripelib.Key = stripeAPIKey
		stripelib.SetAppInfo(&stripelib.AppInfo{Name: "Kyora", Version: "1.0", URL: "https://github.com/abdelrahman146/kyora"})

		// Initialize database connection
		dsn := viper.GetString(config.DatabaseDSN)
		logLevel := viper.GetString(config.DatabaseLogLevel)
		db := database.NewConnection(dsn, logLevel)
		defer db.CloseConnection()

		// Initialize cache
		servers := viper.GetStringSlice(config.CacheHosts)
		cacheDB := cache.NewConnection(servers)

		// Initialize dependencies for billing service
		atomicProcessor := database.NewAtomicProcess(db)
		eventBus := bus.New()
		emailClient, err := email.New()
		if err != nil {
			slog.Error("Failed to initialize email client", "error", err)
			fmt.Println("❌ Failed to initialize email client:", err)
			return
		}

		// Create account service (required by billing service)
		accountStorage := account.NewStorage(db, cacheDB)
		accountSvc := account.NewService(accountStorage, atomicProcessor, eventBus, emailClient)

		// Create billing service
		billingStorage := billing.NewStorage(db, cacheDB)
		billingSvc := billing.NewService(billingStorage, atomicProcessor, eventBus, accountSvc, emailClient)

		// Sync plans to both database and Stripe
		slog.Info("Starting complete plan sync (database + Stripe)...")
		fmt.Println("🔄 Syncing plans to database and Stripe...")

		if err := billingSvc.SyncPlansComplete(ctx); err != nil {
			slog.Error("Failed to sync plans", "error", err)
			fmt.Println("❌ Plan sync failed:", err)
			return
		}

		fmt.Println("✅ Plans synced successfully to database and Stripe")
	},
}

func init() {
	rootCmd.AddCommand(syncPlansCmd)
}
