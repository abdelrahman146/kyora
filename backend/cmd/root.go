package cmd

import (
	"os"

	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kyora",
	Short: "Kyora is a simple social media business management saas application",
	Long: `Kyora is a simple social media business management saas application
	that helps businesses manage their orders, customers, expenses, and analyze their business performance.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load config once for all commands. Missing config file is OK (env vars / flags may be used).
		if err := config.Load(); err != nil {
			return err
		}
		// Initialize logger after config is loaded.
		logger.Init()
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kyora.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
