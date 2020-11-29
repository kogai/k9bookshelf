package main

import (
	"log"

	"github.com/kogai/k9bookshelf/onix"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "onix-kit",
	Short: "onix-kit imports ONIX for Books 2.1 file to Shopify",
	Run: func(cmd *cobra.Command, args []string) {
		input := cmd.Flag("input").Value.String()
		if input == "" {
			log.Fatalln("[input] should be passed.")
		}
		if err := onix.Run(input); err != nil {
			log.Fatalln(err)
		}
	},
}

func main() {
	rootCmd.PersistentFlags().StringP("input", "i", "", "input ONIX for Books 2.1 file")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
