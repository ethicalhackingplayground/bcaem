package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bbscope",
	Short: "Grab AEM Programs from Bugcrowd",
	Long:  `The ultimate aem scope gathering tool for Bugcrowd`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bcaem.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Global flags
	rootCmd.PersistentFlags().StringP("proxy", "", "", "HTTP Proxy (Useful for debugging. Example: http://127.0.0.1:8080)")
	rootCmd.PersistentFlags().StringP("delimiter", "d", " ", "Delimiter character used when printing multiple data using the output flag")
	rootCmd.PersistentFlags().BoolP("bbpOnly", "b", false, "Only fetch aem programs offering monetary rewards")
	rootCmd.PersistentFlags().BoolP("pvtOnly", "p", false, "Only fetch aem programs from private programs")
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

		// Search config in home directory with name ".bbscope" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mcaem")
	}

	viper.AutomaticEnv() // read in environment variables that match
}
