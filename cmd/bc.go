package cmd

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"

	"github.com/ethicalhackingplayground/bcaem/pkg/bugcrowd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// bcCmd represents the bc command
var bcCmd = &cobra.Command{
	Use:   "bc",
	Short: "Bugcrowd",
	Long:  "Gathers AEM data from Bugcrowd (https://bugcrowd.com/)",
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := cmd.Flags().GetString("token")
		categories, _ := cmd.Flags().GetString("categories")
		concurrency, _ := cmd.Flags().GetInt("concurrency")

		delimiterCharacter, _ := rootCmd.PersistentFlags().GetString("delimiter")
		proxy, _ := rootCmd.PersistentFlags().GetString("proxy")
		bbpOnly, _ := rootCmd.Flags().GetBool("bbpOnly")
		pvtOnly, _ := rootCmd.Flags().GetBool("pvtOnly")

		email := viper.GetViper().GetString("bugcrowd-email")
		password := viper.GetViper().GetString("bugcrowd-password")

		if proxy != "" {
			proxyURL, err := url.Parse(proxy)
			if err != nil {
				log.Fatal("Invalid Proxy String")
			}
			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			http.DefaultTransport.(*http.Transport).Proxy = http.ProxyURL(proxyURL)
		}

		if email != "" && password != "" && token == "" {
			token = bugcrowd.Login(email, password)
		}

		bugcrowd.PrintAllScope(token, bbpOnly, pvtOnly, categories, delimiterCharacter, concurrency)
	},
}

func init() {
	rootCmd.AddCommand(bcCmd)
	bcCmd.Flags().StringP("token", "t", "", "Bugcrowd session token (_crowdcontrol_session cookie)")
	bcCmd.Flags().IntP("concurrency", "", 2, "Concurrency")

	bcCmd.Flags().StringP("email", "E", "", "Login email")
	viper.BindPFlag("bugcrowd-email", bcCmd.Flags().Lookup("email"))

	bcCmd.Flags().StringP("password", "P", "", "Login password")
	viper.BindPFlag("bugcrowd-password", bcCmd.Flags().Lookup("password"))

}
