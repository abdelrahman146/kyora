package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/accounting"
	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/email"
	assetTypes "github.com/abdelrahman146/kyora/internal/platform/types/asset"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/hash"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	stripelib "github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/invoice"
	"github.com/stripe/stripe-go/v83/paymentmethod"
)

type seedSize string

const (
	seedSmall  seedSize = "small"
	seedMedium seedSize = "medium"
	seedLarge  seedSize = "large"
)

type seedCounts struct {
	categories int
	products   int
	customers  int
	orders     int
}

func (s seedSize) counts() seedCounts {
	switch s {
	case seedMedium:
		return seedCounts{categories: 5, products: 12, customers: 40, orders: 90}
	case seedLarge:
		return seedCounts{categories: 10, products: 30, customers: 120, orders: 300}
	default:
		return seedCounts{categories: 3, products: 6, customers: 15, orders: 30}
	}
}

var (
	seedClean    bool
	seedSizeFlag string
	seedPassword string
)

type workspaceSeedConfig struct {
	Name              string
	Slug              string
	BusinessName      string
	Descriptor        string
	Country           string
	Currency          string
	Brand             string
	PlanDescriptor    string
	SupportEmail      string
	PhoneNumber       string
	WhatsappNumber    string
	StorefrontEnabled bool
	StorefrontTheme   business.StorefrontTheme
	Counts            seedCounts
	TeamMembers       []seedTeamMember
}

type seedTeamMember struct {
	FirstName string
	LastName  string
	Email     string
	Role      role.Role
}

type seedWorkspaceResult struct {
	Workspace       *account.Workspace
	Business        *business.Business
	Owner           *account.User
	CreatedOrders   int
	PaidOrders      int
	FulfilledOrders int
}

type seedDeps struct {
	accountSvc    *account.Service
	accountStore  *account.Storage
	billingSvc    *billing.Service
	businessSvc   *business.Service
	inventorySvc  *inventory.Service
	customerSvc   *customer.Service
	accountingSvc *accounting.Service
	orderSvc      *order.Service
	billingStore  *billing.Storage
	stripeEnabled bool
	rng           *rand.Rand
}

func isStripeSeedEnabled(stripeKey, baseURL string) bool {
	if strings.TrimSpace(baseURL) != "" {
		// Explicit stripe-mock or custom Stripe API backend; allow the seed flow.
		return true
	}
	key := strings.TrimSpace(stripeKey)
	if key == "" {
		return false
	}
	// Refuse obvious placeholders to avoid confusing 401s during local seed.
	lk := strings.ToLower(key)
	if strings.Contains(lk, "your_stripe") || strings.Contains(lk, "your_stripe_api_key") || strings.Contains(lk, "changeme") || strings.Contains(lk, "replace") || strings.Contains(lk, "********") {
		return false
	}
	// The canonical stripe-mock test key only works with a custom base URL.
	if key == "sk_test_123" {
		return false
	}
	// Secret keys must start with sk_ (pk_ is publishable).
	return strings.HasPrefix(key, "sk_")
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed local development data",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		sz := seedSize(strings.TrimSpace(strings.ToLower(seedSizeFlag)))
		if sz != seedSmall && sz != seedMedium && sz != seedLarge {
			return errors.New("invalid --size. allowed: small|medium|large")
		}

		baseCounts := sz.counts()
		stripeKey, stripeBaseURL := loadStripeConfig()
		stripeEnabled := isStripeSeedEnabled(stripeKey, stripeBaseURL)
		if stripeEnabled {
			initStripeClient(stripeKey, stripeBaseURL)
		} else {
			slog.Warn("Stripe not configured (or placeholder key); Stripe seed steps will be skipped", "baseURL", stripeBaseURL)
		}

		dsn := viper.GetString(config.DatabaseDSN)
		logLevel := viper.GetString(config.DatabaseLogLevel)
		db, err := database.NewConnection(dsn, logLevel)
		if err != nil {
			slog.Error("Failed to connect to database", "error", err)
			return fmt.Errorf("failed to connect to database: %w", err)
		}
		defer db.CloseConnection()

		servers := viper.GetStringSlice(config.CacheHosts)
		cacheDB := cache.NewConnection(servers)

		atomicProcessor := database.NewAtomicProcess(db)
		eventBus := bus.New()
		defer eventBus.Close()

		emailClient, err := email.New()
		if err != nil {
			slog.Error("Failed to initialize email client", "error", err)
			return fmt.Errorf("failed to initialize email client: %w", err)
		}

		accountStorage := account.NewStorage(db, cacheDB)
		accountSvc := account.NewService(accountStorage, atomicProcessor, eventBus, emailClient)

		billingStorage := billing.NewStorage(db, cacheDB)
		billingSvc := billing.NewService(billingStorage, atomicProcessor, eventBus, accountSvc, emailClient)

		businessStorage := business.NewStorage(db, cacheDB)
		businessSvc := business.NewService(businessStorage, atomicProcessor, eventBus)

		inventoryStorage := inventory.NewStorage(db, cacheDB)
		inventorySvc := inventory.NewService(inventoryStorage, atomicProcessor, eventBus)

		customerStorage := customer.NewStorage(db, cacheDB)
		customerSvc := customer.NewService(customerStorage, atomicProcessor, eventBus)

		accountingStorage := accounting.NewStorage(db, cacheDB)
		accountingSvc := accounting.NewService(accountingStorage, atomicProcessor, eventBus)
		accounting.NewBusHandler(eventBus, accountingSvc, businessSvc)

		orderStorage := order.NewStorage(db, nil)
		orderSvc := order.NewService(orderStorage, atomicProcessor, eventBus, inventorySvc, customerSvc, businessSvc)

		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		deps := seedDeps{
			accountSvc:    accountSvc,
			accountStore:  accountStorage,
			billingSvc:    billingSvc,
			businessSvc:   businessSvc,
			inventorySvc:  inventorySvc,
			customerSvc:   customerSvc,
			accountingSvc: accountingSvc,
			orderSvc:      orderSvc,
			billingStore:  billingStorage,
			stripeEnabled: stripeEnabled,
			rng:           rng,
		}

		fmt.Printf("🌱 Seeding multi-workspace dataset (%s)\n", sz)

		if seedClean {
			if err := step("Cleaning database", func() error { return truncateAllPublicTables(ctx, db) }); err != nil {
				return err
			}
		}

		passwordHash, err := hash.Password(seedPassword)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		if err := step("Syncing billing plans to database", func() error {
			return billingSvc.SyncPlans(ctx)
		}); err != nil {
			return err
		}
		if stripeEnabled {
			_ = step("Syncing billing plans to Stripe (best-effort)", func() error {
				return billingSvc.SyncPlansToStripe(ctx)
			})
		}

		configs := buildWorkspaceConfigs(baseCounts)
		results := make([]*seedWorkspaceResult, 0, len(configs))

		for _, cfg := range configs {
			res, err := seedWorkspace(ctx, deps, cfg, passwordHash)
			if err != nil {
				return err
			}
			results = append(results, res)
		}

		fmt.Println("✅ Seed complete")
		fmt.Println("---")
		for _, res := range results {
			fmt.Printf("Workspace: %s\n", res.Workspace.ID)
			fmt.Printf("Business: %s (%s)\n", res.Business.ID, res.Business.Descriptor)
			fmt.Printf("Owner: %s (%s)\n", res.Owner.Email, res.Owner.ID)
			fmt.Printf("Orders: %d (paid %d, fulfilled %d)\n", res.CreatedOrders, res.PaidOrders, res.FulfilledOrders)
			fmt.Println("---")
		}
		fmt.Println("Default owner password:", seedPassword)
		return nil
	},
}

func init() {
	seedCmd.Flags().BoolVar(&seedClean, "clean", false, "truncate all tables before seeding")
	seedCmd.Flags().StringVar(&seedSizeFlag, "size", string(seedSmall), "seed size: small|medium|large")
	seedCmd.Flags().StringVar(&seedPassword, "password", "KyoraDev@123", "password for seeded owner account")
	rootCmd.AddCommand(seedCmd)
}

func loadStripeConfig() (string, string) {
	stripeKey := strings.TrimSpace(viper.GetString(config.StripeAPIKey))
	baseURL := strings.TrimSpace(os.Getenv("KYORA_STRIPE_BASE_URL"))
	if baseURL != "" {
		if override := strings.TrimSpace(os.Getenv("KYORA_STRIPE_API_KEY")); override != "" {
			stripeKey = override
		} else {
			stripeKey = "sk_test_123"
		}
	} else if strings.HasPrefix(stripeKey, "sk_test_") {
		// Default to local stripe-mock when using a test key and no override is provided.
		baseURL = "http://localhost:12111"
	}
	return stripeKey, baseURL
}

func initStripeClient(stripeKey, baseURL string) {
	stripelib.Key = stripeKey
	stripelib.SetAppInfo(&stripelib.AppInfo{Name: "Kyora", Version: "1.0", URL: "https://github.com/abdelrahman146/kyora"})
	if baseURL != "" {
		backend := stripelib.GetBackendWithConfig(stripelib.APIBackend, &stripelib.BackendConfig{URL: &baseURL})
		stripelib.SetBackend(stripelib.APIBackend, backend)
		slog.Info("Stripe base URL overridden", "baseURL", baseURL)
	}
	slog.Info("Stripe client initialized")
}

func buildWorkspaceConfigs(base seedCounts) []workspaceSeedConfig {
	freelancerCounts := scaleCounts(base, 0.7)
	growingCounts := scaleCounts(base, 1.5)
	enterpriseCounts := scaleCounts(base, 3.0)

	primaryTheme := business.StorefrontTheme{
		PrimaryColor:      "#2563EB",
		SecondaryColor:    "#10B981",
		AccentColor:       "#2563EB",
		BackgroundColor:   "#F8FAFC",
		TextColor:         "#0F172A",
		FontFamily:        "Inter",
		HeadingFontFamily: "Inter",
	}

	boldTheme := business.StorefrontTheme{
		PrimaryColor:      "#7C3AED",
		SecondaryColor:    "#F59E0B",
		AccentColor:       "#7C3AED",
		BackgroundColor:   "#0F172A",
		TextColor:         "#F8FAFC",
		FontFamily:        "Inter",
		HeadingFontFamily: "Poppins",
	}

	return []workspaceSeedConfig{
		{
			Name:              "Freelancer",
			Slug:              "freelancer",
			BusinessName:      "Freelancer Goods",
			Descriptor:        "freelancer",
			Country:           "US",
			Currency:          "USD",
			Brand:             "Freelancer Co",
			PlanDescriptor:    "starter",
			SupportEmail:      "freelancer.support@kyora.dev",
			PhoneNumber:       "+12025550111",
			WhatsappNumber:    "+12025550111",
			StorefrontEnabled: true,
			StorefrontTheme:   primaryTheme,
			Counts:            freelancerCounts,
			TeamMembers:       []seedTeamMember{},
		},
		{
			Name:              "Growing Startup",
			Slug:              "startup",
			BusinessName:      "Atlas Outfitters",
			Descriptor:        "atlas",
			Country:           "AE",
			Currency:          "AED",
			Brand:             "Atlas",
			PlanDescriptor:    "professional",
			SupportEmail:      "startup.support@kyora.dev",
			PhoneNumber:       "+971500000111",
			WhatsappNumber:    "+971500000111",
			StorefrontEnabled: true,
			StorefrontTheme:   boldTheme,
			Counts:            growingCounts,
			TeamMembers: []seedTeamMember{
				{FirstName: "Maya", LastName: "Chen", Email: "maya.chen@kyora.dev", Role: role.RoleUser},
			},
		},
		{
			Name:              "Enterprise Corp",
			Slug:              "enterprise",
			BusinessName:      "Northwind Enterprise",
			Descriptor:        "northwind",
			Country:           "GB",
			Currency:          "GBP",
			Brand:             "Northwind",
			PlanDescriptor:    "enterprise",
			SupportEmail:      "enterprise.support@kyora.dev",
			PhoneNumber:       "+442080000111",
			WhatsappNumber:    "+442080000111",
			StorefrontEnabled: false,
			StorefrontTheme:   primaryTheme,
			Counts:            enterpriseCounts,
			TeamMembers: []seedTeamMember{
				{FirstName: "Aiden", LastName: "Stone", Email: "aiden.stone@kyora.dev", Role: role.RoleAdmin},
				{FirstName: "Nora", LastName: "Patel", Email: "nora.patel@kyora.dev", Role: role.RoleUser},
				{FirstName: "Leo", LastName: "Garcia", Email: "leo.garcia@kyora.dev", Role: role.RoleUser},
			},
		},
	}
}

func scaleCounts(base seedCounts, factor float64) seedCounts {
	safeCeil := func(v int) int {
		return int(math.Max(1, math.Ceil(float64(v)*factor)))
	}
	return seedCounts{
		categories: safeCeil(base.categories),
		products:   safeCeil(base.products),
		customers:  safeCeil(base.customers),
		orders:     safeCeil(base.orders),
	}
}

func seedWorkspace(ctx context.Context, deps seedDeps, cfg workspaceSeedConfig, passwordHash string) (*seedWorkspaceResult, error) {
	label := fmt.Sprintf("[%s]", cfg.Name)
	rng := deps.rng

	ownerEmail := fmt.Sprintf("%s.owner@kyora.dev", cfg.Slug)
	var owner *account.User
	var ws *account.Workspace
	if err := step(label+" Creating workspace + owner", func() error {
		u, w, err := deps.accountSvc.BootstrapWorkspaceAndOwner(ctx, cfg.Name, "Owner", ownerEmail, passwordHash, true, "")
		if err != nil {
			return fmt.Errorf("failed to create workspace/owner (try --clean): %w", err)
		}
		owner = u
		ws = w
		return nil
	}); err != nil {
		return nil, err
	}

	biz, err := createBusinessWithRetries(ctx, deps.businessSvc, owner, cfg, rng)
	if err != nil {
		return nil, err
	}

	if err := step(label+" Enabling payment methods", func() error {
		enabled := true
		_, err := deps.businessSvc.UpdatePaymentMethod(ctx, owner, biz, "bank_transfer", &business.UpdateBusinessPaymentMethodRequest{Enabled: &enabled})
		return err
	}); err != nil {
		return nil, err
	}

	if err := step(label+" Updating storefront", func() error {
		input := &business.UpdateBusinessInput{
			StorefrontEnabled: &cfg.StorefrontEnabled,
			StorefrontTheme:   &cfg.StorefrontTheme,
		}
		if _, err := deps.businessSvc.UpdateBusiness(ctx, owner, biz.ID, input); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	shippingZones, err := seedShippingZones(ctx, deps.businessSvc, owner, biz)
	if err != nil {
		return nil, err
	}

	if err := seedStripeForWorkspace(ctx, deps, ws, biz, cfg.PlanDescriptor, label); err != nil {
		return nil, err
	}

	_, variants, err := seedInventoryData(ctx, deps.inventorySvc, owner, biz, cfg.Counts, rng)
	if err != nil {
		return nil, err
	}

	customers, addresses, err := seedCustomersAndAddresses(ctx, deps.customerSvc, owner, biz, cfg.Counts.customers, cfg.Country, rng)
	if err != nil {
		return nil, err
	}

	orders, stats, err := seedOrders(ctx, deps.orderSvc, owner, biz, cfg.Counts.orders, variants, customers, addresses, shippingZones, rng)
	if err != nil {
		return nil, err
	}

	if err := seedAccountingData(ctx, deps.accountingSvc, owner, biz, orders, rng); err != nil {
		return nil, err
	}

	if err := seedTeamMembers(ctx, deps, owner, ws, cfg.TeamMembers, cfg.Name); err != nil {
		slog.Warn("failed to seed some team members", "workspace", ws.ID, "error", err)
	}

	return &seedWorkspaceResult{
		Workspace:       ws,
		Business:        biz,
		Owner:           owner,
		CreatedOrders:   stats.Created,
		PaidOrders:      stats.Paid,
		FulfilledOrders: stats.Fulfilled,
	}, nil
}

func seedStripeForWorkspace(ctx context.Context, deps seedDeps, ws *account.Workspace, biz *business.Business, planDescriptor, label string) error {
	plan, err := deps.billingSvc.GetPlanByDescriptor(ctx, planDescriptor)
	if err != nil {
		return fmt.Errorf("failed to load plan %s: %w", planDescriptor, err)
	}

	bestEffort := func(title string, fn func() error) {
		_ = step(label+" "+title, func() error {
			if !deps.stripeEnabled {
				slog.Info("Skipping Stripe step", "reason", "stripe disabled", "workspace", ws.ID)
				return nil
			}
			if err := fn(); err != nil {
				slog.Warn("Stripe step failed", "step", title, "error", err)
			}
			return nil
		})
	}

	bestEffort("Creating Stripe customer", func() error {
		_, err := deps.billingSvc.EnsureCustomer(ctx, ws)
		return err
	})

	bestEffort("Creating default Stripe payment method", func() error {
		pm, err := paymentmethod.New(&stripelib.PaymentMethodParams{
			Type: stripelib.String("card"),
			Card: &stripelib.PaymentMethodCardParams{Token: stripelib.String("tok_visa")},
		})
		if err != nil {
			return err
		}
		if pm.ID == "" {
			return errors.New("payment method not created")
		}
		if ws.StripeCustomerID.Valid {
			if _, err := paymentmethod.Attach(pm.ID, &stripelib.PaymentMethodAttachParams{Customer: stripelib.String(ws.StripeCustomerID.String)}); err != nil {
				return err
			}
		}
		if err := deps.accountSvc.SetWorkspaceDefaultPaymentMethod(ctx, ws.ID, pm.ID); err != nil {
			return err
		}
		return deps.billingSvc.AttachAndSetDefaultPaymentMethod(ctx, ws, pm.ID)
	})

	bestEffort("Creating subscription", func() error {
		_, err := deps.billingSvc.CreateOrUpdateSubscription(ctx, ws, plan)
		return err
	})

	bestEffort("Syncing Stripe invoices", func() error {
		return seedStripeInvoices(ctx, deps, ws)
	})

	return nil
}

func seedStripeInvoices(ctx context.Context, deps seedDeps, ws *account.Workspace) error {
	if !deps.stripeEnabled {
		return nil
	}
	custID, err := deps.billingSvc.EnsureCustomer(ctx, ws)
	if err != nil {
		return err
	}
	params := &stripelib.InvoiceListParams{
		Customer:   stripelib.String(custID),
		ListParams: stripelib.ListParams{Limit: stripelib.Int64(20)},
	}
	iter := invoice.List(params)
	for iter.Next() {
		inv := iter.Invoice()
		if inv.Status == stripelib.InvoiceStatusDraft {
			if finalized, err := invoice.FinalizeInvoice(inv.ID, nil); err == nil && finalized != nil {
				inv = finalized
			} else if err != nil {
				slog.Warn("failed to finalize invoice", "invoice", inv.ID, "error", err)
			}
		}
		if inv.Status == stripelib.InvoiceStatusOpen {
			if _, err := invoice.Pay(inv.ID, &stripelib.InvoicePayParams{}); err != nil {
				slog.Warn("failed to pay invoice during seed", "invoice", inv.ID, "error", err)
			}
		}
		hostedURL := inv.HostedInvoiceURL
		invoicePDF := inv.InvoicePDF
		if err := deps.billingStore.UpsertInvoiceRecord(ctx, ws.ID, inv.ID, hostedURL, invoicePDF); err != nil {
			slog.Warn("failed to persist invoice record", "invoice", inv.ID, "workspace", ws.ID, "error", err)
		}
	}
	return iter.Err()
}

func createBusinessWithRetries(ctx context.Context, svc *business.Service, owner *account.User, cfg workspaceSeedConfig, rng *rand.Rand) (*business.Business, error) {
	var biz *business.Business
	if err := step("["+cfg.Name+"] Creating business", func() error {
		descriptor := cfg.Descriptor
		for attempt := 0; attempt < 10; attempt++ {
			candidate := descriptor
			if attempt > 0 {
				candidate = fmt.Sprintf("%s-%d", descriptor, rng.Intn(10_000))
			}
			decor := randomBusinessDecor(candidate, rng)
			available, err := svc.IsBusinessDescriptorAvailable(ctx, owner, candidate)
			if err != nil {
				return err
			}
			if !available {
				continue
			}
			var logo *assetTypes.AssetReference
			if decor.logoURL != "" {
				// Generate both CDN and original URLs for logo
				origURL := decor.logoURL
				logo = &assetTypes.AssetReference{
					URL:         decor.logoURL, // CDN URL (primary)
					OriginalURL: &origURL,      // Storage URL (same for placeholder.pics)
				}
			}
			created, err := svc.CreateBusiness(ctx, owner, &business.CreateBusinessInput{
				Name:              cfg.BusinessName,
				Descriptor:        candidate,
				Brand:             cfg.Brand,
				Logo:              logo,
				CountryCode:       cfg.Country,
				Currency:          cfg.Currency,
				VatRate:           decimal.NewFromFloat(0.15),
				StorefrontEnabled: cfg.StorefrontEnabled,
				SupportEmail:      cfg.SupportEmail,
				PhoneNumber:       cfg.PhoneNumber,
				WhatsappNumber:    cfg.WhatsappNumber,
				Address:           fmt.Sprintf("%s HQ", cfg.BusinessName),
				WebsiteURL:        decor.websiteURL,
				InstagramURL:      decor.instagramURL,
				FacebookURL:       decor.facebookURL,
				TikTokURL:         decor.tiktokURL,
				XURL:              decor.xURL,
				SnapchatURL:       decor.snapchatURL,
			})
			if err != nil {
				return err
			}
			biz = created
			return nil
		}
		return errors.New("failed to create business after retries")
	}); err != nil {
		return nil, err
	}
	return biz, nil
}

func seedShippingZones(ctx context.Context, svc *business.Service, owner *account.User, biz *business.Business) ([]*business.ShippingZone, error) {
	zones := make([]*business.ShippingZone, 0, 2)
	domesticCost := decimal.NewFromFloat(10)
	internationalCost := decimal.NewFromFloat(35)
	freeThreshold := decimal.NewFromFloat(150)

	retryRateLimited := func(fn func() error) error {
		for attempt := 0; attempt < 3; attempt++ {
			if err := fn(); err != nil {
				var p *problem.Problem
				if errors.As(err, &p) && p.Status == http.StatusTooManyRequests {
					time.Sleep(1200 * time.Millisecond)
					continue
				}
				return err
			}
			return nil
		}
		return problem.TooManyRequests("too many requests")
	}

	var domestic *business.ShippingZone
	err := retryRateLimited(func() error {
		var createErr error
		domestic, createErr = svc.CreateShippingZone(ctx, owner, biz, &business.CreateShippingZoneRequest{
			Name:                  "Domestic",
			Countries:             []string{strings.ToUpper(biz.CountryCode)},
			ShippingCost:          domesticCost,
			FreeShippingThreshold: freeThreshold,
		})
		return createErr
	})
	if err != nil && !database.IsUniqueViolation(err) {
		return nil, err
	}
	if err == nil {
		zones = append(zones, domestic)
	}
	intlCountries := []string{"GB", "DE", "AE", "US", "SA", "FR"}
	if biz.CountryCode != "" {
		for i, c := range intlCountries {
			if strings.EqualFold(c, biz.CountryCode) {
				intlCountries = append(intlCountries[:i], intlCountries[i+1:]...)
				break
			}
		}
	}
	var international *business.ShippingZone
	err = retryRateLimited(func() error {
		var createErr error
		international, createErr = svc.CreateShippingZone(ctx, owner, biz, &business.CreateShippingZoneRequest{
			Name:         "International",
			Countries:    intlCountries,
			ShippingCost: internationalCost,
		})
		return createErr
	})
	if err != nil && !database.IsUniqueViolation(err) {
		return nil, err
	}
	if err == nil {
		zones = append(zones, international)
	}

	if len(zones) == 0 {
		existing, listErr := svc.ListShippingZones(ctx, owner, biz)
		if listErr == nil && len(existing) > 0 {
			zones = append(zones, existing...)
		}
	}
	return zones, nil
}

func seedInventoryData(ctx context.Context, svc *inventory.Service, owner *account.User, biz *business.Business, counts seedCounts, rng *rand.Rand) ([]*inventory.Category, []*inventory.Variant, error) {
	categoryNames := []string{"Apparel", "Home & Living", "Beauty", "Electronics", "Outdoors", "Stationery", "Food & Drink"}
	categories := make([]*inventory.Category, 0, counts.categories)
	if err := step("["+biz.Descriptor+"] Creating categories", func() error {
		for i := 0; i < counts.categories; i++ {
			name := categoryNames[i%len(categoryNames)]
			cat, err := svc.CreateCategory(ctx, owner, biz, &inventory.CreateCategoryRequest{
				Name:       name,
				Descriptor: fmt.Sprintf("%s-%d", strings.ToLower(strings.ReplaceAll(name, " ", "-")), i+1),
			})
			if err != nil {
				return err
			}
			categories = append(categories, cat)
		}
		return nil
	}); err != nil {
		return nil, nil, err
	}

	variants := make([]*inventory.Variant, 0, counts.products*2)
	adjectives := []string{"Premium", "Lightweight", "Eco", "Classic", "Modern", "Signature", "Artisan"}
	if err := step("["+biz.Descriptor+"] Creating products + variants", func() error {
		for i := 0; i < counts.products; i++ {
			cat := categories[i%len(categories)]
			cost := decimal.NewFromFloat(15 + rng.Float64()*60).Round(2)
			price := cost.Mul(decimal.NewFromFloat(1.8)).Round(2)
			premium := price.Mul(decimal.NewFromFloat(1.25)).Round(2)
			stock := 100 + rng.Intn(250)
			alert := 15 + rng.Intn(20)
			name := fmt.Sprintf("%s %s", adjectives[i%len(adjectives)], cat.Name)
			description := "Hand-picked product seeded for realistic demos"
			photos := []string{
				placeholderImageURL(name),
				placeholderImageURL(name + " Variant"),
			}
			standardPhotos := []string{placeholderVariantImageURL(name, "standard")}
			premiumPhotos := []string{placeholderVariantImageURL(name, "premium")}

			// Convert to AssetReference with full CDN and thumbnail support
			photoRefs := make([]assetTypes.AssetReference, len(photos))
			for i, url := range photos {
				// Generate thumbnail URL (smaller dimension for placeholder.pics)
				thumbURL := strings.Replace(url, "/800/600", "/512/512", 1)
				origURL := url
				thumbOrigURL := thumbURL
				photoRefs[i] = assetTypes.AssetReference{
					URL:                  url,           // CDN URL (primary)
					OriginalURL:          &origURL,      // Storage URL (same for placeholder.pics)
					ThumbnailURL:         &thumbURL,     // CDN thumbnail URL
					ThumbnailOriginalURL: &thumbOrigURL, // Storage thumbnail URL
				}
			}
			standardPhotoRefs := make([]assetTypes.AssetReference, len(standardPhotos))
			for i, url := range standardPhotos {
				thumbURL := strings.Replace(url, "/800/600", "/512/512", 1)
				origURL := url
				thumbOrigURL := thumbURL
				standardPhotoRefs[i] = assetTypes.AssetReference{
					URL:                  url,
					OriginalURL:          &origURL,
					ThumbnailURL:         &thumbURL,
					ThumbnailOriginalURL: &thumbOrigURL,
				}
			}
			premiumPhotoRefs := make([]assetTypes.AssetReference, len(premiumPhotos))
			for i, url := range premiumPhotos {
				thumbURL := strings.Replace(url, "/800/600", "/512/512", 1)
				origURL := url
				thumbOrigURL := thumbURL
				premiumPhotoRefs[i] = assetTypes.AssetReference{
					URL:                  url,
					OriginalURL:          &origURL,
					ThumbnailURL:         &thumbURL,
					ThumbnailOriginalURL: &thumbOrigURL,
				}
			}

			p, err := svc.CreateProductWithVariants(ctx, owner, biz, &inventory.CreateProductWithVariantsRequest{
				Product: inventory.CreateProductRequest{
					Name:        name,
					Description: description,
					Photos:      photoRefs,
					CategoryID:  cat.ID,
				},
				Variants: []inventory.CreateProductVariantRequest{
					{Code: "standard", Photos: standardPhotoRefs, CostPrice: &cost, SalePrice: &price, StockQuantity: &stock, StockQuantityAlert: &alert},
					{Code: "premium", Photos: premiumPhotoRefs, CostPrice: &cost, SalePrice: &premium, StockQuantity: &stock, StockQuantityAlert: &alert},
				},
			})
			if err != nil {
				return err
			}
			variants = append(variants, p.Variants...)
		}
		return nil
	}); err != nil {
		return nil, nil, err
	}

	return categories, variants, nil
}

func seedCustomersAndAddresses(ctx context.Context, svc *customer.Service, owner *account.User, biz *business.Business, count int, homeCountry string, rng *rand.Rand) ([]*customer.Customer, map[string][]*customer.CustomerAddress, error) {
	firstNames := []string{"Layla", "Omar", "Zoe", "Hassan", "Mia", "Jonas", "Rafael", "Noor", "Sofia", "Diego"}
	lastNames := []string{"Khan", "Lopez", "Ibrahim", "Chen", "Patel", "Garcia", "Stone", "Okafor", "Silva", "Laurent"}
	cities := []string{"Riyadh", "Dubai", "London", "San Francisco", "Berlin", "Paris", "Madrid", "Toronto"}
	addressesByCustomer := make(map[string][]*customer.CustomerAddress, count)
	customers := make([]*customer.Customer, 0, count)
	if err := step("["+biz.Descriptor+"] Creating customers + addresses", func() error {
		for i := 0; i < count; i++ {
			name := fmt.Sprintf("%s %s", firstNames[i%len(firstNames)], lastNames[rng.Intn(len(lastNames))])
			email := fmt.Sprintf("%s.%s.%d@kyora.dev", strings.ToLower(firstNames[i%len(firstNames)]), strings.ToLower(lastNames[i%len(lastNames)]), i)
			country := homeCountry
			if rng.Float64() < 0.35 {
				altCountries := []string{"US", "AE", "GB", "DE", "FR"}
				country = altCountries[rng.Intn(len(altCountries))]
			}
			handles := randomCustomerHandles(strings.Split(email, "@")[0], rng)
			genders := []customer.CustomerGender{customer.GenderMale, customer.GenderFemale, customer.GenderOther}
			gender := genders[rng.Intn(len(genders))]
			c, err := svc.CreateCustomer(ctx, owner, biz, &customer.CreateCustomerRequest{
				Name:              name,
				Gender:            gender,
				Email:             email,
				CountryCode:       country,
				PhoneCode:         "+1",
				PhoneNumber:       fmt.Sprintf("555%06d", rng.Intn(999999)),
				TikTokUsername:    handles.tiktok,
				InstagramUsername: handles.instagram,
				FacebookUsername:  handles.facebook,
				XUsername:         handles.x,
				WhatsappNumber:    handles.whatsapp,
			})
			if err != nil {
				return err
			}
			addr, err := svc.CreateCustomerAddress(ctx, owner, biz, c.ID, &customer.CreateCustomerAddressRequest{
				Street:      fmt.Sprintf("%d %s St", 10+rng.Intn(200), lastNames[i%len(lastNames)]),
				City:        cities[rng.Intn(len(cities))],
				State:       "",
				ZipCode:     fmt.Sprintf("%05d", rng.Intn(99999)),
				CountryCode: country,
				PhoneCode:   "+1",
				Phone:       fmt.Sprintf("555%06d", rng.Intn(999999)),
			})
			if err != nil {
				return err
			}
			customers = append(customers, c)
			addressesByCustomer[c.ID] = []*customer.CustomerAddress{addr}
		}
		return nil
	}); err != nil {
		return nil, nil, err
	}

	return customers, addressesByCustomer, nil
}

type orderStats struct {
	Created   int
	Paid      int
	Fulfilled int
}

func seedOrders(ctx context.Context, svc *order.Service, owner *account.User, biz *business.Business, count int, variants []*inventory.Variant, customers []*customer.Customer, addresses map[string][]*customer.CustomerAddress, zones []*business.ShippingZone, rng *rand.Rand) ([]*order.Order, orderStats, error) {
	channels := []string{"instagram", "tiktok", "whatsapp", "facebook"}
	stats := orderStats{}
	orders := make([]*order.Order, 0, count)
	if err := step("["+biz.Descriptor+"] Creating orders", func() error {
		for i := 0; i < count; i++ {
			c := customers[rng.Intn(len(customers))]
			addr := addresses[c.ID][0]
			itemsCount := 1 + rng.Intn(3)
			items := make([]*order.CreateOrderItemRequest, 0, itemsCount)
			for j := 0; j < itemsCount; j++ {
				v := variants[rng.Intn(len(variants))]
				qty := 1 + rng.Intn(4)
				items = append(items, &order.CreateOrderItemRequest{
					VariantID: v.ID,
					Quantity:  qty,
					UnitPrice: v.SalePrice,
					UnitCost:  v.CostPrice,
				})
			}

			shippingFee := decimal.NewFromFloat(0)
			var zoneID *string
			if len(zones) > 0 {
				for _, z := range zones {
					if z.Countries.Contains(strings.ToUpper(addr.CountryCode)) {
						zoneID = &z.ID
						shippingFee = z.ShippingCost
						break
					}
				}
				if zoneID == nil {
					zoneID = &zones[len(zones)-1].ID
					shippingFee = zones[len(zones)-1].ShippingCost
				}
			}

			discount := decimal.NewFromFloat(0)
			if rng.Float64() < 0.2 {
				discount = decimal.NewFromFloat(5 + rng.Float64()*15).Round(2)
			}

			orderedAt := randomPastTime(rng, 12)
			ord, err := svc.CreateOrder(ctx, owner, biz, &order.CreateOrderRequest{
				CustomerID:        c.ID,
				Channel:           channels[rng.Intn(len(channels))],
				ShippingAddressID: addr.ID,
				ShippingZoneID:    zoneID,
				ShippingFee:       shippingFee,
				Discount:          discount,
				PaymentMethod:     order.OrderPaymentMethodBankTransfer,
				OrderedAt:         orderedAt,
				Items:             items,
			})
			if err != nil {
				return err
			}
			orders = append(orders, ord)
			stats.Created++

			progress := rng.Float64()
			if progress < 0.15 {
				updated, err := svc.UpdateOrderStatus(ctx, owner, biz, ord.ID, order.OrderStatusCancelled)
				if err != nil {
					return err
				}
				ord = updated
				orders[len(orders)-1] = ord
				continue
			}

			updatedOrder, err := svc.UpdateOrderStatus(ctx, owner, biz, ord.ID, order.OrderStatusPlaced)
			if err != nil {
				return err
			}
			ord = updatedOrder

			payProbability := 0.65
			if progress > 0.6 {
				payProbability = 0.9
			}
			if rng.Float64() < payProbability {
				updated, err := svc.UpdateOrderPaymentStatus(ctx, owner, biz, ord.ID, order.OrderPaymentStatusPaid)
				if err != nil {
					return err
				}
				ord = updated
				stats.Paid++
			}

			if progress > 0.45 {
				_, err := svc.UpdateOrderStatus(ctx, owner, biz, ord.ID, order.OrderStatusShipped)
				if err != nil {
					return err
				}
				fulfilled, err := svc.UpdateOrderStatus(ctx, owner, biz, ord.ID, order.OrderStatusFulfilled)
				if err != nil {
					return err
				}
				ord = fulfilled
				stats.Fulfilled++
			}

			orders[len(orders)-1] = ord
		}
		return nil
	}); err != nil {
		return nil, stats, err
	}

	return orders, stats, nil
}

func seedAccountingData(ctx context.Context, svc *accounting.Service, owner *account.User, biz *business.Business, orders []*order.Order, rng *rand.Rand) error {
	assets := []struct {
		Name string
		Type accounting.AssetType
		Cost float64
	}{
		{Name: "MacBook Pro", Type: accounting.AssetTypeEquipment, Cost: 2800},
		{Name: "Office Furniture", Type: accounting.AssetTypeFurniture, Cost: 1500},
		{Name: "Studio Lighting", Type: accounting.AssetTypeEquipment, Cost: 800},
	}

	if err := step("["+biz.Descriptor+"] Accounting: assets", func() error {
		for _, a := range assets {
			_, err := svc.CreateAsset(ctx, owner, biz, &accounting.CreateAssetRequest{
				Name:        a.Name,
				Type:        a.Type,
				Value:       decimal.NewFromFloat(a.Cost).Round(2),
				PurchasedAt: randomPastTime(rng, 10),
				Note:        "Seeded fixed asset",
			})
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if err := step("["+biz.Descriptor+"] Accounting: investments & withdrawals", func() error {
		_, err := svc.CreateInvestment(ctx, owner, biz, &accounting.CreateInvestmentRequest{
			InvestorID: owner.ID,
			Amount:     decimal.NewFromFloat(20000 + rng.Float64()*15000).Round(2),
			Note:       "Initial capital injection",
			InvestedAt: randomPastTime(rng, 12),
		})
		if err != nil {
			return err
		}
		_, err = svc.CreateWithdrawal(ctx, owner, biz, &accounting.CreateWithdrawalRequest{
			WithdrawerID: owner.ID,
			Amount:       decimal.NewFromFloat(1500 + rng.Float64()*1500).Round(2),
			Note:         "Owner draw",
			WithdrawnAt:  randomPastTime(rng, 6),
		})
		return err
	}); err != nil {
		return err
	}

	renting := accounting.CreateRecurringExpenseRequest{
		Frequency:                    accounting.RecurringExpenseFrequencyMonthly,
		RecurringStartDate:           time.Now().AddDate(0, -6, 0),
		Amount:                       decimal.NewFromFloat(1200),
		Category:                     accounting.ExpenseCategoryRent,
		Note:                         "Monthly rent",
		AutoCreateHistoricalExpenses: true,
	}
	saasp := accounting.CreateRecurringExpenseRequest{
		Frequency:                    accounting.RecurringExpenseFrequencyMonthly,
		RecurringStartDate:           time.Now().AddDate(0, -6, 0),
		Amount:                       decimal.NewFromFloat(320),
		Category:                     accounting.ExpenseCategorySoftware,
		Note:                         "Software subscriptions",
		AutoCreateHistoricalExpenses: true,
	}
	if err := step("["+biz.Descriptor+"] Accounting: recurring expenses", func() error {
		if _, err := svc.CreateRecurringExpense(ctx, owner, biz, &renting); err != nil {
			return err
		}
		if _, err := svc.CreateRecurringExpense(ctx, owner, biz, &saasp); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	oneTimeExpenses := []accounting.CreateExpenseRequest{}
	oneTimeExpenseTemplates := []struct {
		amount   float64
		category accounting.ExpenseCategory
		note     string
		months   int
	}{
		{amount: 450, category: accounting.ExpenseCategorySupplies, note: "Office supplies", months: 8},
		{amount: 780, category: accounting.ExpenseCategoryTravel, note: "Client visits", months: 4},
		{amount: 560, category: accounting.ExpenseCategoryMarketing, note: "Ad spend", months: 3},
	}
	for _, tpl := range oneTimeExpenseTemplates {
		when := randomPastTime(rng, tpl.months)
		oneTimeExpenses = append(oneTimeExpenses, accounting.CreateExpenseRequest{
			Amount:     decimal.NewFromFloat(tpl.amount).Round(2),
			Category:   tpl.category,
			Type:       accounting.ExpenseTypeOneTime,
			Note:       tpl.note,
			OccurredOn: &when,
		})
	}
	if err := step("["+biz.Descriptor+"] Accounting: one-time expenses", func() error {
		for _, e := range oneTimeExpenses {
			if _, err := svc.CreateExpense(ctx, owner, biz, &e); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	paidOrders := make([]*order.Order, 0)
	for _, o := range orders {
		if o.PaymentStatus == order.OrderPaymentStatusPaid {
			paidOrders = append(paidOrders, o)
		}
	}
	if len(paidOrders) > 0 {
		_ = step("["+biz.Descriptor+"] Accounting: transaction fee expenses", func() error {
			for _, o := range paidOrders {
				if rng.Float64() > 0.35 {
					continue
				}
				fee := o.Total.Mul(decimal.NewFromFloat(0.02)).Round(2)
				if fee.LessThanOrEqual(decimal.Zero) {
					continue
				}
				if err := svc.UpsertTransactionFeeExpenseForOrder(ctx, o.BusinessID, o.ID, fee, o.Currency, o.OrderedAt, string(o.PaymentMethod)); err != nil {
					slog.Warn("failed to create transaction fee expense", "order", o.ID, "error", err)
				}
			}
			return nil
		})
	}

	return nil
}

func seedTeamMembers(ctx context.Context, deps seedDeps, owner *account.User, ws *account.Workspace, members []seedTeamMember, workspaceName string) error {
	for _, m := range members {
		inv, err := deps.accountSvc.InviteUserToWorkspace(ctx, owner, ws.ID, m.Email, m.Role)
		if err != nil {
			return err
		}
		token, _, err := deps.accountStore.CreateWorkspaceInvitationToken(&account.WorkspaceInvitationPayload{
			InvitationID: inv.ID,
			WorkspaceID:  ws.ID,
			Email:        m.Email,
			Role:         string(m.Role),
			InviterID:    owner.ID,
		})
		if err != nil {
			return err
		}
		if _, _, err := deps.accountSvc.AcceptInvitation(ctx, token, m.FirstName, m.LastName, seedPassword); err != nil {
			return fmt.Errorf("failed to accept invitation for %s in %s: %w", m.Email, workspaceName, err)
		}
	}
	return nil
}

func randomPastTime(rng *rand.Rand, monthsBack int) time.Time {
	if monthsBack <= 0 {
		monthsBack = 1
	}
	days := monthsBack * 30
	return time.Now().UTC().Add(-time.Duration(rng.Intn(days*24)) * time.Hour)
}

type businessDecor struct {
	logoURL      string
	websiteURL   string
	instagramURL string
	facebookURL  string
	tiktokURL    string
	xURL         string
	snapchatURL  string
}

type customerHandles struct {
	tiktok    string
	instagram string
	facebook  string
	x         string
	whatsapp  string
}

func randomBusinessDecor(descriptor string, rng *rand.Rand) businessDecor {
	seed := fmt.Sprintf("%s-%d", descriptor, rng.Intn(10_000))
	decor := businessDecor{
		logoURL: fmt.Sprintf("https://picsum.photos/seed/%s/320/320", seed),
	}
	if rng.Float64() < 0.8 {
		decor.websiteURL = fmt.Sprintf("https://%s.kyora.test", descriptor)
	}
	if rng.Float64() < 0.7 {
		decor.instagramURL = fmt.Sprintf("https://instagram.com/%s_shop", descriptor)
	}
	if rng.Float64() < 0.4 {
		decor.facebookURL = fmt.Sprintf("https://facebook.com/%s.store", descriptor)
	}
	if rng.Float64() < 0.5 {
		decor.tiktokURL = fmt.Sprintf("https://www.tiktok.com/@%s", descriptor)
	}
	if rng.Float64() < 0.35 {
		decor.xURL = fmt.Sprintf("https://x.com/%s", descriptor)
	}
	if rng.Float64() < 0.2 {
		decor.snapchatURL = fmt.Sprintf("https://www.snapchat.com/add/%s", descriptor)
	}
	return decor
}

func randomCustomerHandles(base string, rng *rand.Rand) customerHandles {
	makeHandle := func(prefix string) string {
		return fmt.Sprintf("%s_%s_%d", prefix, strings.ToLower(base), rng.Intn(10_000))
	}
	h := customerHandles{}
	if rng.Float64() < 0.5 {
		h.tiktok = makeHandle("tik")
	}
	if rng.Float64() < 0.6 {
		h.instagram = makeHandle("ig")
	}
	if rng.Float64() < 0.35 {
		h.facebook = makeHandle("fb")
	}
	if rng.Float64() < 0.4 {
		h.x = makeHandle("x")
	}
	if rng.Float64() < 0.55 {
		h.whatsapp = fmt.Sprintf("+1%s%04d", strings.Repeat("5", 3), rng.Intn(9000)+1000)
	}
	return h
}

func placeholderImageURL(text string) string {
	encoded := url.QueryEscape(text)
	return fmt.Sprintf("https://placehold.co/600x400?text=%s", encoded)
}

func placeholderVariantImageURL(productName, variantCode string) string {
	encoded := url.QueryEscape(fmt.Sprintf("%s-%s", productName, variantCode))
	return fmt.Sprintf("https://placehold.co/640x480?text=%s", encoded)
}

func step(title string, fn func() error) error {
	fmt.Printf("→ %s... ", title)
	start := time.Now()
	if err := fn(); err != nil {
		fmt.Printf("❌ (%s)\n", time.Since(start).Truncate(time.Millisecond))
		fmt.Printf("   [ERROR] %v\n", err)
		return err
	}
	fmt.Printf("done (%s)\n", time.Since(start).Truncate(time.Millisecond))
	return nil
}

func truncateAllPublicTables(ctx context.Context, db *database.Database) error {
	// TRUNCATE all tables in public schema to keep `--clean` robust as the schema evolves.
	// Uses CASCADE to handle FK relationships and RESTART IDENTITY for sequences.
	var tables []string
	if err := db.Conn(ctx).
		Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").
		Scan(&tables).Error; err != nil {
		return err
	}
	if len(tables) == 0 {
		return nil
	}
	// Keep a small allowlist for tables we never want to truncate (if present).
	skip := map[string]struct{}{
		"schema_migrations": {},
	}
	quoted := make([]string, 0, len(tables))
	for _, t := range tables {
		if _, ok := skip[t]; ok {
			continue
		}
		// Quote identifiers to avoid surprises.
		quoted = append(quoted, fmt.Sprintf("\"%s\"", t))
	}
	if len(quoted) == 0 {
		return nil
	}
	stmt := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(quoted, ", "))
	return db.Conn(ctx).Exec(stmt).Error
}
