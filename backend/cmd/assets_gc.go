package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/asset"
	"github.com/abdelrahman146/kyora/internal/platform/blob"
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	assetsGCDryRun       bool
	assetsGCPendingLimit int
	assetsGCOrphanLimit  int
	assetsGCOrphanMinAge time.Duration
)

// assetsGCCmd runs a best-effort garbage collection job to delete abandoned uploads.
var assetsGCCmd = &cobra.Command{
	Use:   "assets-gc",
	Short: "Garbage-collect expired and orphaned uploaded assets",
	RunE: func(cmd *cobra.Command, args []string) error {
		dsn := viper.GetString(config.DatabaseDSN)
		logLevel := viper.GetString(config.DatabaseLogLevel)
		db, err := database.NewConnection(dsn, logLevel)
		if err != nil {
			return err
		}

		servers := viper.GetStringSlice(config.CacheHosts)
		cacheDB := cache.NewConnection(servers)

		provider, err := blob.FromConfig()
		if err != nil {
			return err
		}

		st := asset.NewStorage(db, cacheDB)
		svc := asset.NewService(st, database.NewAtomicProcess(db), provider)

		res, err := svc.GarbageCollect(context.Background(), asset.GarbageCollectOptions{
			PendingLimit: assetsGCPendingLimit,
			OrphanLimit:  assetsGCOrphanLimit,
			OrphanMinAge: assetsGCOrphanMinAge,
			DryRun:       assetsGCDryRun,
		})
		if err != nil {
			slog.Error("assets gc failed", "error", err)
			return err
		}

		slog.Info("assets gc completed",
			"expiredPendingCandidates", res.ExpiredPendingCandidates,
			"readyOrphanCandidates", res.ReadyOrphanCandidates,
			"deletedAssets", res.DeletedAssets,
			"deletedBlobs", res.DeletedBlobs,
			"deletedFiles", res.DeletedFiles,
			"errors", res.Errors,
		)

		if res.Errors > 0 {
			return fmt.Errorf("assets gc completed with %d errors", res.Errors)
		}
		return nil
	},
}

func init() {
	assetsGCCmd.Flags().BoolVar(&assetsGCDryRun, "dry-run", false, "Log actions without deleting")
	assetsGCCmd.Flags().IntVar(&assetsGCPendingLimit, "pending-limit", 500, "Max expired pending uploads deleted per run")
	assetsGCCmd.Flags().IntVar(&assetsGCOrphanLimit, "orphan-limit", 500, "Max orphan ready assets deleted per run")
	assetsGCCmd.Flags().DurationVar(&assetsGCOrphanMinAge, "orphan-min-age", 30*time.Minute, "Min age for ready assets before orphan cleanup")

	rootCmd.AddCommand(assetsGCCmd)
}
