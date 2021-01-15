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
		dryRun := cmd.Flag("dryRun").Value.String()
		if input == "" {
			log.Fatalln("[input] should be passed.")
		}
		if err := onix.Run(input, dryRun == "true"); err != nil {
			log.Fatalln(err)
		}
	},
}

func main() {
	rootCmd.PersistentFlags().StringP("input", "i", "", "input ONIX for Books 2.1 file")
	rootCmd.PersistentFlags().BoolP("dryRun", "d", false, "dry run")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
