package onix

import (
	"encoding/xml"
	"fmt"
	"strconv"

	codegen "github.com/kogai/onix-codegen/go"
)

// Productidentifiers is not documented yet.
type Productidentifiers []codegen.ProductIdentifier

// FindByIDType findx identifier by id-type.
func (c *Productidentifiers) FindByIDType(idType string) *string {
	for _, p := range *c {
		if p.ProductIDType.Body == idType {
			return &p.IDValue
		}
	}
	return nil
}

// Prices is not documented yet.
type Prices []codegen.Price

// FindByType findx identifier by id-type.
func (c *Prices) FindByType(ty string) *codegen.Price {
	for _, p := range *c {
		if p.CurrencyCode != nil && p.CurrencyCode.Body == ty {
			return &p
		}
	}
	return nil
}

const kirogramPerPound float64 = 0.4535924277
const gramPerOunce float64 = 28.349

// ToKg convert measure to kirogram.
func ToKg(m *codegen.Measure) (float64, error) {
	n, err := strconv.ParseFloat(m.Measurement, 64)
	if err != nil {
		return 0, err
	}
	switch m.MeasureUnitCode.Body {
	case "Grams":
		return n / 1000, nil
	case "Kilograms":
		return n, nil
	case "Pounds (US)":
		return n * kirogramPerPound, nil
	case "Ounces (US)":
		return n * gramPerOunce / 1000, nil
	default:
		return 0, fmt.Errorf("Unexpected value of the unit was passed, got [%s] [%s]", m.MeasureTypeCode, m.MeasureUnitCode)
	}
}

// Measures is not documented yet.
type Measures []codegen.Measure

// FindByType findx identifier by id-type.
func (c *Measures) FindByType(ty string) *codegen.Measure {
	for _, p := range *c {
		if p.MeasureTypeCode.Body == ty {
			return &p
		}
	}
	return nil
}

// OtherTexts is not documented yet.
type OtherTexts []codegen.OtherText

// FindByType findx identifier by id-type.
func (c *OtherTexts) FindByType(ty string) *codegen.OtherText {
	for _, p := range *c {
		if p.TextTypeCode.Body == ty {
			return &p
		}
	}
	return nil
}

// Subjects is not documented yet.
type Subjects []codegen.Subject

// FindByIDType findx identifier by id-type.
func (c *Subjects) FindByIDType(idType string) *codegen.Subject {
	for _, p := range *c {
		if p.SubjectSchemeIdentifier.Body == idType {
			return &p
		}
	}
	return nil
}

// Imprints is not documented yet.
type Imprints = []codegen.Imprint

// BoolIfElementPresent represent whether exsits self-closing tag.
type BoolIfElementPresent bool

// UnmarshalXML convert self-closing tag to bool.
// NOTE: https://stackoverflow.com/questions/23724591/golang-unmarshal-self-closing-tags
func (c *BoolIfElementPresent) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	*c = true
	return nil
}
