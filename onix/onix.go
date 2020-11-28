package onix

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
)

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
