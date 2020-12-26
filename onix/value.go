package onix

import "os"

const apiVersion string = "2020-10"
const shopDomain string = "k9books.myshopify.com"
const fixedExchangeRate float64 = 110

var metaFieldNamespace string = "k9bookshelf"
var metaFieldKeyPublishedAt string = "published_at"
var metaFieldKeySubTitle string = "subtitle"
var metaFieldKeyNumberOfPages string = "number_of_pages"

var appKey string = os.Getenv("INGRAM_CONTENT_IMPORTER_APP_KEY")
var appSecret string = os.Getenv("INGRAM_CONTENT_IMPORTER_APP_SECRET")
