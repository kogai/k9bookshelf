[k9bookshelf](https://k9bookshelf.com)のテーマやコンテンツ、サブツールの管理レポジトリ。

## 必要なもの

- Go 1.24
- Node.js / Yarn（テーマのビルドと lint）
  - Node.js 17 以降では `NODE_OPTIONS=--openssl-legacy-provider` が必要（webpack 5.18 が md4 を使うため）
- Theme Kit は `make bin/theme` で `bin/theme` に取得される（チェックサム検証あり）

## テスト

```sh
$ make test
```

`make test` は `contents/` の各ファイルを検証し、`contents/` から生成した HTML を `content/testdata/rendered/` の golden ファイルと比較します。`contents/` や変換処理を変えたら golden ファイルを更新して、差分を確認してからコミットしてください。

```sh
$ make golden
```

テーマのビルドと lint だけを行う（デプロイしない）場合:

```sh
$ make check/theme
```

## テーマの更新

```sh
$ make deploy/theme
```

## 商品本文、ページ、ブログの更新

```sh
$ make deploy/contents
```

## 商品情報のインポート

```sh
$ make ONIX_FILE=20201208.onix import

# dry-run writes snapshot json to onix/snapshots
$ make ONIX_FILE=20201208.onix import/dry
```
