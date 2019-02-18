package cmd

import (
	"os"

	"github.com/Oliver-Chang/ddns/util"

	"github.com/Oliver-Chang/ddns/util/logger"

	"github.com/Oliver-Chang/ddns/dns"
	"github.com/fsnotify/fsnotify"
	"github.com/robfig/cron"

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
		var (
			config *dns.DDNSConfig
			err    error
		)
		config = &dns.DDNSConfig{}
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			var err error
			for {
				err = viper.Unmarshal(config)
				if err != nil {
					logger.Logger.WithError(err).Error()
					continue
				}
				break
			}
		})
		err = viper.Unmarshal(config)
		if err != nil {
			logger.Logger.WithError(err).Error()
		}
		logger.Logger.WithField("config", config).Info()
		c := cron.New()
		ipChan := make(chan string, 1)
		var storeIP *string
		err = c.AddFunc("@every 5s", func() {
		redo:
			ip, err := util.GetIP()
			if err != nil {
				logger.Logger.WithError(err).Error("Get ip err")
			}
			logger.Logger.WithField("ipv6", ip).Info()
			if !util.IsIPv6(ip) {
				goto redo
			}

			if storeIP == nil || *storeIP != ip {
				storeIP = &ip
				ipChan <- *storeIP
			}
		})
		if err != nil {
			logger.Logger.Error(err)
			return
		}
		c.Start()
		for {
			select {
			case ip, ok := <-ipChan:
				if ok {
					logger.Logger.WithField("ipv6", ip).Info("Process")
					ddns := dns.NewDDNS(config)
					ddns.CreateRecord(ip, config.Domain)
				}
			}
		}

	},
}

// Execute Execute
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Logger.WithError(err).Error()
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
			logger.Logger.WithError(err).Error()
			os.Exit(1)
		}

		// Search config in home directory with name ".ddns" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ddns")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Logger.Info("Using config file:", viper.ConfigFileUsed())
	}
}
