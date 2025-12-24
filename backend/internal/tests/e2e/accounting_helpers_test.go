package e2e_test

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/hash"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
)

// AccountingTestHelper provides reusable helpers for accounting E2E tests.
type AccountingTestHelper struct {
	db             *database.Database
	Client         *testutils.HTTPClient
	accountStorage *account.Storage
}

func NewAccountingTestHelper(db *database.Database, cacheAddr, baseURL string) *AccountingTestHelper {
	cacheClient := cache.NewConnection([]string{cacheAddr})
	acctStorage := account.NewStorage(db, cacheClient)

	return &AccountingTestHelper{
		db:             db,
		Client:         testutils.NewHTTPClient(baseURL),
		accountStorage: acctStorage,
	}
}

type WorkspaceUsers struct {
	Workspace   *account.Workspace
	Admin       *account.User
	Member      *account.User
	AdminToken  string
	MemberToken string
	Business    *business.Business
}

func (h *AccountingTestHelper) CreateWorkspaceWithAdminAndMemberAndBusiness(ctx context.Context) (*WorkspaceUsers, error) {
	admin, ws, adminToken, err := testutils.CreateAuthenticatedUser(ctx, h.db, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	if err != nil {
		return nil, err
	}

	memberPassword, err := hash.Password("Password123!")
	if err != nil {
		return nil, err
	}
	member := &account.User{
		WorkspaceID:     ws.ID,
		Role:            role.RoleUser,
		FirstName:       "Member",
		LastName:        "User",
		Email:           "member@example.com",
		Password:        memberPassword,
		IsEmailVerified: true,
	}
	userRepo := database.NewRepository[account.User](h.db)
	if err := userRepo.CreateOne(ctx, member); err != nil {
		return nil, err
	}

	memberToken, err := auth.NewJwtToken(member.ID, member.WorkspaceID, member.AuthVersion)
	if err != nil {
		return nil, err
	}

	biz, err := h.CreateBusiness(ctx, ws.ID)
	if err != nil {
		return nil, err
	}

	return &WorkspaceUsers{
		Workspace:   ws,
		Admin:       admin,
		Member:      member,
		AdminToken:  adminToken,
		MemberToken: memberToken,
		Business:    biz,
	}, nil
}

func (h *AccountingTestHelper) CreateBusiness(ctx context.Context, workspaceID string) (*business.Business, error) {
	bizRepo := database.NewRepository[business.Business](h.db)
	biz := &business.Business{
		WorkspaceID:   workspaceID,
		Descriptor:    "default",
		Name:          "Test Business",
		CountryCode:   "AE",
		Currency:      "aed",
		EstablishedAt: time.Now().UTC(),
	}
	if err := bizRepo.CreateOne(ctx, biz); err != nil {
		return nil, err
	}
	return biz, nil
}
