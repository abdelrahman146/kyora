package cmd

import (
	"log"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/expense"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// recurringCmd represents the recurring expenses job command
var recurringCmd = &cobra.Command{
	Use:   "recurring",
	Short: "Run recurring expenses daily job",
	Run:   runRecurring,
}

func init() {
	rootCmd.AddCommand(recurringCmd)
}

func runRecurring(cmd *cobra.Command, args []string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
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
	// Auto-migrate expense tables to ensure columns exist
	if err := postgres.AutoMigrate(expense.Expense{}, expense.RecurringExpense{}); err != nil {
		log.Fatal("failed to migrate tables", utils.Log.Err(err))
	}
	atomicProcess := db.NewAtomicProcess(postgres.DB())

	storeDomain := store.NewDomain(postgres, atomicProcess, nil)
	expenseDomain := expense.NewDomain(postgres, atomicProcess, nil, storeDomain)

	if err := expenseDomain.ExpenseService.ProcessRecurringExpensesDaily(cmd.Context()); err != nil {
		log.Fatal("recurring job failed", utils.Log.Err(err))
	}
}
