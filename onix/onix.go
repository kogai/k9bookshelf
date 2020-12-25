package onix

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/kogai/k9bookshelf/gqlgenc/client"
)

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

	for _, d := range data.Products {
		isbn := d.Productidentifiers.FindByIDType("ISBN-13")
		if isbn == nil {
			return fmt.Errorf("%s not have ISBN-13 value", d.Title.TitleText)
		}
		found, idx := hasSameISBN13(*isbn, products)
		if found {
			fmt.Println("Update", d.Title.TitleText, d.Title.Subtitle)

			currentProduct := products.Products.Edges[idx]
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

			date, err := extractDatetime(d.PublicationDate)
			if err != nil {
				return err
			}
			numberOfPages := fmt.Sprintf("%d", d.NumberOfPages)
			value := date.String()
			valueType := client.MetafieldValueTypeString
			res, err := gqlClient.ProductCreateDo(context.Background(), client.ProductInput{
				CollectionsToJoin: []string{"gid://shopify/Collection/236195152071"}, // NOTE: /collections/recommend
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
