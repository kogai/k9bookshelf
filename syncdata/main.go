package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	file, err := ioutil.ReadFile("./syncdata/product.gql")
	if err != nil {
		panic(err)
	}
	body, err := json.Marshal(map[string]interface{}{
		"query": fmt.Sprintf("%s", file),
		"variables": map[string]int{
			"first": 10,
		},
	})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://%s/admin/api/%s/graphql.json", "k9books.myshopify.com", "2020-10"),
		bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	req.Header.Add("X-Shopify-Access-Token", os.Getenv("MARKDOWN_APP_SECRET"))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Response", string(body))
}
