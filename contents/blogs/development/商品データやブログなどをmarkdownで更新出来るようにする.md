こんにちは。
趣味でShopify[マーチャントをやっている](https://k9bookshelf.com/)ものです。

この記事は [Shopify開発を盛り上げる（Liquid, React, Node.js, Graph QL） Advent Calendar 2020](https://qiita.com/advent-calendar/2020/shopify-liquid) の11日目の記事です。

昨日は[benzookapi](https://qiita.com/benzookapi)さんの[placeholder]でした。

皆さんShopifyで商品情報やページの編集をする時はどうされていますか？
当然ですが、ショップの管理画面で提供されているエディタ(以下リッチエディタ)で更新されていることと思います。
(ちなみにShopifyのデザインシステム[Poralisで提供して欲しいという要望](https://github.com/Shopify/polaris-react/issues/303#issuecomment-415554317)が多く集まっているようですね)

私は元々仕事でShopifyアプリの開発をしていた延長で、趣味と実益を兼ねてこのサイトを出店しているのですが、マーチャントの立場で改めて管理画面を使ってみると、どうにもこのエディタでは満足できないことに気づいたのです。

と言うのも普段公私ともに文書を書くのに最もよく使っているのは[markdown](https://www.markdownguide.org/)だからです。

また本業のプログラマーという職業柄、文章を書くには使い慣れたテキストエディタを使いたいという気持ちもあります。

Shopifyではテーマ開発用ツールとして[theme-kit](https://github.com/Shopify/themekit)が提供されており、
Shopifyで編集したテーマファイルをダウンロードしたり手元のエディタで編集したものをアップロードすることが出来ます。

これに近い体験が、コンテンツの管理でも出来ないものでしょうか？

コミュニティを[検索してみると](https://community.shopify.com/c/forums/searchpage/tab/message?advanced=false&allow_punctuation=false&filter=location&location=category:en&q=markdown)、markdownで編集出来るようなアプリをリリースしたというコメントは散見されるものの、デッドリンクとなっているなど、2020年時点でも使える解決策は見つかりませんでした。

幸いShopifyには[Admin API](https://shopify.dev/docs/admin-api)を始めとする、多種多様なAPIが公開されています。
また、リッチエディタで保存された文言は、内部的にはHTMLで保存されているようです。

商品情報編集やページ・ブログ記事の更新までAPI経由で操作が可能なので、theme-kitよろしくcontent-kitのようなツールが書けそうです。
(と言うか、この記事は[content-kitと名付けたその自作ツール](https://github.com/kogai/k9bookshelf/blob/main/content/README.md)で更新しています)

## ツールの構成

開発に際して参考にしたドキュメントやツールは以下の通りです。

- [Private App](https://help.shopify.com/en/manual/apps/app-types#private-apps)でショップへのアクセス権限を取得
- [Admin API](https://shopify.dev/docs/admin-api)
    - [GraphQL API](https://shopify.dev/docs/admin-api/graphql/reference/common-objects/queryroot/index)
        - GraphQLのクエリファイルから以下のツールでクライアントを生成
            - https://github.com/Yamashou/gqlgenc
            - https://github.com/99designs/gqlgen
    - [REST API](https://shopify.dev/docs/admin-api/rest/reference)
        - ページとブログの更新はGraphQL APIでサポートされていないのでzREST APIも併用

Shopifyはドキュメントが非常に充実していて、[GraphiQLアプリ](https://shopify.dev/tools/graphiql-admin-api)によるインタラクティブな試行環境もあるので、
特に詰まることもなく開発出来ました。

ページやブログはREST APIの[Goクライアントライブラリ](https://github.com/bold-commerce/go-shopify)でサポートされていませんでした。
この辺りはShopifyアプリではあまり取り扱われないリソースなのかも知れません。

唯一困ったのは、テーマファイルや他のツール(テーマや書誌情報の規格化されたデータファイルの取り込みツールなど)を同じレポジトリで管理している都合上、
[bazel](https://bazel.build/)でビルドをしているのですが、GraphQLクライアントの生成ツールが依存している `*.gotpl` などのファイルの依存関係が自動生成できなかったことです。

Shopify関係ないですね。

(bazelのgo系ルールには依存解決時にパッチを当てる仕組みがあるようで、[パッチを書くことで](https://github.com/kogai/k9bookshelf/blob/b7bb804c0ad45b5eed5215d1b62a9c434c4cc6aa/content/com_github_yamashou_gqlgenc.patch#L1-L25)解決出来ました)

まだ`theme-kit`に比べると、watchがないなど基本的な機能がまだ足りてないのですが、一応使い始めることは出来ています。

## デモ

実際の使用感はこんな感じです。

[](gif)

バイナリを[Releaseページ]()に置いておくので、良かったら試してみて下さい。

```
content -h
content-kit is a content management tool like theme-kit which is theme management tool

Usage:
  content-kit [flags]
  content-kit [command]

Available Commands:
  deploy      Upload contents to store
  download    Download contents from store
  help        Help about any command

Flags:
  -h, --help   help for content-kit

Use "content-kit [command] --help" for more information about a command.
```

次回は[ShoheiTai](https://qiita.com/ShoheiTai)さんの「Shopifyアプリの選定・運用ノウハウとか」です。
お楽しみに！
