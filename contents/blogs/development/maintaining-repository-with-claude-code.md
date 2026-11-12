---
title: "【草稿】昔の技術ブログの移行とレポジトリの整備を Claude Code に任せてみた"
published_at: "2026-11-13T09:00:00+09:00"
tags:
  - "Claude Code"
  - "Shopify"
  - "GitHub Actions"
---

<!-- TODO: 草稿です。公開前にタイトル（【草稿】を外す）、published_at、本文を見直してください。main にマージするとすぐに公開されます -->

こんにちは。[趣味でマーチャントをやっている](https://k9bookshelf.com/blogs/development/how-and-why-running-bookstore) ものです。

このサイトのテーマやコンテンツは [GitHub のレポジトリ](https://github.com/kogai/k9bookshelf) で管理していて、main ブランチに push すると GitHub Actions がストアにデプロイする仕組みになっています（[以前の記事](https://k9bookshelf.com/blogs/development/to-edit-contents-in-markdown) で紹介した Content Kit を使っています）。

ただ、ここ数年は放置気味で、気づけば CI がすべて動かなくなっていました。ちょうど昔の個人ブログ（[kogai/blog](https://github.com/kogai/blog)）の記事をこちらに移したいと思っていたので、そのついでにレポジトリの整備まで [Claude Code](https://claude.com/claude-code) に任せてみることにしました。

<!-- TODO: 任せてみようと思った動機や、普段の使い方などを追記 -->

## やったこと

最終的にマージされた PR は次のとおりです。一日でこれだけ進みました。

- 昔のブログ記事26本の移行（[#371](https://github.com/kogai/k9bookshelf/pull/371)、[#373](https://github.com/kogai/k9bookshelf/pull/373)、[#377](https://github.com/kogai/k9bookshelf/pull/377)）
- Bazel を廃止して素の Go に移行し、CI を復旧（[#372](https://github.com/kogai/k9bookshelf/pull/372)）
- 依存ライブラリの更新（[#374](https://github.com/kogai/k9bookshelf/pull/374)）
- テストのハーネスと Claude Code 用の設定（[#375](https://github.com/kogai/k9bookshelf/pull/375)、[#376](https://github.com/kogai/k9bookshelf/pull/376)）
- 公開レポジトリ向けの secret 混入チェック（[#379](https://github.com/kogai/k9bookshelf/pull/379)）
- レポジトリの運用ルールの明文化（[#380](https://github.com/kogai/k9bookshelf/pull/380)、[#382](https://github.com/kogai/k9bookshelf/pull/382)）
- テーマのダウンロードの見直し（[#381](https://github.com/kogai/k9bookshelf/pull/381)）

以下、順に振り返ります。

## 昔のブログ記事の移行

最初の依頼は「kogai/blog のコンテンツを k9bookshelf のブログとして移したい」という一言だけでした。

Claude Code はまず両方のレポジトリを読み、Content Kit の `deploy` が **既存の記事の本文を更新するだけ** で、ストアにないブログや記事があるとエラーになることに気づきました。そこで記事を移すだけでなく、Markdown の先頭に YAML の front matter（`title` / `published_at` / `tags`）を書けるようにして、ストアにない記事は作成するように Content Kit を拡張してくれました。

```md
---
title: "Article title"
published_at: "2016-12-22T19:31:47+09:00"
tags:
  - TypeScript
---

本文...
```

旧ブログの記事は日付の書式がまちまち（タイムゾーンなし、日付だけ、`+0900` など）でしたが、変換スクリプトで RFC3339 に揃え、記事内の旧ブログへのリンクも新しい URL に書き換えてくれました。

ちなみに、最初は `devlog` という新しいブログに入れていたのですが、途中で「全部 development に入れて」とお願いして移し直しています。

### 本番で一度失敗した

マージ後の初回デプロイは `commentable: must be in a valid state` というエラーで失敗しました。使っていた go-shopify v2 の `Blog` 構造体に `omitempty` が付いておらず、ブログの作成時に空の `commentable` まで送っていたのが原因です。

失敗したのはブログ作成の時点で、記事は1件も作られていませんでした。Claude Code はログから原因を特定し、title と handle だけを送る修正をすぐ別の PR にしてくれました。

<!-- TODO: このあたりの感想 -->

## 動かなくなっていた CI の復旧

CI は `actions/cache@v2` や `google-github-actions/setup-gcloud@master` など、提供が終わったアクションを使っていたため、起動直後に失敗していました。

`setup-gcloud` を使っていた理由は、Bazel のリモートキャッシュ（[以前の記事](https://k9bookshelf.com/blogs/development/use-remote-cache-of-bazel) で書いたもの）の認証でした。Claude Code から「Bazel ごとやめて素の `go test` / `go run` にした方がシンプルになる」と提案されたので、思い切って切り替えました。Theme Kit は Bazel の `http_file` で取得していましたが、同じチェックサムで検証しながら `make bin/theme` で取得する形になっています。

<!-- TODO: Bazel をやめたことへの感想（以前の記事との対比など） -->

## 依存ライブラリの更新

数年分の更新が溜まっていたので、Go・Node.js・GitHub Actions をまとめて更新しました。go-shopify は v2 のままだと依存が解決できなくなっていたため、v4 に上げています。

この時、v4 ではリクエストパスに `admin/api/<version>` が自動で付くため、自前で組み立てていたパスが二重になるところでした。これは後述のフェイクの Shopify サーバーを使ったテストで見つかっています。

## ハーネスを作る

「移行が終わったらハーネスを整備した方がいいかもね」と言ったところ、テストと Claude Code 用の設定の両方を作ってくれました。

### テスト

- **golden ファイル**: `contents/` 以下の全 Markdown から生成した HTML をレポジトリに保存し、変換結果を比較します。Markdown ライブラリを更新した時に、ストアに反映される HTML の差分が PR 上で見えるようになりました
- **コンテンツの検証**: 置き場所、拡張子、handle、front matter の書式をチェックします
- **フェイクの Shopify Admin API**: `httptest` でフェイクのサーバーを立て、deploy と download が送るリクエストを検証します

download のテストを `go test -race` で動かしたところ、goroutine が共有の `err` に書き込んでいてエラーが握りつぶされるバグが見つかりました。昔の自分のコードです。

### Claude Code 用の設定

- `CLAUDE.md`: レポジトリの構成、確認用のコマンド、「main へのマージは本番への反映」といったルール
- `.claude/settings.json`: テスト系のコマンドは許可し、デプロイ系のコマンドは拒否
- SessionStart フック: クラウドのセッション開始時に Go と Node.js の依存をインストール
- PreToolUse フック: `git commit` の前に secret らしい文字列が含まれていないかチェック

## 運用ルールも育てていく

作業の途中で、気になったことをその都度伝えて `CLAUDE.md` に追記してもらいました。

- 「公開レポジトリだからやばいことコミットしないようにして」→ secret 混入チェックと、公開レポジトリとしてのルール
- 「PR はテーマごとに作ろう」→ 1つの PR には1つの目的、という分け方のルール
- 「Copilot のレビューを入るようにしたから、サブエージェントのレビューと Copilot のレビューのどちらにも応答してね」→ レビューの進め方のルール

レビューのルール自体もサブエージェントにレビューさせたところ、「『CI が通ったらマージする』だと、自分でマージしてよいと読める」という指摘があり、「マージはユーザーの指示があった時だけ」と明記されました。レビューの仕組みが自分のルールを直していくのは、なかなか面白い光景でした。

<!-- TODO: サブエージェントと Copilot のレビューを両方入れてみた感想 -->

## ハマったところ

一日の中で、思わぬところで何度かつまずきました。

- **API のレート制限**: PR を開くたびにテーマをダウンロードするワークフローが、本番デプロイと Shopify Admin API のレート制限（2 calls per second）を取り合って、無関係な PR のチェックを落としていました。週次の定期実行と手動実行だけにしています
- **Rulesets と force-push**: テーマ更新の PR を作る create-pull-request はブランチを force-push で更新しますが、レポジトリの Rulesets で force-push を禁止していたため拒否されていました。該当ブランチを対象から外して解決しました
- **独立にマージした PR の噛み合わせ**: 依存更新とハーネスの PR を別々にマージしたら、golden ファイルが新しい Markdown ライブラリの出力と合わずに main のテストが落ちました。ただ、差分は事前に手で確認していた4記事だけで、ハーネスがちゃんと仕事をしていることの確認にもなりました

## まとめ

<!-- TODO: まとめを書く。例: 任せられたこと / 人間が判断したこと（マージの可否、ブログの置き場所、Rulesets の変更、devlog ブログの削除など）/ 今後やりたいこと -->

この記事の草稿も Claude Code に書いてもらい、人間が手を入れて公開しています。
