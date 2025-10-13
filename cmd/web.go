/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/asset"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/expense"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/domain/owner"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/domain/supplier"
	"github.com/abdelrahman146/kyora/internal/utils"
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
	if err != nil {
		log.Fatal("failed to connect to database", utils.Log.Err(err))
	}
	cache := db.NewMemcache(viper.GetStringSlice("db.memcache_hosts"))
	atomicProcess := db.NewAtomicProcess(postgres.DB())

	storeDomain := store.NewDomain(postgres, atomicProcess, cache)
	storeProvisioner := store.NewOnboardingStoreProvisioner(storeDomain.StoreService)
	accountDomain := account.NewDomain(postgres, atomicProcess, cache, storeProvisioner)
	assetDomain := asset.NewDomain(postgres, atomicProcess, cache, storeDomain.StoreService)
	customerDomain := customer.NewDomain(postgres, atomicProcess, cache, storeDomain.StoreService)
	expenseDomain := expense.NewDomain(postgres, atomicProcess, cache, storeDomain.StoreService)
	inventoryDomain := inventory.NewDomain(postgres, atomicProcess, cache, storeDomain.StoreService)
	orderDomain := order.NewDomain(postgres, atomicProcess, cache, storeDomain.StoreService)
	ownerDomain := owner.NewDomain(postgres, atomicProcess, cache, storeDomain.StoreService)
	supplierDomain := supplier.NewDomain(postgres, atomicProcess, cache, storeDomain.StoreService)

	router := webrouter.NewRouter()
	webrouter.RegisterRoutes(router,
		accountDomain,
		storeDomain,
		inventoryDomain,
		orderDomain,
		customerDomain,
		ownerDomain,
		assetDomain,
		expenseDomain,
		supplierDomain,
	)

	err = router.Run(viper.GetString("server.port"))
	if err != nil {
		log.Fatal("failed to start server", utils.Log.Err(err))
	}
}
