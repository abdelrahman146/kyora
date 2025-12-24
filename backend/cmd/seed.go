package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
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
	"github.com/abdelrahman146/kyora/internal/platform/utils/hash"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	stripelib "github.com/stripe/stripe-go/v83"
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

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed local development data",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		sz := seedSize(strings.TrimSpace(strings.ToLower(seedSizeFlag)))
		if sz != seedSmall && sz != seedMedium && sz != seedLarge {
			return errors.New("invalid --size. allowed: small|medium|large")
		}
		counts := sz.counts()

		// Initialize Stripe (best-effort; seed still works without it).
		stripeKey := viper.GetString(config.StripeAPIKey)
		baseURL := strings.TrimSpace(os.Getenv("KYORA_STRIPE_BASE_URL"))
		if baseURL != "" {
			// Stripe-mock expects a specific test key. Prefer an explicit override when provided,
			// otherwise force the well-known stripe-mock test key.
			if override := strings.TrimSpace(os.Getenv("KYORA_STRIPE_API_KEY")); override != "" {
				stripeKey = override
			} else {
				stripeKey = "sk_test_123"
			}
		}

		stripeEnabled := stripeKey != ""
		if stripeEnabled {
			stripelib.Key = stripeKey
			stripelib.SetAppInfo(&stripelib.AppInfo{Name: "Kyora", Version: "1.0", URL: "https://github.com/abdelrahman146/kyora"})
			if baseURL != "" {
				backend := stripelib.GetBackendWithConfig(stripelib.APIBackend, &stripelib.BackendConfig{URL: &baseURL})
				stripelib.SetBackend(stripelib.APIBackend, backend)
				slog.Info("Stripe base URL overridden", "baseURL", baseURL)
			}
			slog.Info("Stripe client initialized")
		} else {
			slog.Warn("Stripe API key not configured; billing seed steps will be skipped")
		}

		// Initialize database connection.
		dsn := viper.GetString(config.DatabaseDSN)
		logLevel := viper.GetString(config.DatabaseLogLevel)
		db, err := database.NewConnection(dsn, logLevel)
		if err != nil {
			slog.Error("Failed to connect to database", "error", err)
			return fmt.Errorf("failed to connect to database: %w", err)
		}
		defer db.CloseConnection()

		// Initialize cache.
		servers := viper.GetStringSlice(config.CacheHosts)
		cacheDB := cache.NewConnection(servers)

		// Initialize platform dependencies.
		atomicProcessor := database.NewAtomicProcess(db)
		eventBus := bus.New()
		defer eventBus.Close()
		emailClient, err := email.New()
		if err != nil {
			slog.Error("Failed to initialize email client", "error", err)
			return fmt.Errorf("failed to initialize email client: %w", err)
		}

		// DI - storages/services.
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

		// Seed runs are intentionally high-throughput. The order service has a 1-second
		// anti-double-submit throttle backed by cache; bypass it for seeding.
		orderStorage := order.NewStorage(db, nil)
		orderSvc := order.NewService(orderStorage, atomicProcessor, eventBus, inventorySvc, customerSvc, businessSvc)

		rng := rand.New(rand.NewSource(time.Now().UnixNano()))

		fmt.Printf("🌱 Seeding (%s)\n", sz)

		if seedClean {
			if err := step("Cleaning database", func() error { return truncateAllPublicTables(ctx, db) }); err != nil {
				return err
			}
		}

		passwordHash, err := hash.Password(seedPassword)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// Seed workspace + owner.
		var owner *account.User
		var ws *account.Workspace
		ownerEmail := "owner@kyora.dev"
		if err := step("Creating workspace + owner", func() error {
			u, w, err := accountSvc.BootstrapWorkspaceAndOwner(ctx, "Dev", "Owner", ownerEmail, passwordHash, true, "")
			if err != nil {
				return fmt.Errorf("failed to create workspace/owner (try --clean): %w", err)
			}
			owner = u
			ws = w
			return nil
		}); err != nil {
			return err
		}

		// Create business.
		var biz *business.Business
		bizDescriptor := "demo"
		if err := step("Creating business", func() error {
			descriptor := bizDescriptor
			for attempt := 0; attempt < 10; attempt++ {
				try := descriptor
				if attempt > 0 {
					try = fmt.Sprintf("%s-%d", descriptor, rng.Intn(10_000))
				}
				available, err := businessSvc.IsBusinessDescriptorAvailable(ctx, owner, try)
				if err != nil {
					return err
				}
				if !available {
					continue
				}
				created, err := businessSvc.CreateBusiness(ctx, owner, &business.CreateBusinessInput{
					Name:              "Kyora Demo Shop",
					Descriptor:        try,
					Brand:             "Kyora",
					CountryCode:       "SA",
					Currency:          "SAR",
					VatRate:           decimal.NewFromFloat(0.15),
					StorefrontEnabled: true,
					SupportEmail:      "support@kyora.dev",
					PhoneNumber:       "+966500000000",
					WhatsappNumber:    "+966500000000",
					Address:           "Riyadh, Saudi Arabia",
				})
				if err != nil {
					return err
				}
				biz = created
				bizDescriptor = try
				return nil
			}
			return errors.New("failed to create business after retries")
		}); err != nil {
			return err
		}

		// Enable a payment method so orders can be created.
		if err := step("Enabling payment methods", func() error {
			enabled := true
			_, err := businessSvc.UpdatePaymentMethod(ctx, owner, biz, "bank_transfer", &business.UpdateBusinessPaymentMethodRequest{Enabled: &enabled})
			return err
		}); err != nil {
			return err
		}

		// Billing: sync plans + ensure Stripe customer + subscription.
		if stripeEnabled {
			_ = step("Syncing billing plans (best-effort)", func() error {
				return billingSvc.SyncPlansComplete(ctx)
			})
			_ = step("Creating Stripe customer (best-effort)", func() error {
				_, err := billingSvc.EnsureCustomer(ctx, ws)
				return err
			})
			_ = step("Creating default Stripe payment method (best-effort)", func() error {
				// Create a PaymentMethod using a test token, then attach and set default.
				pm, err := paymentmethod.New(&stripelib.PaymentMethodParams{
					Type: stripelib.String("card"),
					Card: &stripelib.PaymentMethodCardParams{Token: stripelib.String("tok_visa")},
				})
				if err != nil {
					return err
				}
				_, err = paymentmethod.Attach(pm.ID, &stripelib.PaymentMethodAttachParams{Customer: stripelib.String(ws.StripeCustomerID.String)})
				if err != nil {
					return err
				}
				// Persist on workspace and set default via billing service.
				if err := accountSvc.SetWorkspaceDefaultPaymentMethod(ctx, ws.ID, pm.ID); err != nil {
					return err
				}
				return billingSvc.AttachAndSetDefaultPaymentMethod(ctx, ws, pm.ID)
			})
			_ = step("Creating subscription (best-effort)", func() error {
				plan, err := billingSvc.GetPlanByDescriptor(ctx, "starter")
				if err != nil {
					return err
				}
				_, err = billingSvc.CreateOrUpdateSubscription(ctx, ws, plan)
				return err
			})
		}

		// Inventory.
		var categories []*inventory.Category
		if err := step("Creating categories", func() error {
			categories = make([]*inventory.Category, 0, counts.categories)
			for i := 0; i < counts.categories; i++ {
				cat, err := inventorySvc.CreateCategory(ctx, owner, biz, &inventory.CreateCategoryRequest{
					Name:       fmt.Sprintf("Category %d", i+1),
					Descriptor: fmt.Sprintf("cat-%d", i+1),
				})
				if err != nil {
					return err
				}
				categories = append(categories, cat)
			}
			return nil
		}); err != nil {
			return err
		}

		var variants []*inventory.Variant
		if err := step("Creating products + variants", func() error {
			variants = make([]*inventory.Variant, 0, counts.products*2)
			for i := 0; i < counts.products; i++ {
				cat := categories[i%len(categories)]
				cost := decimal.NewFromInt(int64(20 + rng.Intn(50)))
				price := cost.Mul(decimal.NewFromFloat(1.8)).Round(2)
				premium := price.Mul(decimal.NewFromFloat(1.25)).Round(2)
				q := 250
				alert := 20
				p, err := inventorySvc.CreateProductWithVariants(ctx, owner, biz, &inventory.CreateProductWithVariantsRequest{
					Product: inventory.CreateProductRequest{
						Name:        fmt.Sprintf("Product %d", i+1),
						Description: "Seeded product",
						Photos:      []string{},
						CategoryID:  cat.ID,
					},
					Variants: []inventory.CreateProductVariantRequest{
						{Code: "default", CostPrice: &cost, SalePrice: &price, StockQuantity: &q, StockQuantityAlert: &alert},
						{Code: "premium", CostPrice: &cost, SalePrice: &premium, StockQuantity: &q, StockQuantityAlert: &alert},
					},
				})
				if err != nil {
					return err
				}
				variants = append(variants, p.Variants...)
			}
			return nil
		}); err != nil {
			return err
		}

		// Customers.
		var customers []*customer.Customer
		var addressesByCustomer map[string][]*customer.CustomerAddress
		if err := step("Creating customers + addresses", func() error {
			customers = make([]*customer.Customer, 0, counts.customers)
			addressesByCustomer = make(map[string][]*customer.CustomerAddress, counts.customers)
			for i := 0; i < counts.customers; i++ {
				c, err := customerSvc.CreateCustomer(ctx, owner, biz, &customer.CreateCustomerRequest{
					Name:        fmt.Sprintf("Customer %d", i+1),
					Email:       fmt.Sprintf("customer%d@kyora.dev", i+1),
					CountryCode: "SA",
					PhoneCode:   "+966",
					PhoneNumber: fmt.Sprintf("500%06d", i+1),
				})
				if err != nil {
					return err
				}
				customers = append(customers, c)
				addr, err := customerSvc.CreateCustomerAddress(ctx, owner, biz, c.ID, &customer.CreateCustomerAddressRequest{
					Street:      fmt.Sprintf("Street %d", i+1),
					City:        "Riyadh",
					State:       "Riyadh",
					ZipCode:     "11564",
					CountryCode: "SA",
					PhoneCode:   "+966",
					Phone:       fmt.Sprintf("500%06d", i+1),
				})
				if err != nil {
					return err
				}
				addressesByCustomer[c.ID] = []*customer.CustomerAddress{addr}
			}
			return nil
		}); err != nil {
			return err
		}

		// Orders.
		createdOrders := 0
		paidOrders := 0
		fulfilledOrders := 0
		if err := step("Creating orders", func() error {
			for i := 0; i < counts.orders; i++ {
				c := customers[rng.Intn(len(customers))]
				addr := addressesByCustomer[c.ID][0]

				itemsCount := 1 + rng.Intn(3)
				items := make([]*order.CreateOrderItemRequest, 0, itemsCount)
				for j := 0; j < itemsCount; j++ {
					v := variants[rng.Intn(len(variants))]
					qty := 1 + rng.Intn(3)
					items = append(items, &order.CreateOrderItemRequest{
						VariantID: v.ID,
						Quantity:  qty,
						UnitPrice: v.SalePrice,
						UnitCost:  v.CostPrice,
					})
				}

				orderedAt := time.Now().UTC().Add(-time.Duration(rng.Intn(45*24)) * time.Hour)
				ord, err := orderSvc.CreateOrder(ctx, owner, biz, &order.CreateOrderRequest{
					CustomerID:        c.ID,
					Channel:           "instagram",
					ShippingAddressID: addr.ID,
					ShippingFee:       decimal.NewFromInt(0),
					Discount:          decimal.NewFromInt(0),
					PaymentMethod:     order.OrderPaymentMethodBankTransfer,
					OrderedAt:         orderedAt,
					Items:             items,
				})
				if err != nil {
					return err
				}
				createdOrders++

				// Mark some orders as paid/fulfilled for history.
				// Note: Payment can only be set to paid for orders in Placed/Shipped/Fulfilled.
				shouldFulfill := rng.Float64() < 0.5
				shouldPay := rng.Float64() < 0.7
				if shouldFulfill {
					// Fulfilled orders should be paid in seeded data.
					shouldPay = true
				}
				if shouldPay || shouldFulfill {
					if _, err := orderSvc.UpdateOrderStatus(ctx, owner, biz, ord.ID, order.OrderStatusPlaced); err != nil {
						return err
					}
				}
				if shouldPay {
					if _, err := orderSvc.UpdateOrderPaymentStatus(ctx, owner, biz, ord.ID, order.OrderPaymentStatusPaid); err != nil {
						return err
					}
					paidOrders++
				}
				if shouldFulfill {
					if _, err := orderSvc.UpdateOrderStatus(ctx, owner, biz, ord.ID, order.OrderStatusShipped); err != nil {
						return err
					}
					if _, err := orderSvc.UpdateOrderStatus(ctx, owner, biz, ord.ID, order.OrderStatusFulfilled); err != nil {
						return err
					}
					fulfilledOrders++
				}
			}
			return nil
		}); err != nil {
			return err
		}

		fmt.Println("✅ Seed complete")
		fmt.Println("---")
		fmt.Println("Workspace:", ws.ID)
		fmt.Println("Business:", biz.ID, "("+bizDescriptor+")")
		fmt.Println("Owner email:", ownerEmail)
		fmt.Println("Owner password:", seedPassword)
		fmt.Printf("Categories: %d\n", len(categories))
		fmt.Printf("Products: %d\n", counts.products)
		fmt.Printf("Customers: %d\n", len(customers))
		fmt.Printf("Orders: %d (paid %d, fulfilled %d)\n", createdOrders, paidOrders, fulfilledOrders)
		return nil
	},
}

func init() {
	seedCmd.Flags().BoolVar(&seedClean, "clean", false, "truncate all tables before seeding")
	seedCmd.Flags().StringVar(&seedSizeFlag, "size", string(seedSmall), "seed size: small|medium|large")
	seedCmd.Flags().StringVar(&seedPassword, "password", "KyoraDev@123", "password for seeded owner account")
	rootCmd.AddCommand(seedCmd)
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
