package account

import "github.com/abdelrahman146/kyora/internal/db"

type AccountDomain struct {
	AuthService       *AuthenticationService
	OrgService        *OrganizationService
	UserService       *UserService
	OnboardingService *OnboardingService
}

func NewDomain(postgres *db.Postgres, atomicProcess *db.AtomicProcess, cache *db.Memcache) *AccountDomain {
	userRepo := NewUserRepository(postgres)
	organizationRepo := NewOrganizationRepository(postgres)
	postgres.AutoMigrate(&User{}, &Organization{})

	return &AccountDomain{
		AuthService:       NewAuthenticationService(userRepo, cache),
		OrgService:        NewOrganizationService(organizationRepo),
		UserService:       NewUserService(userRepo),
		OnboardingService: NewOnboardingService(userRepo, organizationRepo, atomicProcess),
	}
}
