package onix

import (
	"encoding/xml"
)

// Header is not documented yet.
type Header struct {
	FromCompany           string `xml:"m174"`
	SentDate              string `xml:"m182"`
	DefaultLanguageOfText string `xml:"m184"`
	DefaultPriceTypeCode  string `xml:"m185"`
	DefaultCurrencyCode   string `xml:"m186"`
}

// Productidentifier is not documented yet.
type Productidentifier struct {
	ProductIDType string `xml:"b221"`
	IDValue       string `xml:"b244"`
}

// SupplyDetail is not documented yet.
type SupplyDetail struct {
	SupplierName        string `xml:"j137"`
	SupplierRole        string `xml:"j292"`
	ReturnsCodeType     string `xml:"j268"`
	ReturnsCode         string `xml:"j269"`
	ProductAvailability string `xml:"j396"`
	PackQuantity        string `xml:"j145"`
	Price               []struct {
		PriceTypeCode string `xml:"j148"`
		DiscountCoded []struct {
			DiscountCodeType     string `xml:"j363"`
			DiscountCodeTypeName string `xml:"j378"`
			DiscountCode         string `xml:"j364"`
		} `xml:"discountcoded"`
		PriceAmount  float32 `xml:"j151"`
		CurrencyCode string  `xml:"j152"`
		CountryCode  string  `xml:"b251"`
	} `xml:"price"`
}

// Product is not documented yet.
type Product struct {
	RecordReference       string              `xml:"a001"`
	NotificationType      string              `xml:"a002"`
	Productidentifier     []Productidentifier `xml:"productidentifier"`
	ProductForm           string              `xml:"b012"`
	ProductFormDetail     []string            `xml:"b333"`
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
		Subtitle  string `xml:"b029"`
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
	Subject []struct {
		SubjectSchemeIdentifier int    `xml:"b067"`
		SubjectSchemeName       string `xml:"b171"`
		SubjectCode             string `xml:"b069,omitempty"`
		SubjectHeadingText      string `xml:"b070"`
	} `xml:"subject"`
	AudienceCode string `xml:"b073"`
	OtherText    []struct {
		TextTypeCode string `xml:"d102"`
		Text         struct {
			Body       string `xml:",cdata"`
			TextFormat string `xml:"textformat,attr"`
		} `xml:"d104"`
	} `xml:"othertext"`
	Imprint []struct {
		NameCodeType     string `xml:"b241"`
		NameCodeTypeName string `xml:"b242"`
		NameCodeValue    string `xml:"b243"`
		ImprintName      string `xml:"b079"`
	} `xml:"imprint"`
	Publisher struct {
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
	Measure []struct {
		MeasureTypeCode string `xml:"c093"`
		Measurement     string `xml:"c094"`
		MeasureUnitCode string `xml:"c095"`
	} `xml:"measure"`
	RelatedProduct []struct {
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
