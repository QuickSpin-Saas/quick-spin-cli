package cmd

import (
	"fmt"
	"os"

	"github.com/quickspin/quickspin-cli/internal/cmd/auth"
	"github.com/quickspin/quickspin-cli/internal/cmd/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Version information (set by main)
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"

	// Global flags
	cfgFile    string
	profile    string
	output     string
	noColor    bool
	verbose    bool
	debug      bool
	apiURL     string
	org        string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qspin",
	Short: "QuickSpin CLI - Managed microservices platform",
	Long: `QuickSpin CLI (qspin) is the official command-line interface for QuickSpin.

QuickSpin is a managed microservices SaaS platform providing Redis, RabbitMQ,
Elasticsearch, PostgreSQL, MongoDB, and other services through a shared
Kubernetes multi-tenant architecture.

Perfect for developers with limited RAM or Docker configuration challenges.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return err
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	// Add subcommands
	rootCmd.AddCommand(auth.NewAuthCmd())
	rootCmd.AddCommand(config.NewConfigCmd())
	rootCmd.AddCommand(NewVersionCmd())

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.quickspin/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "named profile to use")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output format: table, json, yaml (default: table)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable color output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug mode")
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "override API URL")
	rootCmd.PersistentFlags().StringVar(&org, "org", "", "override organization context")

	// Bind flags to viper
	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	_ = viper.BindPFlag("api.url", rootCmd.PersistentFlags().Lookup("api-url"))
	_ = viper.BindPFlag("defaults.organization", rootCmd.PersistentFlags().Lookup("org"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		// Search config in ~/.quickspin directory
		configDir := home + "/.quickspin"
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// Environment variable prefix
	viper.SetEnvPrefix("QUICKSPIN")
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		if debug || verbose {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}

	// Apply profile if specified
	if profile != "" {
		applyProfile(profile)
	}
}

func setDefaults() {
	viper.SetDefault("api.url", "https://api.quickspin.dev")
	viper.SetDefault("api.timeout", "30s")
	viper.SetDefault("auth.method", "jwt")
	viper.SetDefault("defaults.region", "us-east-1")
	viper.SetDefault("defaults.output", "table")
	viper.SetDefault("defaults.service_type", "redis")
	viper.SetDefault("defaults.tier", "developer")
	viper.SetDefault("telemetry.enabled", true)
	viper.SetDefault("telemetry.anonymous", true)
}

func applyProfile(profileName string) {
	profileKey := "profiles." + profileName
	if !viper.IsSet(profileKey) {
		fmt.Fprintf(os.Stderr, "Warning: Profile '%s' not found in config\n", profileName)
		return
	}

	// Apply profile settings
	if viper.IsSet(profileKey + ".api.url") {
		viper.Set("api.url", viper.GetString(profileKey+".api.url"))
	}
	if viper.IsSet(profileKey + ".api.timeout") {
		viper.Set("api.timeout", viper.GetString(profileKey+".api.timeout"))
	}
	if viper.IsSet(profileKey + ".defaults") {
		defaults := viper.GetStringMap(profileKey + ".defaults")
		for k, v := range defaults {
			viper.Set("defaults."+k, v)
		}
	}
}
