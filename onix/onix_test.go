package onix

import (
	"testing"

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
