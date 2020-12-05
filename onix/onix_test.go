package onix

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/kogai/k9bookshelf/gqlgenc/client"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindPriceBy(t *testing.T) {
	t.Parallel()
	price := findPriceBy([]Price{{
		PriceAmount:  1.0,
		CurrencyCode: "USD",
	}, {
		PriceAmount:  2.0,
		CurrencyCode: "JPY",
	}}, "USD")
	assert.Equal(t, price.PriceAmount, 1.0)
}

func TestFindPriceByNotFound(t *testing.T) {
	t.Parallel()
	price := findPriceBy([]Price{}, "USD")
	assert.Equal(t, price, nil)
}

func TestHasSameISBN13(t *testing.T) {
	t.Parallel()
	fixture, err := ioutil.ReadFile("./fixtures/TestHasSameISBN13.json")
	assert.Equal(t, nil, err)
	var products client.ProductISBNs
	err = json.Unmarshal([]byte(fixture), &products)
	assert.Equal(t, nil, err)
	found, idx := hasSameISBN13("9781839214110", &products)
	assert.Equal(t, found, true)
	assert.Equal(t, products.Products.Edges[idx].Node.Title, "Node.js Design Patterns - Third edition: Design and implement production-grade Node.js applications using proven patterns and techniques")
}

func TestExtractID(t *testing.T) {
	t.Parallel()
	id, err := extractID("gid://shopify/Product/6080250872007")
	assert.Equal(t, nil, err)
	assert.Equal(t, id, 6080250872007)
}
