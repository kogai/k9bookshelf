package content

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

const frontMatterDelimiter = "---"

// FrontMatter is optional metadata written at the head of a markdown file.
// It is used to create blog articles which don't exist on the store yet.
type FrontMatter struct {
	Title       string   `yaml:"title,omitempty"`
	PublishedAt string   `yaml:"published_at,omitempty"`
	Tags        []string `yaml:"tags,omitempty"`
}

func (f *FrontMatter) publishedAt() (*time.Time, error) {
	if f == nil || f.PublishedAt == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, f.PublishedAt)
	if err != nil {
		return nil, fmt.Errorf("published_at must be RFC3339 format: %w", err)
	}
	return &t, nil
}

func (f *FrontMatter) tags() string {
	if f == nil {
		return ""
	}
	return strings.Join(f.Tags, ", ")
}

// splitFrontMatter separates YAML front matter from markdown body.
// When raw doesn't start with front matter, it returns nil and raw as is.
func splitFrontMatter(raw string) (*FrontMatter, string, error) {
	normalized := strings.ReplaceAll(raw, "\r\n", "\n")
	if !strings.HasPrefix(normalized, frontMatterDelimiter+"\n") {
		return nil, raw, nil
	}
	rest := normalized[len(frontMatterDelimiter)+1:]
	end := strings.Index(rest, "\n"+frontMatterDelimiter+"\n")
	if end < 0 {
		return nil, "", fmt.Errorf("front matter is not closed")
	}
	var fm FrontMatter
	if err := yaml.UnmarshalStrict([]byte(rest[:end]), &fm); err != nil {
		return nil, "", err
	}
	body := strings.TrimLeft(rest[end+len(frontMatterDelimiter)+2:], "\n")
	return &fm, body, nil
}

func renderFrontMatter(fm *FrontMatter) (string, error) {
	if fm == nil {
		return "", nil
	}
	b, err := yaml.Marshal(fm)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	buf.WriteString(frontMatterDelimiter + "\n")
	buf.Write(b)
	buf.WriteString(frontMatterDelimiter + "\n\n")
	return buf.String(), nil
}

func articleToFrontMatter(a Article) *FrontMatter {
	fm := FrontMatter{Title: a.Title}
	if a.PublishedAt != nil {
		fm.PublishedAt = a.PublishedAt.Format(time.RFC3339)
	}
	for _, t := range strings.Split(a.Tags, ",") {
		if t = strings.TrimSpace(t); t != "" {
			fm.Tags = append(fm.Tags, t)
		}
	}
	return &fm
}
