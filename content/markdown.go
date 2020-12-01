package content

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
)

func htmlToMarkdown(html string) (string, error) {
	opt := &md.Options{
		HorizontalRule: "---",
	}

	converter := md.NewConverter("", true, opt)
	markdownStr, err := converter.ConvertString(html)
	if err != nil {
		return "", err
	}
	return markdownStr + "\n", nil
}
