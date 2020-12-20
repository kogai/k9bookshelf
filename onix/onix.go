package onix

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	shopify "github.com/bold-commerce/go-shopify"
	"github.com/kogai/k9bookshelf/gqlgenc/client"
)

const apiVersion string = "2020-10"
const shopDomain string = "k9books.myshopify.com"
const fixedExchangeRate float64 = 110

var metaFieldNamespace string = "k9bookshelf"
var metaFieldKeyPublishedAt string = "published_at"
var metaFieldKeySubTitle string = "subtitle"
var metaFieldKeyNumberOfPages string = "number_of_pages"

var appKey string = os.Getenv("INGRAM_CONTENT_IMPORTER_APP_KEY")
var appSecret string = os.Getenv("INGRAM_CONTENT_IMPORTER_APP_SECRET")

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
	// restClient := establishRestClient(shopDomain, appKey, appSecret)
	products, err := fetchProducts(ctx, gqlClient)
	if err != nil {
		fmt.Println("HERE", shopDomain, apiVersion, appSecret)
		return err
	}

	for _, d := range data.Products {
		isbn := d.Productidentifiers.FindByIDType("ISBN-13")
		if isbn == nil {
			return fmt.Errorf("%s not have ISBN-13 value", d.Title.TitleText)
		}
		found, idx := hasSameISBN13(*isbn, products)
		if found {
			fmt.Println("Update", d.Title.TitleText, d.Title.Subtitle)

			currentProduct := products.Products.Edges[idx]
			// var descriptionHTML string = ""
			// otherText := d.OtherTexts.FindByType("Long description")
			// if otherText != nil {
			// 	descriptionHTML = otherText.Text.Body
			// }

			// var tags []string
			// subject := d.Subjects.FindByIDType("Keywords")
			// if subject != nil {
			// 	tags = strings.Split(*subject.SubjectHeadingText, "; ")
			// }
			// inventoryPolicy := client.ProductVariantInventoryPolicyContinue

			// var weight *float64
			// measure := d.Measures.FindByType("Unit weight")
			// weightUnit := client.WeightUnitKilograms
			// if measure != nil {
			// 	w, err := measure.ToKg()
			// 	if err != nil {
			// 		return err
			// 	}
			// 	weight = &w
			// }

			// var price *string
			// _price := d.SupplyDetail.Prices.FindByType("USD")
			// if _price != nil {
			// 	p := fmt.Sprintf("%f", _price.PriceAmount)
			// 	price = &p
			// }

			title := d.Title.TitleText
			numberOfPages := fmt.Sprintf("%d", d.NumberOfPages)
			date, err := extractDatetime(d.PublicationDate)
			if err != nil {
				return err
			}
			value := date.String()
			valueType := client.MetafieldValueTypeString
			var publishedAtID *string
			var subTitleID *string
			var numberOfPagesID *string
			for _, edge := range currentProduct.Node.Metafields.Edges {
				if edge.Node.Key == metaFieldKeyPublishedAt {
					publishedAtID = &edge.Node.ID
					break
				}
			}
			for _, edge := range currentProduct.Node.Metafields.Edges {
				if edge.Node.Key == metaFieldKeySubTitle {
					subTitleID = &edge.Node.ID
					break
				}
			}
			for _, edge := range currentProduct.Node.Metafields.Edges {
				if edge.Node.Key == metaFieldKeyNumberOfPages {
					numberOfPagesID = &edge.Node.ID
					break
				}
			}

			// NOTE: DescriptionHTML and Tags are possible to edit manually,
			// So we should touch only at create-time.
			res, err := gqlClient.ProductUpdateDo(context.Background(), client.ProductInput{
				ID: &currentProduct.Node.ID,
				Metafields: []*client.MetafieldInput{{
					ID:        publishedAtID,
					Key:       &metaFieldKeyPublishedAt,
					Namespace: &metaFieldNamespace,
					Value:     &value,
					ValueType: &valueType,
				}, {
					ID:        subTitleID,
					Value:     &d.Title.Subtitle,
					Key:       &metaFieldKeySubTitle,
					Namespace: &metaFieldNamespace,
					ValueType: &valueType,
				}, {
					ID:        numberOfPagesID,
					Value:     &numberOfPages,
					Key:       &metaFieldKeyNumberOfPages,
					Namespace: &metaFieldNamespace,
					ValueType: &valueType,
				}},
				// 	Variants: []*client.ProductVariantInput{
				// 		{
				// 			InventoryPolicy: &inventoryPolicy,
				// 			Weight:          weight,
				// 			WeightUnit:      &weightUnit,
				// 			Price:           price,
				// 			Barcode:         isbn,
				// 		},
				// 	},
				Title: &title,
			})
			if err != nil {
				return err
			}
			if len(res.ProductUpdate.UserErrors) > 0 {
				errMsg := ""
				for _, e := range res.ProductUpdate.UserErrors {
					errMsg += fmt.Sprintln(e.Field, ":", e.Message)
				}
				return fmt.Errorf(errMsg)
			}
		} else {
			fmt.Println("Create", d.Title.TitleText, d.Title.Subtitle)
			var descriptionHTML string = "<h2>出版社より</h2><br />"
			otherText := d.OtherTexts.FindByType("Long description")
			translated, err := Translate(otherText.Text.Body)
			if err != nil {
				return err
			}
			if otherText != nil {
				descriptionHTML = fmt.Sprintf(`<h2>出版社より</h2>
%s
<hr/>
<h2>DeepL 粗訳</h2>
%s`, otherText.Text.Body, *translated)
			}

			tags := extractTags(&d)
			title := d.Title.TitleText
			inventoryPolicy := client.ProductVariantInventoryPolicyContinue

			var weight *float64
			measure := d.Measures.FindByType("Unit weight")
			weightUnit := client.WeightUnitKilograms
			if measure != nil {
				w, err := measure.ToKg()
				if err != nil {
					return err
				}
				weight = &w
			}

			var price *string
			_price := d.SupplyDetail.Prices.FindByType("USD")
			if _price != nil {
				p := fmt.Sprintf("%f", _price.PriceAmount*fixedExchangeRate)
				price = &p
			}

			// onlineStoreID := "gid://shopify/Publication/68389830855"
			// facebookID := "gid://shopify/Publication/68395040967"
			// googleID := "gid://shopify/Publication/68395073735"
			// amazonID := "gid://shopify/Publication/68864508103"
			// chatID := "gid://shopify/Publication/68864934087"
			// buyButtonID := "gid://shopify/Publication/68977950919"
			// published := true
			date, err := extractDatetime(d.PublicationDate)
			if err != nil {
				return err
			}
			numberOfPages := fmt.Sprintf("%d", d.NumberOfPages)
			value := date.String()
			valueType := client.MetafieldValueTypeString
			res, err := gqlClient.ProductCreateDo(context.Background(), client.ProductInput{
				CollectionsToJoin: []string{"gid://shopify/Collection/236195152071"},
				DescriptionHTML:   &descriptionHTML,
				Metafields: []*client.MetafieldInput{{
					Key:       &metaFieldKeyPublishedAt,
					Namespace: &metaFieldNamespace,
					Value:     &value,
					ValueType: &valueType,
				}, {
					Value:     &d.Title.Subtitle,
					Key:       &metaFieldKeySubTitle,
					Namespace: &metaFieldNamespace,
					ValueType: &valueType,
				}, {
					Value:     &numberOfPages,
					Key:       &metaFieldKeyNumberOfPages,
					Namespace: &metaFieldNamespace,
					ValueType: &valueType,
				}},
				Variants: []*client.ProductVariantInput{
					{
						InventoryPolicy: &inventoryPolicy,
						Weight:          weight,
						WeightUnit:      &weightUnit,
						Price:           price,
						Barcode:         isbn,
					},
				},
				// Published: &published,
				// Publications: []*client.ProductPublicationInput{
				// 	{
				// 		PublicationID: &onlineStoreID,
				// 	},
				// 	{
				// 		PublicationID: &facebookID,
				// 	},
				// 	{
				// 		PublicationID: &googleID,
				// 	},
				// 	{
				// 		PublicationID: &amazonID,
				// 	},
				// 	{
				// 		PublicationID: &chatID,
				// 	},
				// 	{
				// 		PublicationID: &buyButtonID,
				// 	},
				// },
				Tags:   tags,
				Title:  &title,
				Vendor: &d.Publisher.PublisherName,
			})
			if err != nil {
				return err
			}
			if len(res.ProductCreate.UserErrors) > 0 {
				errMsg := ""
				for _, e := range res.ProductCreate.UserErrors {
					errMsg += fmt.Sprintln(e.Field, ":", e.Message)
				}
				return fmt.Errorf(errMsg)
			}

			fmt.Printf("Done. open 'https://ipage.ingramcontent.com/ipage/servlet/ibg.common.titledetail.imageloader?ean=%s&size=640&howerType=Y'\n", *isbn)
		}
	}
	return nil
}
