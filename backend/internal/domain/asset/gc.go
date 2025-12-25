package asset

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/logger"
)

// GarbageCollectOptions configures an asset garbage collection run.
//
// This is designed for a scheduled maintenance job (cron/k8s) and should be safe to run repeatedly.
type GarbageCollectOptions struct {
	Now time.Time

	// PendingLimit bounds the number of expired pending uploads deleted per run.
	PendingLimit int
	// OrphanLimit bounds the number of ready orphans deleted per run.
	OrphanLimit int
	// OrphanMinAge ensures we only consider ready assets older than this duration.
	OrphanMinAge time.Duration

	DryRun bool
}

// GarbageCollectResult summarizes a GC run.
type GarbageCollectResult struct {
	ExpiredPendingCandidates int
	ReadyOrphanCandidates    int

	DeletedAssets int
	DeletedBlobs  int
	DeletedFiles  int

	Errors int
}

func (o GarbageCollectOptions) withDefaults() GarbageCollectOptions {
	out := o
	if out.Now.IsZero() {
		out.Now = time.Now().UTC()
	} else {
		out.Now = out.Now.UTC()
	}
	if out.PendingLimit == 0 {
		out.PendingLimit = 500
	}
	if out.OrphanLimit == 0 {
		out.OrphanLimit = 500
	}
	if out.OrphanMinAge == 0 {
		out.OrphanMinAge = 30 * time.Minute
	}
	return out
}

// GarbageCollect removes expired pending uploads and ready orphan assets.
//
// Orphan detection is done in the storage query using URL-only reference checks.
func (s *Service) GarbageCollect(ctx context.Context, opts GarbageCollectOptions) (*GarbageCollectResult, error) {
	opts = opts.withDefaults()
	log := logger.FromContext(ctx).With(
		slog.Time("now", opts.Now),
		slog.Bool("dryRun", opts.DryRun),
	)

	res := &GarbageCollectResult{}

	// 1) Expired pending uploads.
	pending, err := s.storage.ListExpiredPending(ctx, opts.Now, opts.PendingLimit)
	if err != nil {
		return nil, err
	}
	res.ExpiredPendingCandidates = len(pending)
	for _, a := range pending {
		if err := s.gcDeleteOne(ctx, log, a, opts.DryRun, res); err != nil {
			res.Errors++
			log.Warn("asset gc: failed to delete expired pending", slog.String("assetId", a.ID), slog.Any("error", err))
		}
	}

	// 2) Ready orphans.
	orphans, err := s.storage.ListReadyOrphans(ctx, opts.Now, opts.OrphanMinAge, opts.OrphanLimit)
	if err != nil {
		return nil, err
	}
	res.ReadyOrphanCandidates = len(orphans)
	for _, a := range orphans {
		if err := s.gcDeleteOne(ctx, log, a, opts.DryRun, res); err != nil {
			res.Errors++
			log.Warn("asset gc: failed to delete ready orphan", slog.String("assetId", a.ID), slog.Any("error", err))
		}
	}

	log.Info("asset gc complete",
		slog.Int("expiredPendingCandidates", res.ExpiredPendingCandidates),
		slog.Int("readyOrphanCandidates", res.ReadyOrphanCandidates),
		slog.Int("deletedAssets", res.DeletedAssets),
		slog.Int("deletedBlobs", res.DeletedBlobs),
		slog.Int("deletedFiles", res.DeletedFiles),
		slog.Int("errors", res.Errors),
	)

	return res, nil
}

func (s *Service) gcDeleteOne(ctx context.Context, log *slog.Logger, a *Asset, dryRun bool, res *GarbageCollectResult) error {
	if a == nil {
		return nil
	}
	if dryRun {
		log.Info("asset gc: would delete asset",
			slog.String("assetId", a.ID),
			slog.String("status", string(a.Status)),
			slog.String("publicUrl", a.PublicURL),
			slog.String("objectKey", a.ObjectKey),
			slog.String("localFilePath", a.LocalFilePath),
		)
		return nil
	}

	// Delete local file if present (ignore missing).
	if a.LocalFilePath != "" {
		if err := os.Remove(a.LocalFilePath); err == nil {
			res.DeletedFiles++
		} else if !os.IsNotExist(err) {
			return err
		}
	}

	// Delete blob object if configured.
	if s.blob != nil && a.ObjectKey != "" {
		if err := s.blob.Delete(ctx, a.ObjectKey); err != nil {
			return err
		}
		res.DeletedBlobs++
	}

	if err := s.storage.Delete(ctx, a); err != nil {
		return err
	}
	res.DeletedAssets++
	return nil
}
