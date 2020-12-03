package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/kogai/k9bookshelf/content"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var rootCmd = &cobra.Command{
	Use:   "content-kit",
	Short: "content-kit is a content management tool like theme-kit which is theme management tool",
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatalln("Nothing to do without subcommand.")
	},
}

type configTy struct {
	Content struct {
		Domain *string `yaml:"domain,omitempty"`
		Key    *string `yaml:"key,omitempty"`
		Secret *string `yaml:"secret,omitempty"`
		Token  *string `yaml:"token,omitempty"`
		Dir    *string `yaml:"dir,omitempty"`
	} `yaml:"content,omitempty"`
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Upload contents to store",
	Run: func(cmd *cobra.Command, args []string) {
		var input, shopDomain, appKey, appSecret, shopToken string

		config := cmd.Flag("config").Value.String()
		if config != "" {
			buf, err := ioutil.ReadFile(config)
			if err != nil {
				log.Fatalln(err)
			}
			expanded := os.ExpandEnv(string(buf))
			var conf configTy
			err = yaml.Unmarshal([]byte(expanded), &conf)
			if err != nil {
				log.Fatalln(err)
			}

			if conf.Content.Domain != nil {
				shopDomain = *conf.Content.Domain
			} else {
				shopDomain = cmd.Flag("domain").Value.String()
			}
			if conf.Content.Key != nil {
				appKey = *conf.Content.Key
			} else {
				appKey = cmd.Flag("key").Value.String()
			}
			if conf.Content.Secret != nil {
				appSecret = *conf.Content.Secret
			} else {
				appSecret = cmd.Flag("secret").Value.String()
			}
			if conf.Content.Token != nil {
				shopToken = *conf.Content.Token
			} else {
				shopToken = cmd.Flag("token").Value.String()
			}
			if conf.Content.Dir != nil {
				input = *conf.Content.Dir
			} else {
				input = cmd.Flag("dir").Value.String()
			}
		} else {
			shopDomain = cmd.Flag("domain").Value.String()
			appKey = cmd.Flag("key").Value.String()
			appSecret = cmd.Flag("secret").Value.String()
			shopToken = cmd.Flag("token").Value.String()
			input = cmd.Flag("dir").Value.String()
		}

		if shopDomain == "" || appKey == "" || appSecret == "" || shopToken == "" {
			log.Fatalln(fmt.Sprintf("One of required parameter is empty, shopDomain='%s' appKey='%s' appSecret='%s' shopToken='%s'", shopDomain, appKey, appSecret, shopToken))
		}
		err := content.Deploy(shopDomain, appKey, appSecret, shopToken, input)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download contents from store",
	Run: func(cmd *cobra.Command, args []string) {
		var output, shopDomain, appKey, appSecret, shopToken string

		config := cmd.Flag("config").Value.String()
		if config != "" {
			buf, err := ioutil.ReadFile(config)
			if err != nil {
				log.Fatalln(err)
			}
			expanded := os.ExpandEnv(string(buf))
			var conf configTy
			err = yaml.Unmarshal([]byte(expanded), &conf)
			if err != nil {
				log.Fatalln(err)
			}

			if conf.Content.Domain != nil {
				shopDomain = *conf.Content.Domain
			} else {
				shopDomain = cmd.Flag("domain").Value.String()
			}
			if conf.Content.Key != nil {
				appKey = *conf.Content.Key
			} else {
				appKey = cmd.Flag("key").Value.String()
			}
			if conf.Content.Secret != nil {
				appSecret = *conf.Content.Secret
			} else {
				appSecret = cmd.Flag("secret").Value.String()
			}
			if conf.Content.Token != nil {
				shopToken = *conf.Content.Token
			} else {
				shopToken = cmd.Flag("token").Value.String()
			}
			if conf.Content.Dir != nil {
				output = *conf.Content.Dir
			} else {
				output = cmd.Flag("dir").Value.String()
			}
		} else {
			shopDomain = cmd.Flag("domain").Value.String()
			appKey = cmd.Flag("key").Value.String()
			appSecret = cmd.Flag("secret").Value.String()
			shopToken = cmd.Flag("token").Value.String()
			output = cmd.Flag("dir").Value.String()
		}

		if shopDomain == "" || appKey == "" || appSecret == "" || shopToken == "" {
			log.Fatalln(fmt.Sprintf("One of required parameter is empty, shopDomain='%s' appKey='%s' appSecret='%s' shopToken='%s'", shopDomain, appKey, appSecret, shopToken))
		}
		err := content.Download(shopDomain, appKey, appSecret, shopToken, output)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	downloadCmd.PersistentFlags().StringP("dir", "d", fmt.Sprintf("%s", cwd), "directory where contents exists")
	downloadCmd.PersistentFlags().String("domain", "", "ShopDomain of your shop ex:your-shop.myshopify.com")
	downloadCmd.PersistentFlags().String("key", "", "Key of Admin API")
	downloadCmd.PersistentFlags().String("secret", "", "Secret of Admin API")
	downloadCmd.PersistentFlags().String("token", "", "AccessToken for Admin API generally same as secret if using Private App.")
	downloadCmd.PersistentFlags().String("config", "", "configuration file which includes api key and so on")

	deployCmd.PersistentFlags().StringP("dir", "d", fmt.Sprintf("%s", cwd), "directory where contents exists")
	deployCmd.PersistentFlags().String("domain", "", "ShopDomain of your shop ex:your-shop.myshopify.com")
	deployCmd.PersistentFlags().String("key", "", "Key of Admin API")
	deployCmd.PersistentFlags().String("secret", "", "Secret of Admin API")
	deployCmd.PersistentFlags().String("token", "", "AccessToken for Admin API generally same as secret if using Private App.")
	deployCmd.PersistentFlags().String("config", "", "configuration file which includes api key and so on")

	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(deployCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
