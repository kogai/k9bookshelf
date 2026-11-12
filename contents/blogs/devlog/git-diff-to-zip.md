---
title: "git-diffで差分ファイルのみをarchive.zipにする"
published_at: "2015-03-17T00:00:00+09:00"
tags:
  - "git"
---

## やりたいこと

- クライアント案件で、本番サーバーへのアクセス権がない。
- 更新作業の際には差分ファイルのみを納品するように要求されている。

## git archiveで差分抽出

`git archive`コマンドで、前回納品時との`diff`を抽出して`zip`化する

```shell
git archive --format=zip HEAD `git diff master --name-only` -o archive.zip

```
`git diff`で`master`ブランチ(前回納品時の状態を想定)との差分ファイルリストを取得して、`git archive`に渡してあげる。
これでアーカイブ化できました。

## --diff-filterを使う

上記の手法だと差分の中に削除したファイルも含まれてしまい、以下のエラーメッセージが出ます。

```shell
fatal: pathspec 'path/to/not/exist/file' did not match any files

```

削除したファイルは`git archive`対象のリストに渡らないようにすれば良さそうです。
[公式ドキュメントのgit-diffのページ](http://git-scm.com/docs/git-diff)を見ると、`--diff-filter`というオプションがありました。

> Select only files that are Added (A), Copied (C), Deleted (D), Modified (M), Renamed (R)

このオプションを使って追加したファイル(A)と編集したファイル(M)だけにフィルタリングされたリストが作れそうです。

```shell
git archive --format=zip HEAD `git diff master --name-only --diff-filter=AM` -o archive.zip
```

これでOK。
困ったときは公式ドキュメントに限りますね。
