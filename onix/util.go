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
	codegen "github.com/kogai/onix-codegen/generated/go/v2"
)

// TODO: Convert implementation to instance method.
func findPriceBy(prices []codegen.Price, currencyCode string) *codegen.Price {
	var price *codegen.Price
	for _, p := range prices {
		if p.CurrencyCode.Body == currencyCode {
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

func extractTags(p *codegen.Product) []string {
	var tags []string
	sbs := Subjects(p.Subjects)
	subject := sbs.FindByIDType("Keywords")
	if subject != nil && subject.SubjectHeadingText != nil {
		tags = strings.Split(*subject.SubjectHeadingText, "; ")
	}
	return tags
}

func extractDatetime(date string) (time.Time, error) {
	const shortForm = "20060102"
	return time.Parse(shortForm, date)
}

func generateDescription(onixProduct *codegen.Product) (*string, error) {
	var descriptionHTML string = ""
	otherTexts := OtherTexts(onixProduct.OtherTexts)
	otherText := otherTexts.FindByType("Long description")
	if otherText == nil {
		otherText = otherTexts.FindByType("Main description")
	}
	if otherText == nil {
		otherText = otherTexts.FindByType("Short description/annotation")
	}
	if otherText == nil {
		otherText = otherTexts.FindByType("Biographical note")
	}
	translated, err := Translate(string(*otherText.Text))
	if err != nil {
		return nil, err
	}
	if otherText != nil {
		descriptionHTML = fmt.Sprintf(`<h2>出版社より</h2>
%s
<hr/>
<h2>DeepL 粗訳</h2>
%s`, *otherText.Text, *translated)
	}
	return &descriptionHTML, nil
}
