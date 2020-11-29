package main

import (
	"fmt"
	"log"
	"os"

	"k9bookshelf/syncdata"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "data-kit",
	Short: "data-kit is a content management tool like theme-kit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Nothing to do without subcommand.")
	},
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Upload contents to store",
	Run: func(cmd *cobra.Command, args []string) {
		input := cmd.Flag("input").Value.String()
		shopDomain := cmd.Flag("domain").Value.String()
		appKey := cmd.Flag("key").Value.String()
		appSecret := cmd.Flag("secret").Value.String()
		shopToken := cmd.Flag("token").Value.String()
		if shopDomain == "" || appKey == "" || appSecret == "" || shopToken == "" {
			log.Fatal(fmt.Sprintf("One of required parameter is empty, shopDomain='%s' appKey='%s' appSecret='%s' shopToken='%s'", shopDomain, appKey, appSecret, shopToken))
		}
		err := syncdata.Deploy(shopDomain, appKey, appSecret, shopToken, input)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download contents from store",
	Run: func(cmd *cobra.Command, args []string) {
		output := cmd.Flag("output").Value.String()
		shopDomain := cmd.Flag("domain").Value.String()
		appKey := cmd.Flag("key").Value.String()
		appSecret := cmd.Flag("secret").Value.String()
		shopToken := cmd.Flag("token").Value.String()
		if shopDomain == "" || appKey == "" || appSecret == "" || shopToken == "" {
			log.Fatal(fmt.Sprintf("One of required parameter is empty, shopDomain='%s' appKey='%s' appSecret='%s' shopToken='%s'", shopDomain, appKey, appSecret, shopToken))
		}
		err := syncdata.Download(shopDomain, appKey, appSecret, shopToken, output)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	downloadCmd.PersistentFlags().StringP("output", "o", fmt.Sprintf("%s", cwd), "output directory")
	downloadCmd.PersistentFlags().String("domain", "", "ShopDomain of your shop ex:your-shop.myshopify.com")
	downloadCmd.PersistentFlags().String("key", "", "Key of Admin API")
	downloadCmd.PersistentFlags().String("secret", "", "Secret of Admin API")
	downloadCmd.PersistentFlags().String("token", "", "AccessToken for Admin API generally same as secret if using Private App.")

	deployCmd.PersistentFlags().StringP("input", "i", fmt.Sprintf("%s", cwd), "input directory")
	deployCmd.PersistentFlags().String("domain", "", "ShopDomain of your shop ex:your-shop.myshopify.com")
	deployCmd.PersistentFlags().String("key", "", "Key of Admin API")
	deployCmd.PersistentFlags().String("secret", "", "Secret of Admin API")
	deployCmd.PersistentFlags().String("token", "", "AccessToken for Admin API generally same as secret if using Private App.")

	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(deployCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
