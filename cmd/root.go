package cmd

import (
	"os"

	"github.com/Oliver-Chang/ddns/utils/logger"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"

	"github.com/Oliver-Chang/ddns/daemon"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{

	Use:   "ddns",
	Short: "DDNS",
	Long:  `DDNS`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cfg := daemon.Config{}
		viper.WatchConfig()

		if err := viper.Unmarshal(&cfg); err != nil {
			logger.Logger.Error("viper config unmarshal failed", zap.NamedError("config", err))
			return
		}
		// logger.Logger.Info(fmt.Sprintf("%+v", cfg))
		d := daemon.New(&cfg)
		// cfgReflesh := make(chan bool)
		viper.OnConfigChange(func(in fsnotify.Event) {
			viper.Unmarshal(&cfg)
			logger.Logger.Info("cfg file have new change")
			if err := d.FetchIPv6(); err != nil {
				logger.Logger.Error("fetch ipv6 address failed", zap.NamedError("fetch_error", err))
			}
		})
		d.Daemon()
	},
}

// Execute Execute
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// logger.Logger.WithError(err).Error()
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ddns.yaml)")
	rootCmd.PersistentFlags().StringP("uid", "u", "", "For cloudflare is AuthID, dnspod is token_id")
	rootCmd.PersistentFlags().StringP("token", "t", "", "For cloudflare is AuthKey, dnspod is token")
	rootCmd.PersistentFlags().StringP("zone-id", "z", "", "For cloudflare is ZoneID, dnspod is DomainID")
	rootCmd.PersistentFlags().StringP("class", "c", "", "type of dns")
	rootCmd.PersistentFlags().StringP("domain", "d", "", "domain")
	// // rootCmd.MarkPersistentFlagRequired("uid")
	// // rootCmd.MarkPersistentFlagRequired("token")
	// // rootCmd.MarkPersistentFlagRequired("zone-id-idin")
	// // rootCmd.MarkPersistentFlagRequired("class")
	// // rootCmd.MarkPersistentFlagRequired("domain")
	viper.BindPFlag("uid", rootCmd.PersistentFlags().Lookup("uid"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("zone_id", rootCmd.PersistentFlags().Lookup("zone-id"))
	viper.BindPFlag("dns", rootCmd.PersistentFlags().Lookup("class"))
	viper.BindPFlag("domain", rootCmd.PersistentFlags().Lookup("domain"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			// logger.Logger.WithError(err).Error()
			os.Exit(1)
		}

		// Search config in home directory with name ".ddns" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ddns")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Logger.Info("Using config file:", zap.String("config_file", viper.ConfigFileUsed()))
	}
}
