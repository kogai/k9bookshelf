package onix

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/kogai/k9bookshelf/gqlgenc/client"
)

const apiVersion string = "2020-10"
const shopDomain string = "k9books.myshopify.com"

var appSecret string = os.Getenv("MARKDOWN_APP_SECRET")

// TODO: Convert implementation to instance method.
func findPriceBy(prices []Price, currencyCode string) *Price {
	var price *Price
	for _, p := range prices {
		if p.CurrencyCode == currencyCode {
			price = &p
			break
		}
	}
	return price
}

func fetchProducts(ctx context.Context, adminClient *client.Client) (*client.ProductISBNs, error) {
	var cursor *string
	var res *client.ProductISBNs

	for {
		tmpRes, err := adminClient.ProductISBNs(ctx, 10, cursor)
		if err != nil {
			return nil, err
		}
		if res == nil {
			res = tmpRes
		} else {
			res.Products.Edges = append(res.Products.Edges, tmpRes.Products.Edges...)
		}

		if !tmpRes.Products.PageInfo.HasNextPage {
			break
		} else {
			last := tmpRes.Products.Edges[len(tmpRes.Products.Edges)-1]
			cursor = &last.Cursor
		}
	}
	return res, nil
}

// Run imports ONIX for Books 2.1 format file to Shopify.
func Run(input string) error {
	file, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	var data IngramContentOnix
	decoder := xml.NewDecoder(bytes.NewReader(file))
	decoder.CharsetReader = func(label string, input io.Reader) (io.Reader, error) {
		return input, nil
	}

	if err := decoder.Decode(&data); err != nil {
		return err
	}
	gqlClient, ctx := client.EstablishGqlClient(shopDomain, apiVersion, appSecret)
	products, err := fetchProducts(ctx, gqlClient)
	if err != nil {
		return err
	}

	fmt.Println(products.Products.Edges)
	fmt.Println("======")

	for _, d := range data.Products {
		fmt.Println("Title", d.Title.TitleText, d.Title.Subtitle)

		price := findPriceBy(d.SupplyDetail.Prices, "USD")
		fmt.Println("price", price.PriceAmount, price.DiscountCodeds)
		// d.PublishingStatus
		// for _, t := range d.OtherText {
		// 	fmt.Println("t=", t.Text)
		// }
		fmt.Println("======")
	}
	return nil
}
