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

func findMetaFieldIDBy(fetchedProducts *client.ProductISBNs, idx int, key string) *string {
	currentProduct := fetchedProducts.Products.Edges[idx]
	for _, edge := range currentProduct.Node.Metafields.Edges {
		if edge.Node.Key == key {
			return &edge.Node.ID
		}
	}
	return nil
}

func metaFieldInput(onixProduct *Product, fetchedProducts *client.ProductISBNs, idx int) ([]*client.MetafieldInput, error) {
	numberOfPages := fmt.Sprintf("%d", onixProduct.NumberOfPages)
	date, err := extractDatetime(onixProduct.PublicationDate)
	if err != nil {
		return nil, err
	}

	value := date.String()
	valueType := client.MetafieldValueTypeString
	publishedAtID := findMetaFieldIDBy(fetchedProducts, idx, metaFieldKeyPublishedAt)
	subTitleID := findMetaFieldIDBy(fetchedProducts, idx, metaFieldKeySubTitle)
	numberOfPagesID := findMetaFieldIDBy(fetchedProducts, idx, metaFieldKeyNumberOfPages)

	return []*client.MetafieldInput{{
		ID:        publishedAtID,
		Key:       &metaFieldKeyPublishedAt,
		Namespace: &metaFieldNamespace,
		Value:     &value,
		ValueType: &valueType,
	}, {
		ID:        subTitleID,
		Value:     &onixProduct.Title.Subtitle,
		Key:       &metaFieldKeySubTitle,
		Namespace: &metaFieldNamespace,
		ValueType: &valueType,
	}, {
		ID:        numberOfPagesID,
		Value:     &numberOfPages,
		Key:       &metaFieldKeyNumberOfPages,
		Namespace: &metaFieldNamespace,
		ValueType: &valueType,
	}}, nil
}

func updateInput(onixProduct *Product, fetchedProducts *client.ProductISBNs, idx int) (*client.ProductInput, error) {
	currentProduct := fetchedProducts.Products.Edges[idx]
	title := onixProduct.Title.TitleText
	mtIpt, err := metaFieldInput(onixProduct, fetchedProducts, idx)
	if err != nil {
		return nil, err
	}

	// NOTE: DescriptionHTML and Tags are possible to edit manually,
	// So we should touch only at create-time.
	return &client.ProductInput{
		ID:         &currentProduct.Node.ID,
		Metafields: mtIpt,
		Title:      &title,
		// NOTE: DescriptionHTML and Tags are possible to edit manually,
		// So we should touch only at create-time.
	}, nil
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

	for _, d := range data.Products {
		isbn := d.Productidentifiers.FindByIDType("ISBN-13")
		if isbn == nil {
			return fmt.Errorf("%s not have ISBN-13 value", d.Title.TitleText)
		}
		found, idx := hasSameISBN13(*isbn, products)
		if found {
			fmt.Println("Update", d.Title.TitleText, d.Title.Subtitle)
			ipt, err := updateInput(&d, products, idx)
			if err != nil {
				return err
			}
			res, err := gqlClient.ProductUpdateDo(context.Background(), *ipt)
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
			descriptionHTML, err := generateDescription(&d)
			if err != nil {
				return err
			}
			mtIpt, err := metaFieldInput(&d, products, idx)
			if err != nil {
				return err
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

			res, err := gqlClient.ProductCreateDo(context.Background(), client.ProductInput{
				CollectionsToJoin: []string{"gid://shopify/Collection/236195152071"}, // NOTE: /collections/recommend
				DescriptionHTML:   descriptionHTML,
				Metafields:        mtIpt,
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
