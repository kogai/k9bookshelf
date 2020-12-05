package onix

import (
	"encoding/xml"
	"fmt"
)

// Header is not documented yet.
type Header struct {
	FromCompany           string `xml:"m174"`
	SentDate              string `xml:"m182"`
	DefaultLanguageOfText string `xml:"m184"`
	DefaultPriceTypeCode  string `xml:"m185"`
	DefaultCurrencyCode   string `xml:"m186"`
}

// ProductIDType is not documented yet.
type ProductIDType string

// UnmarshalXML is not documented yet.
func (c *ProductIDType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	switch v {
	// TODO: Define as enum
	case "01":
		*c = "Proprietary"
	case "02":
		*c = "ISBN-10"

	// NOTE: GS1 Global Trade Item Number, formerly known as EAN article number (13 digits)
	case "03":
		*c = "GTIN-13"
	case "04":
		*c = "UPC"
	case "05":
		*c = "ISMN-10"
	case "06":
		*c = "DOI"
	case "13":
		*c = "LCCN"
	case "14":
		*c = "GTIN-14"
	case "15":
		*c = "ISBN-13"
	case "17":
		*c = "Legal deposit number"
	case "22":
		*c = "URN"
	case "23":
		*c = "OCLC number"
	case "24":
		*c = "Co-publisher’s ISBN-13"
	case "25":
		*c = "ISMN-13"
	case "26":
		*c = "ISBN-A"
	case "27":
		*c = "JP e-code"
	case "28":
		*c = "OLCC number"
	case "29":
		*c = "JP Magazine ID"
	case "30":
		*c = "UPC12+5"
	case "31":
		*c = "BNF Control number"
	case "35":
		*c = "ARK"
	default:
		return fmt.Errorf("undefined code has been passed, got [%s]", v)
	}
	return nil
}

// Productidentifier is not documented yet.
type Productidentifier struct {
	ProductIDType ProductIDType `xml:"b221"`
	IDValue       string        `xml:"b244"` // EAN
}

// Productidentifiers is not documented yet.
type Productidentifiers []Productidentifier

// FindByIDType findx identifier by id-type.
func (c *Productidentifiers) FindByIDType(idType ProductIDType) *string {
	for _, p := range *c {
		if p.ProductIDType == idType {
			return &p.IDValue
		}
	}
	return nil
}

// DiscountCodeType is not documented yet.
type DiscountCodeType string

// UnmarshalXML is not documented yet.
func (c *DiscountCodeType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	switch v {
	// TODO: Define as enum
	case "01":
		*c = "BIC discount group code"
	case "02":
		*c = "Proprietary discount code"
	case "03":
		*c = "Boeksoort"
	case "04":
		*c = "German terms code"
	case "05":
		*c = "Proprietary commission code"
	case "06":
		*c = "BIC commission group code"
	default:
		return fmt.Errorf("undefined code has been passed, got [%s]", v)
	}
	return nil
}

// Price is not documented yet.
type Price struct {
	PriceTypeCode  string `xml:"j148"`
	DiscountCodeds []struct {
		DiscountCodeType     DiscountCodeType `xml:"j363"`
		DiscountCodeTypeName string           `xml:"j378"`
		DiscountCode         string           `xml:"j364"`
	} `xml:"discountcoded"`
	PriceAmount  float64 `xml:"j151"`
	CurrencyCode string  `xml:"j152"`
	CountryCode  string  `xml:"b251"`
}

// Prices is not documented yet.
type Prices []Price

// FindByType findx identifier by id-type.
func (c *Prices) FindByType(ty string) *Price {
	for _, p := range *c {
		if p.CurrencyCode == ty {
			return &p
		}
	}
	return nil
}

// SupplyDetail is not documented yet.
type SupplyDetail struct {
	SupplierName        string `xml:"j137"`
	SupplierRole        string `xml:"j292"`
	ReturnsCodeType     string `xml:"j268"`
	ReturnsCode         string `xml:"j269"`
	ProductAvailability string `xml:"j396"`
	PackQuantity        int    `xml:"j145"`
	Prices              Prices `xml:"price"`
}

// MeasureTypeCode is not documented yet.
type MeasureTypeCode string

// UnmarshalXML is not documented yet.
func (c *MeasureTypeCode) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	switch v {
	// TODO: Define as enum
	case "01":
		*c = "Height"
	case "02":
		*c = "Width"
	case "03":
		*c = "Thickness"
	case "04":
		*c = "Page trim height"
	case "05":
		*c = "Page trim width"
	case "08":
		*c = "Unit weight"
	case "09":
		*c = "Diameter (sphere)"
	case "10":
		*c = "Unfolded/unrolled sheet height"
	case "11":
		*c = "Unfolded/unrolled sheet width"
	case "12":
		*c = "Diameter (tube or cylinder)"
	case "13":
		*c = "Rolled sheet package side measure"
	default:
		return fmt.Errorf("undefined code has been passed, got [%s]", v)
	}
	return nil
}

// Measure is not documented yet.
type Measure struct {
	MeasureTypeCode MeasureTypeCode `xml:"c093"`
	Measurement     float64         `xml:"c094"`
	MeasureUnitCode string          `xml:"c095"`
}

const kirogramPerPound float64 = 0.4535924277
const gramPerOunce float64 = 28.349

// ToKg convert measure to kirogram.
func (m *Measure) ToKg() (float64, error) {
	switch m.MeasureUnitCode {
	case "gr":
		return m.Measurement / 1000, nil
	case "kg":
		return m.Measurement, nil
	case "lb":
		return m.Measurement * kirogramPerPound, nil
	case "oz":
		return m.Measurement * gramPerOunce / 1000, nil
	default:
		return 0, fmt.Errorf("Unexpected value of the unit was passed, got [%s] [%s]", m.MeasureTypeCode, m.MeasureUnitCode)
	}
}

// Measures is not documented yet.
type Measures []Measure

// FindByType findx identifier by id-type.
func (c *Measures) FindByType(ty MeasureTypeCode) *Measure {
	for _, p := range *c {
		if p.MeasureTypeCode == ty {
			return &p
		}
	}
	return nil
}

// TextTypeCode is not documented yet.
type TextTypeCode string

// UnmarshalXML is not documented yet.
func (c *TextTypeCode) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	switch v {
	// TODO: Define as enum
	case "01":
		*c = "Main description"
	case "02":
		*c = "Short description/annotation"
	case "03":
		*c = "Long description"
	case "04":
		*c = "Table of contents"
	case "05":
		*c = "Review quote, restricted length"
	case "06":
		*c = "Quote from review of previous edition"
	case "07":
		*c = "Review text"
	case "08":
		*c = "Review quote"
	case "09":
		*c = "Promotional ‘headline’"
	case "10":
		*c = "Previous review quote"
	case "11":
		*c = "Author comments"
	case "12":
		*c = "Description for reader"
	case "13":
		*c = "Biographical note"
	case "14":
		*c = "Description for Reading Group Guide"
	case "15":
		*c = "Discussion question for Reading Group Guide"
	case "16":
		*c = "Competing titles"
	case "17":
		*c = "Flap copy"
	case "18":
		*c = "Back cover copy"
	case "19":
		*c = "Feature"
	case "20":
		*c = "New feature"
	case "21":
		*c = "Publisher’s notice"
	case "22":
		*c = "Index"
	case "23":
		*c = "Excerpt from book"
	case "24":
		*c = "First chapter"
	case "25":
		*c = "Description for sales people"
	case "26":
		*c = "Description for press or other media"
	case "27":
		*c = "Description for subsidiary rights department"
	case "28":
		*c = "Description for teachers/educators"
	case "30":
		*c = "Unpublished endorsement"
	case "31":
		*c = "Description for bookstore"
	case "32":
		*c = "Description for library"
	case "33":
		*c = "Introduction or preface"
	case "34":
		*c = "Full text"
	case "35":
		*c = "Promotional text"
	case "40":
		*c = "Author interview / QandA"
	case "41":
		*c = "Reading Group Guide"
	case "42":
		*c = "Commentary / discussion"
	case "43":
		*c = "Short description for series or set"
	case "44":
		*c = "Long description for series or set"
	case "45":
		*c = "Contributor event schedule"
	case "46":
		*c = "License"
	case "47":
		*c = "Open access statement"
	case "48":
		*c = "Digital exclusivity statement"
	case "49":
		*c = "Official recommendation"
	case "98":
		*c = "Master brand name"
	case "99":
		*c = "Country of final manufacture"
	default:
		return fmt.Errorf("undefined code has been passed, got [%s]", v)
	}
	return nil
}

// TextFormat is not documented yet.
type TextFormat string

// UnmarshalXMLAttr is not documented yet.
func (c *TextFormat) UnmarshalXMLAttr(d xml.Attr) error {
	switch d.Value {
	// TODO: Define as enum
	case "00":
		*c = "ASCII text"
	case "01":
		*c = "SGML"
	case "02":
		*c = "HTML"
	case "03":
		*c = "XML"
	case "04":
		*c = "PDF"
	case "05":
		*c = "XHTML"
	case "06":
		*c = "Default text format"
	case "07":
		*c = "Basic ASCII text"
	case "08":
		*c = "PDF"
	case "09":
		*c = "Microsoft rich text format (RTF)"
	case "10":
		*c = "Microsoft Word binary format (DOC)"
	case "11":
		*c = "ECMA 376 WordprocessingML"
	case "12":
		*c = "ISO 26300 ODF"
	case "13":
		*c = "Corel Wordperfect binary format (DOC)"
	case "14":
		*c = "EPUB"
	case "15":
		*c = "XPS"
	default:
		return fmt.Errorf("undefined code has been passed, got [%s]", d.Value)
	}
	return nil
}

// Text is not documented yet.
type Text struct {
	Body       string     `xml:",cdata"`
	TextFormat TextFormat `xml:"textformat,attr"`
}

// OtherText is not documented yet.
type OtherText struct {
	TextTypeCode TextTypeCode `xml:"d102"`
	Text         Text         `xml:"d104"`
}

// OtherTexts is not documented yet.
type OtherTexts []OtherText

// FindByType findx identifier by id-type.
func (c *OtherTexts) FindByType(ty TextTypeCode) *OtherText {
	for _, p := range *c {
		if p.TextTypeCode == ty {
			return &p
		}
	}
	return nil
}

// SubjectSchemeIdentifier is Main subject scheme identifier code.
type SubjectSchemeIdentifier string

// UnmarshalXML is not documented yet.
func (c *SubjectSchemeIdentifier) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	switch v {
	// TODO: Define as enum
	case "01":
		*c = "Dewey"
	case "02":
		*c = "Abridged Dewey"
	case "03":
		*c = "LC classification"
	case "04":
		*c = "LC subject heading"
	case "05":
		*c = "NLM classification"
	case "06":
		*c = "MeSH heading"
	case "07":
		*c = "NAL subject heading"
	case "08":
		*c = "AAT"
	case "09":
		*c = "UDC"
	case "10":
		*c = "BISAC Subject Heading"
	case "11":
		*c = "BISAC region code"
	case "12":
		*c = "BIC subject category"
	case "13":
		*c = "BIC geographical qualifier"
	case "14":
		*c = "BIC language qualifier (language as subject)"
	case "15":
		*c = "BIC time period qualifier"
	case "16":
		*c = "BIC educational purpose qualifier"
	case "17":
		*c = "BIC reading level and special interest qualifier"
	case "18":
		*c = "DDC-Sachgruppen der Deutschen Nationalbibliografie"
	case "19":
		*c = "LC fiction genre heading"
	case "20":
		*c = "Keywords"
	case "21":
		*c = "BIC children’s book marketing category"
	case "22":
		*c = "BISAC Merchandising Theme"
	case "23":
		*c = "Publisher’s own category code"
	case "24":
		*c = "Proprietary subject scheme"
	case "25":
		*c = "Tabla de materias ISBN"
	case "26":
		*c = "Warengruppen-Systematik des deutschen Buchhandels"
	case "27":
		*c = "SWD"
	case "93":
		*c = "Thema subject category"
	default:
		return fmt.Errorf("undefined code has been passed, got [%s]. See ONIX_for_Books_Release2-1_rev03_docs+codes_Issue_36/codelists/onix-codelist-26.htm", v)
	}
	return nil
}

// Subject is not documented yet.
type Subject struct {
	SubjectSchemeIdentifier SubjectSchemeIdentifier `xml:"b067"`
	SubjectSchemeName       string                  `xml:"b171,omitempty"`
	SubjectCode             string                  `xml:"b069,omitempty"`
	SubjectHeadingText      string                  `xml:"b070,omitempty"`
}

// Subjects is not documented yet.
type Subjects []Subject

// FindByIDType findx identifier by id-type.
func (c *Subjects) FindByIDType(idType SubjectSchemeIdentifier) *Subject {
	for _, p := range *c {
		if p.SubjectSchemeIdentifier == idType {
			return &p
		}
	}
	return nil
}

// Imprint is not documented yet.
type Imprint struct {
	NameCodeType     string `xml:"b241"`
	NameCodeTypeName string `xml:"b242"`
	NameCodeValue    string `xml:"b243"`
	ImprintName      string `xml:"b079"`
}

// Imprints is not documented yet.
type Imprints = []Imprint

// Product is not documented yet.
type Product struct {
	RecordReference       string             `xml:"a001"`
	NotificationType      string             `xml:"a002"`
	Productidentifiers    Productidentifiers `xml:"productidentifier"`
	ProductForm           string             `xml:"b012"`
	ProductFormDetail     []string           `xml:"b333"`
	ProductClassification struct {
		ProductClassificationType string `xml:"b274"`
		ProductClassificationCode string `xml:"b275"`
	} `xml:"productclassification"`
	NoSeries BoolIfElementPresent `xml:"n338"`
	Title    struct {
		TextCase  string `xml:"textcase,attr"`
		Language  string `xml:"language,attr"`
		TitleType string `xml:"b202"`
		TitleText string `xml:"b203"`
		Subtitle  string `xml:"b029,omitempty"`
	} `xml:"title"`
	Contributor struct {
		ContributorRole  string `xml:"b035"`
		NamesBeforeKey   string `xml:"b039"`
		KeyNames         string `xml:"b040"`
		BiographicalNote string `xml:"b044"`
	} `xml:"contributor"`
	Language struct {
		LanguageRole string `xml:"b253"`
		LanguageCode string `xml:"b252"`
	} `xml:"language"`
	NumberOfPages    int    `xml:"b061"`
	BASICMainSubject string `xml:"b064"`
	MainSubject      struct {
		MainSubjectSchemeIdentifier int    `xml:"b191"`
		SubjectCode                 string `xml:"b069"`
		SubjectHeadingText          string `xml:"b070"`
	} `xml:"mainsubject"`
	Subjects     Subjects   `xml:"subject"`
	AudienceCode string     `xml:"b073"`
	OtherTexts   OtherTexts `xml:"othertext"`
	Imprints     Imprints   `xml:"imprint"`
	Publisher    struct {
		PublishingRole string `xml:"b291"`
		PublisherName  string `xml:"b081"`
	} `xml:"publisher"`
	PublishingStatus struct {
		Body      string `xml:",innerxml"`
		Datestamp string `xml:"datestamp,attr"`
	} `xml:"b394"`
	PublicationDate string `xml:"b003"`
	SalesRights     []struct {
		SalesRightsType string `xml:"b089"`
		RightsCountry   string `xml:"b090"`
		RightsTerritory string `xml:"b388"`
	} `xml:"salesrights"`
	Measures        Measures `xml:"measure"`
	RelatedProducts []struct {
		RelationCode      string              `xml:"h208"`
		Productidentifier []Productidentifier `xml:"productidentifier"`
		ProductForm       string              `xml:"b012"`
	} `xml:"relatedproduct"`
	SupplyDetail SupplyDetail `xml:"supplydetail"`
}

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

// IngramContentOnix is not documented yet.
type IngramContentOnix struct {
	XMLName  xml.Name  `xml:"ONIXmessage"`
	Header   Header    `xml:"header"`
	Products []Product `xml:"product"`
}

// ToUploadable converts raw XML to Shopify specified data structure.
// TODO: It may needs some *.gql file and client generator.
func (c *IngramContentOnix) ToUploadable() error {
	return nil
}
