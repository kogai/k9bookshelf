package syncdata

const apiVersion string = "2020-10"

// const shopDomain string = "k9books.myshopify.com"
// var appKey string = os.Getenv("MARKDOWN_APP_KEY")
// var appSecret string = os.Getenv("MARKDOWN_APP_SECRET")
// var shopToken string = appSecret

// Content is not documented yet.
type Content struct {
	handle string
	html   string
}

// Contents is not documented yet.
type Contents struct {
	kind  string
	items []Content
}
