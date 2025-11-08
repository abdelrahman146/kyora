package cmd

import (
	"context"
	"log/slog"

	"github.com/abdelrahman146/kyora/internal/domain/onboarding"
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/spf13/cobra"
)

// onboardingCleanupCmd removes expired/committed onboarding sessions
var onboardingCleanupCmd = &cobra.Command{
	Use:   "onboarding-cleanup",
	Short: "Delete expired and finalized onboarding sessions",
	RunE: func(cmd *cobra.Command, args []string) error {
		db := database.NewConnection()
		cacheDB := cache.NewConnection()
		storage := onboarding.NewStorage(db, cacheDB)
		if err := storage.DeleteAllExpired(context.Background()); err != nil {
			slog.Error("onboarding cleanup failed", "error", err)
			return err
		}
		slog.Info("onboarding cleanup completed successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(onboardingCleanupCmd)
}
