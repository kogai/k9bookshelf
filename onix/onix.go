package onix

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
)

// Run imports ONIX for Books 2.1 format file to Shopify.
func Run() error {
	file, err := ioutil.ReadFile("/Users/kogaishinichi/k9books/onix/14362217.onix")
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
	// fmt.Println("data=", data)
	for _, d := range data.Products {
		fmt.Println("d=", d.PublishingStatus)
		// d.PublishingStatus
		// for _, t := range d.OtherText {
		// 	fmt.Println("t=", t.Text)
		// }
	}
	return nil
}
