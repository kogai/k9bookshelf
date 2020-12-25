package onix

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	shopify "github.com/bold-commerce/go-shopify"
	"github.com/kogai/k9bookshelf/gqlgenc/client"
)

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

func hasSameISBN13(isbn string, products *client.ProductISBNs) (bool, int) {
	fmt.Println(isbn)
	for i, p := range products.Products.Edges {
		for _, v := range p.Node.Variants.Edges {
			if isbn == *v.Node.Barcode {
				return true, i
			}
		}
	}
	return false, 0
}

func establishRestClient(shopDomain, appKey, appSecret string) *shopify.Client {
	app := shopify.App{
		ApiKey:    appKey,
		ApiSecret: appSecret,
	}

	return shopify.NewClient(app, shopDomain, appSecret, shopify.WithVersion(apiVersion))
}

func extractID(graphqlID string) (int, error) {
	result := regexp.MustCompile("gid://shopify/Product/(.*)").FindStringSubmatch(graphqlID)
	return strconv.Atoi(result[len(result)-1])
}

func extractTags(p *Product) []string {
	var tags []string
	subject := p.Subjects.FindByIDType("Keywords")
	if subject != nil {
		tags = strings.Split(subject.SubjectHeadingText, "; ")
	}
	return tags
}

func findMetaFieldIDBy(_p interface{}, key string) (*string, error) {
	p, ok := _p.(*client.Product)
	if !ok {
		return nil, fmt.Errorf("invalid type casting, got =[%v]", _p)
	}
	for _, edge := range p.Metafields.Edges {
		if edge.Node.Key == key {
			return &edge.Node.ID, nil
		}
	}
	return nil, nil
}

func extractDatetime(date string) (time.Time, error) {
	const shortForm = "20060102"
	return time.Parse(shortForm, date)
}
