package content

const apiVersion string = "2020-10"

// Content is not documented yet.
type Content struct {
	handle string
	html   string
	// meta is front matter of the markdown. It is nil when the file has none.
	meta *FrontMatter
}

// Contents is not documented yet.
type Contents struct {
	kind  string
	items []Content
}
