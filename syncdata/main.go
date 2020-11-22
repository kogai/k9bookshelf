package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"k9bookshelf/generated"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Yamashou/gqlgenc/client"
	"github.com/gomarkdown/markdown"
	"github.com/mattn/godown"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "datakit",
	Short: "datakit is a content management tool like theme-kit",
	Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Println("Do stuff")
	},
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Upload contents to store",
	Run: func(cmd *cobra.Command, args []string) {
		deploy(cmd.Flag("input").Value.String())
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download contents from store",
	Run: func(cmd *cobra.Command, args []string) {
		download(cmd.Flag("output").Value.String())
	},
}

func gqlClient() (*generated.Client, context.Context) {
	authHeader := func(req *http.Request) {
		req.Header.Set("X-Shopify-Access-Token", os.Getenv("MARKDOWN_APP_SECRET"))
	}

	return &generated.Client{
		Client: client.NewClient(http.DefaultClient,
			fmt.Sprintf("https://%s/admin/api/%s/graphql.json", "k9books.myshopify.com", "2020-10"),
			authHeader),
	}, context.Background()
}

func download(output string) error {
	adminClient, ctx := gqlClient()
	res, err := adminClient.Products(ctx, 10)

	if err != nil {
		return err
	}

	// TODO: Use goroutine
	var downloadGroup sync.WaitGroup

	for _, edge := range res.Products.Edges {
		downloadGroup.Add(1)
		c := make(chan error)

		go func(handle, descriptionHTML string) {
			defer downloadGroup.Done()

			file, err := os.Create(path.Join(output, "products", handle+".md"))
			if err != nil {
				c <- err
				return
			}
			err = godown.Convert(file, strings.NewReader(descriptionHTML), nil)
			if err != nil {
				c <- err
				return
			}
			fmt.Printf("Done: %s.md\n", handle)
			c <- nil
		}(edge.Node.Handle, edge.Node.DescriptionHTML)
		err = <-c
		if err != nil {
			return err
		}
	}

	downloadGroup.Wait()
	return nil
}

func deploy(input string) error {
	files, err := ioutil.ReadDir(path.Join(input, "products"))
	if err != nil {
		return err
	}
	adminClient, ctx := gqlClient()
	wg := sync.WaitGroup{}
	for _, file := range files {
		wg.Add(1)
		c := make(chan error)
		filename := file.Name()

		go func(handle, pathToFile string) {
			defer wg.Done()
			productByHandle, err := adminClient.ProductByHandle(ctx, handle)
			if err != nil {
				c <- err
				return
			}

			md, err := ioutil.ReadFile(pathToFile)
			if err != nil {
				c <- err
				return
			}
			descriptionHTML := string(markdown.ToHTML(md, nil, nil))

			res, err := adminClient.Deploy(
				ctx,
				generated.ProductInput{
					ID:              &productByHandle.ProductByHandle.ID,
					Handle:          &handle,
					DescriptionHTML: &descriptionHTML,
				},
			)
			if err != nil {
				c <- err
				return
			}
			if len(res.ProductUpdate.UserErrors) > 0 {
				var errorBuf string
				for _, userError := range res.ProductUpdate.UserErrors {
					errorBuf += fmt.Sprintf("'%s': '%s'\n", userError.Field, userError.Message)
				}
				c <- fmt.Errorf("{\n%s}", errorBuf)
				return
			}
			c <- nil
		}(
			filename[0:len(filename)-len(filepath.Ext(filename))],
			path.Join(input, "products", filename),
		)

		err = <-c
		if err != nil {
			return err
		}
	}
	wg.Wait()
	return nil
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
