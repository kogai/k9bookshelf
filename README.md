[k9bookshelf](https://k9bookshelf.com)のテーマやコンテンツ、サブツールの管理レポジトリ。

## 必要なもの

- Go 1.24
- Node.js 20.9 以降 / Yarn（テーマのビルドと lint）
- Theme Kit は `make bin/theme` で `bin/theme` に取得される（チェックサム検証あり）

## テスト

```sh
$ make test
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
