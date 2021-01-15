[k9bookshelf](https://k9bookshelf.com)のテーマやコンテンツ、サブツールの管理レポジトリ。

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
