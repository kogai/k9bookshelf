package content

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

var htmlToMakrdownTestcases = []struct {
	in  string
	out string
}{
	{
		`<p>abc</p>`,
		`abc
`,
	},
	{
		`<ul><li>abc</li></ul><p>def</p>`,
		`- abc

def
`,
	}, {
		`<ul><li>abc</li></ul><b>def</b>`,
		`- abc

 **def**
`}, {
		`<ul>
	<li>
	<a href="https://help.shopify.com/en/manual/apps/app-types#private-apps">Private App</a> でショップへのアクセス権限を取得</li>
	<li><a href="https://shopify.dev/docs/admin-api">Admin API</a></li>
	<li>
	<a href="https://shopify.dev/docs/admin-api/graphql/reference/common-objects/queryroot/index">GraphQL API</a>
	<ul>
	<li>GraphQLのクエリファイルから以下のツールでクライアントを生成</li>
	<li><a href="https://github.com/Yamashou/gqlgenc">https://github.com/Yamashou/gqlgenc</a></li>
	<li><a href="https://github.com/99designs/gqlgen">https://github.com/99designs/gqlgen</a></li>
	</ul>
	</li>
	<li>
	<a href="https://shopify.dev/docs/admin-api/rest/reference">REST API</a>
	<ul>
	<li>ページとブログの更新はGraphQL APIでサポートされていないのでREST APIも併用</li>
	</ul>
	</li>
	</ul>`,
		`- [Private App](https://help.shopify.com/en/manual/apps/app-types#private-apps) でショップへのアクセス権限を取得
- [Admin API](https://shopify.dev/docs/admin-api)
- [GraphQL API](https://shopify.dev/docs/admin-api/graphql/reference/common-objects/queryroot/index)
    - GraphQLのクエリファイルから以下のツールでクライアントを生成
    - [https://github.com/Yamashou/gqlgenc](https://github.com/Yamashou/gqlgenc)
    - [https://github.com/99designs/gqlgen](https://github.com/99designs/gqlgen)
- [REST API](https://shopify.dev/docs/admin-api/rest/reference)
    - ページとブログの更新はGraphQL APIでサポートされていないのでREST APIも併用
`},
	{
		`<p>Taking some of the most popular, bestselling recent games, Schreier immerses readers in the hellfire of the development process, whether it’s RPG studio Bioware’s challenge to beat an impossible schedule and overcome countless technical nightmares to build <em>Dragon Age: Inquisition</em>; indie developer Eric Barone’s single-handed efforts to grow country-life RPG <em>Stardew Valley</em> from one man’s vision into a multi-million-dollar franchise; or Bungie spinning out from their corporate overlords at Microsoft to create <em>Destiny</em>, a brand new universe that they hoped would become as iconic as <em>Star Wars</em> and <em>Lord of the Rings</em> –even as it nearly ripped their studio apart.</p>`,
		`Taking some of the most popular, bestselling recent games, Schreier immerses readers in the hellfire of the development process, whether it’s RPG studio Bioware’s challenge to beat an impossible schedule and overcome countless technical nightmares to build _Dragon Age: Inquisition_; indie developer Eric Barone’s single-handed efforts to grow country-life RPG _Stardew Valley_ from one man’s vision into a multi-million-dollar franchise; or Bungie spinning out from their corporate overlords at Microsoft to create _Destiny_, a brand new universe that they hoped would become as iconic as _Star Wars_ and _Lord of the Rings_ –even as it nearly ripped their studio apart.
`,
	},
}

func TestHtmlToMarkdown(t *testing.T) {
	t.Parallel()
	for _, tt := range htmlToMakrdownTestcases {
		t.Run(tt.in, func(t *testing.T) {
			actual, err := htmlToMarkdown(tt.in)
			assert.Equal(t, nil, err)
			assert.Equal(t, tt.out, actual)
		})
	}
}
