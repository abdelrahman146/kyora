package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd starts the HTTP server
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Kyora HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, err := server.New()
		if err != nil {
			return err
		}

		// Auto-sync billing plans if enabled (default: true in production, false in tests)
		viper.SetDefault(config.BillingAutoSyncPlans, true)
		if viper.GetBool(config.BillingAutoSyncPlans) {
			slog.Info("Starting billing plan auto-sync in background")
			go func() {
				ctx := context.Background()
				if err := srv.SyncPlansComplete(ctx); err != nil {
					slog.Error("Failed to auto-sync plans", "error", err)
				} else {
					slog.Info("Billing plans auto-synced successfully")
				}
			}()
		} else {
			slog.Info("Billing plan auto-sync disabled")
		}

		// Start the server (non-blocking); Stop() will gracefully drain
		if err := srv.Start(); err != nil {
			return err
		}

		slog.Info("HTTP server started; waiting for shutdown signal")

		// Wait for SIGINT/SIGTERM
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop

		slog.Info("Shutdown signal received; stopping server...")
		if err := srv.Stop(); err != nil {
			return err
		}
		slog.Info("Server stopped gracefully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
