package content

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestHtmlToMarkdown(t *testing.T) {
	t.Parallel()
	md, err := htmlToMarkdown(`<p>abc</p>`)
	assert.Equal(t, nil, err)
	assert.Equal(t, `abc
`, md)
}

func TestHtmlToMarkdownList(t *testing.T) {
	t.Parallel()
	md, err := htmlToMarkdown(`<ul><li>abc</li></ul><p>def</p>`)
	assert.Equal(t, nil, err)
	assert.Equal(t, `- abc

def
`, md)
}

func TestHtmlToMarkdownListAndBold(t *testing.T) {
	t.Parallel()
	md, err := htmlToMarkdown(`<ul><li>abc</li></ul><b>def</b>`)
	assert.Equal(t, nil, err)
	assert.Equal(t, `- abc

 **def**
`, md)
}

func TestHtmlToMarkdownNestedList(t *testing.T) {
	t.Parallel()
	md, err := htmlToMarkdown(`<ul>
	<li><a href="https://help.shopify.com/en/manual/apps/app-types#private-apps">Private App</a> でショップへのアクセス権限を取得</li>
	<li><a href="https://shopify.dev/docs/admin-api">Admin API</a></li>
	<li><a href="https://shopify.dev/docs/admin-api/graphql/reference/common-objects/queryroot/index">GraphQL API</a>
		<ul>
			<li>GraphQLのクエリファイルから以下のツールでクライアントを生成</li>
			<li><a href="https://github.com/Yamashou/gqlgenc">https://github.com/Yamashou/gqlgenc</a></li>
			<li><a href="https://github.com/99designs/gqlgen">https://github.com/99designs/gqlgen</a></li>
		</ul>
	</li>
	<li><a href="https://shopify.dev/docs/admin-api/rest/reference">REST API</a>
		<ul>
			<li>ページとブログの更新はGraphQL APIでサポートされていないのでREST APIも併用</li>
		</ul>
	</li>
</ul>`)
	assert.Equal(t, nil, err)
	assert.Equal(t, `- [Private App](https://help.shopify.com/en/manual/apps/app-types#private-apps) でショップへのアクセス権限を取得
- [Admin API](https://shopify.dev/docs/admin-api)
- [GraphQL API](https://shopify.dev/docs/admin-api/graphql/reference/common-objects/queryroot/index)
    - GraphQLのクエリファイルから以下のツールでクライアントを生成
    - [https://github.com/Yamashou/gqlgenc](https://github.com/Yamashou/gqlgenc)
    - [https://github.com/99designs/gqlgen](https://github.com/99designs/gqlgen)
- [REST API](https://shopify.dev/docs/admin-api/rest/reference)
    - ページとブログの更新はGraphQL APIでサポートされていないのでREST APIも併用
`, md)
}
