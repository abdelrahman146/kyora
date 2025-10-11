/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/abdelrahman146/kyora/internal/web/handlers"
	"github.com/abdelrahman146/kyora/internal/web/webrouter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runWeb,
}

func init() {
	rootCmd.AddCommand(webCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// webCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// webCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runWeb(cmd *cobra.Command, args []string) {
	viper.SetConfigName("config") // name of config file (without extension
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	utils.Log.Setup(viper.GetString("log.level"))
	postgres, err := db.NewPostgres(viper.GetString("db.dsn"), &gorm.Config{
		Logger: db.NewSlogGormLogger(gormlogger.Warn),
	})
	postgres.AutoMigrate(account.User{}, account.Organization{})
	if err != nil {
		log.Fatal("failed to connect to database", utils.Log.Err(err))
	}
	_ = db.NewMemcache(viper.GetStringSlice("db.memcache_hosts"))
	atomicProcess := db.NewAtomicProcess(postgres.DB())

	organizationRepo := account.NewOrganizationRepository(postgres)
	userRepo := account.NewUserRepository(postgres)
	storeRepo := store.NewStoreRepository(postgres)
	productRepo := inventory.NewProductRepository(postgres)
	variantRepo := inventory.NewVariantRepository(postgres)

	_ = account.NewAuthenticationService(userRepo)
	_ = account.NewOnboardingService(userRepo, organizationRepo, atomicProcess)
	_ = account.NewOrganizationService(organizationRepo)
	_ = account.NewUserService(userRepo)
	storeService := store.NewStoreService(storeRepo)
	_ = inventory.NewInventoryService(productRepo, variantRepo, storeService, atomicProcess)

	router := webrouter.NewRouter()
	dashboardHandler := handlers.NewDashboardHandler()
	dashboardHandler.RegisterRoutes(router)
	loginHandler := handlers.NewLoginHandler()
	loginHandler.RegisterRoutes(router)
	registerHandler := handlers.NewRegisterHandler()
	registerHandler.RegisterRoutes(router)
	forgotPasswordHandler := handlers.NewForgotPasswordHandler()
	forgotPasswordHandler.RegisterRoutes(router)
	resetPasswordHandler := handlers.NewResetPasswordHandler()
	resetPasswordHandler.RegisterRoutes(router)
	orderHandler := handlers.NewOrderHandler()
	orderHandler.RegisterRoutes(router)
	productHandler := handlers.NewProductHandler()
	productHandler.RegisterRoutes(router)
	customerHandler := handlers.NewCustomerHandler()
	customerHandler.RegisterRoutes(router)
	expenseHandler := handlers.NewExpenseHandler()
	expenseHandler.RegisterRoutes(router)
	supplierHandler := handlers.NewSupplierHandler()
	supplierHandler.RegisterRoutes(router)
	invoiceHandler := handlers.NewInvoiceHandler()
	invoiceHandler.RegisterRoutes(router)
	settingsHandler := handlers.NewSettingsHandler()
	settingsHandler.RegisterRoutes(router)
	analyticsHandler := handlers.NewAnalyticsHandler()
	analyticsHandler.RegisterRoutes(router)

	err = router.Run(viper.GetString("server.port"))
	if err != nil {
		log.Fatal("failed to start server", utils.Log.Err(err))
	}
}
