package content

const apiVersion string = "2020-10"

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
