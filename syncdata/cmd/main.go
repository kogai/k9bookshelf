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
		err := syncdata.Deploy(cmd.Flag("input").Value.String())
		if err != nil {
			log.Fatal(err)
		}
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download contents from store",
	Run: func(cmd *cobra.Command, args []string) {
		err := syncdata.Download(cmd.Flag("output").Value.String())
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
	deployCmd.PersistentFlags().StringP("input", "i", fmt.Sprintf("%s", cwd), "input directory")
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(deployCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
