package onboarding

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"gorm.io/gorm"

	// cross-domain models

	acc "github.com/abdelrahman146/kyora/internal/domain/account"

	biz "github.com/abdelrahman146/kyora/internal/domain/business"
)

type Storage struct {
	session *database.Repository[OnboardingSession]
	cache   *cache.Cache
	// cross-domain repos used only at final commit
	workspace *database.Repository[acc.Workspace]
	user      *database.Repository[acc.User]
	business  *database.Repository[biz.Business]
}

func NewStorage(db *database.Database, c *cache.Cache) *Storage {
	return &Storage{
		session:   database.NewRepository[OnboardingSession](db),
		cache:     c,
		workspace: database.NewRepository[acc.Workspace](db),
		user:      database.NewRepository[acc.User](db),
		business:  database.NewRepository[biz.Business](db),
	}
}

func (s *Storage) CreateSession(ctx context.Context, sess *OnboardingSession) error {
	return s.session.CreateOne(ctx, sess)
}

func (s *Storage) UpdateSession(ctx context.Context, sess *OnboardingSession) error {
	return s.session.UpdateOne(ctx, sess)
}

func (s *Storage) GetByToken(ctx context.Context, token string) (*OnboardingSession, error) {
	return s.session.FindOne(ctx, s.session.ScopeEquals(SessionSchema.Token, token))
}

func (s *Storage) GetActiveByEmail(ctx context.Context, email string) (*OnboardingSession, error) {
	return s.session.FindOne(ctx,
		s.session.ScopeEquals(SessionSchema.Email, email),
		s.session.ScopeNotEquals(SessionSchema.Stage, StageCommitted),
		func(db *gorm.DB) *gorm.DB { return db.Where("expires_at > ? AND (committed_at IS NULL)", time.Now()) },
	)
}

func (s *Storage) DeleteSession(ctx context.Context, sess *OnboardingSession) error {
	return s.session.DeleteOne(ctx, sess)
}

// Finalization helpers
func (s *Storage) CreateWorkspace(ctx context.Context, ws *acc.Workspace, opts ...func(db *gorm.DB) *gorm.DB) error {
	return s.workspace.CreateOne(ctx, ws, opts...)
}
func (s *Storage) UpdateWorkspace(ctx context.Context, ws *acc.Workspace, opts ...func(db *gorm.DB) *gorm.DB) error {
	return s.workspace.UpdateOne(ctx, ws, opts...)
}
func (s *Storage) CreateUser(ctx context.Context, u *acc.User, opts ...func(db *gorm.DB) *gorm.DB) error {
	return s.user.CreateOne(ctx, u, opts...)
}
func (s *Storage) CreateBusiness(ctx context.Context, b *biz.Business, opts ...func(db *gorm.DB) *gorm.DB) error {
	return s.business.CreateOne(ctx, b, opts...)
}

// Cleanup helpers
func (s *Storage) DeleteExpiredSessionsByEmail(ctx context.Context, email string) error {
	return s.session.DeleteMany(ctx, func(db *gorm.DB) *gorm.DB { return db.Where("email = ? AND expires_at <= now()", email) })
}

// DeleteAllExpired removes all sessions past expiry or already committed more than 24h ago
func (s *Storage) DeleteAllExpired(ctx context.Context) error {
	return s.session.DeleteMany(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("expires_at <= now() OR (committed_at IS NOT NULL AND committed_at <= now() - interval '24 hours')")
	})
}
