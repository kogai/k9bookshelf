# k9bookshelf

Shopify ストア [k9bookshelf](https://k9bookshelf.com) のテーマ、コンテンツ、サブツールのレポジトリ。

## 構成

- `contents/`: ストアに反映される Markdown
  - `products/<handle>.md`: 商品の説明文（本文のみ。front matter 不可）
  - `pages/<handle>.md`: ページ本文（本文のみ。front matter 不可）
  - `blogs/<blog-handle>/<article-handle>.md`: ブログ記事。先頭に YAML front matter（`title` / `published_at`（RFC3339）/ `tags`）を書ける。ストアにない記事を作成するには `title` が必須
  - `assets/`: 記事から raw.githubusercontent.com 経由で参照する画像。デプロイはされない
- `content/`: Content Kit。`contents/` と Shopify Admin API を同期する Go の CLI（`content/cmd/content`）
  - `content/testdata/rendered/`: `contents/` から生成した HTML の golden ファイル
- `theme/`: Shopify テーマ（Liquid）。`storefront/` の TypeScript を webpack で `theme/assets/app.min.js` にビルドする
- `onix/`: ONIX ファイルから商品をインポートするツール
- `gqlgenc/`: GraphQL Admin API クライアント。`gqlgenc/client/client.go` と `model.go` は生成コードなので手で編集しない（`gql.go` は手書き）
  - 生成にはストアのスキーマ（Admin API の token）が必要なため、手元ではなく GitHub Actions の `generate graphql client` を作業ブランチを指定して手動実行する。生成結果はそのブランチにコミットされる
  - そのコミットでは CI と Copilot のレビューが走らないので、呼び出し側の修正など後続のコミットを push して確認する

## 開発

```sh
make test         # gofmt / go vet / go test（コンテンツの検証と golden 比較を含む）
make golden       # contents/ を変えたら実行して golden ファイルを更新する
make check/theme  # テーマのビルドと lint（デプロイはしない）
```

- `contents/` や Markdown の変換処理を変えたら `make golden` を実行し、`content/testdata/rendered/` の差分が意図どおりか確認してからコミットする。golden の差分は PR で実際にストアへ反映される HTML の変化として見える
- 変更をコミットする前に `make test` を通す

## Pull Request

- PR はテーマ（目的）ごとに1本にする。例: コンテンツの追加・移動、ライブラリの更新、CI の修正、Content Kit の不具合修正は、それぞれ別の PR にする
- 作業中に別テーマの問題を見つけたら、今の PR に混ぜず、main から別ブランチを切って別の PR にする
- ライブラリの更新は、レビューしやすい単位（Go / Node.js / GitHub Actions など）で分ける
- 依存関係のある PR は、先にマージすべきものを PR の説明に書く

### レビュー

自分が作った PR は、マージ前に次の2つのレビューを受け、両方に応答する。指摘への応答は「対応した」「対応しない（理由を書く）」「別 PR で対応する（PR を示す）」のいずれか。別テーマの指摘は上記のとおり別 PR にする。

- **サブエージェントによるレビュー**: PR を作る前か作った直後に、差分をサブエージェントにレビューさせる（ファイルの変更や GitHub への投稿はさせない）。応答は PR の説明に記録する。指摘への対応で差分が大きく変わったら、再度レビューさせる
- **GitHub Copilot のレビュー**: PR を作ると自動で依頼される。指摘には各スレッドで返信し、返信したスレッドは resolve する。Copilot は承認ではなくコメントを返すので、「通る」ではなく「指摘に応答し終えた」ことを確認する
- push したら Copilot のレビューを再依頼する（GitHub MCP の `request_copilot_review` など）。依頼されていない、または最新コミットへのレビューが届かない場合も再依頼し、それでも届かなければマージせずにユーザーに伝える
- マージはユーザーの指示があったときだけ行う。その時点の最新コミットで、両方のレビューへの応答と CI が済んでいることを確認してからマージする

## 公開レポジトリ

このレポジトリは public。コミット、PR、Issue、コメントに書いたものは誰でも読める。

- 認証情報（API キー、アクセストークン、パスワード、秘密鍵、`.env`、サービスアカウントの JSON など）をコミットしない。値は GitHub Actions の secrets に置き、ファイルには `${ENV_NAME}` のように環境変数で参照する（`config.yml` 参照）
- 顧客情報、注文情報、売上など、ストアの非公開データをコミットしない。`onix/*.onix` は gitignore 済み
- PR やコミットメッセージにも、secret の値、管理画面の内部 URL、個人のメールアドレスを書かない
- `make test` に含まれる `make check/secrets` が、追跡中のファイルから secret らしい文字列や禁止ファイルを検出する。Claude Code では `git commit` の前に同じチェックが走り、検出するとコミットが止まる
- 誤ってコミットした場合は、履歴から消すだけでは不十分。push 済みならすぐにその secret を失効・再発行する

## デプロイ

- main へ push すると GitHub Actions がストアへデプロイする（`deploy_contents.yml` / `deploy_theme.yml`）。**main へのマージは本番ストアへの反映を意味する**
- `make deploy/*` / `make download/*` / `make import` / `./bin/theme deploy` などストアを直接書き換えるコマンドはローカルから実行しない。認証情報も手元にはない
- Content Kit のデプロイは、ストアにないブログと記事を作成し、既存のものは本文と front matter で指定したフィールドを更新する。削除はしない
