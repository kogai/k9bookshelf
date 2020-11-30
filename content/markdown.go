package content

import (
	"bytes"
	"strings"

	"github.com/mattn/godown"
)

// type CustomRule interface {
// 	Rule(next WalkFunc) (tagName string, customRule WalkFunc)
// }

func htmlToMarkdown(html string) (string, error) {
	var buf bytes.Buffer
	err := godown.Convert(&buf, strings.NewReader(html), nil)
	if err != nil {
		return "", err
	}
	s := buf.String()
	return strings.TrimRight(s, "\n") + "\n", nil
}
