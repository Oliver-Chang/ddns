package cmd

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/robfig/cron"

	"github.com/Oliver-Chang/ddns/dns"

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
			config.UID = viper.GetString("uid")
			config.Token = viper.GetString("token")
			config.ZoneID = viper.GetString("zone_id")
			config.DNS = viper.GetString("dns")
			config.Domain = viper.GetString("domain")
		})
		config.UID = viper.GetString("uid")
		config.Token = viper.GetString("token")
		config.ZoneID = viper.GetString("zone_id")
		config.DNS = viper.GetString("dns")
		config.Domain = viper.GetString("domain")
		c := cron.New()
		err = c.AddFunc("@every 10m", func() {
			ddns := dns.NewDDNS(config)
			ddns.CreateRecord()
		})
		if err != nil {
			fmt.Println(err)
		}
		c.Start()
		fmt.Println(config.DNS)
		for {

		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
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
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ddns" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ddns")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
