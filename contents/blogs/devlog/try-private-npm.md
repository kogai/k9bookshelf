---
title: "private npmを試してみた"
published_at: "2015-04-15T17:03:25+09:00"
tags:
  - "npm"
  - "node.js"
---

以前より話が出ていた、npmのプライベートモジュールが公開されていました。
[private-modules](https://www.npmjs.com/private-modules)
![](https://www.evernote.com/shard/s168/sh/241a329b-f2a3-4a33-ad08-4b7fc59d573c/6ec829c5154bd20530db0c9b20b7fffc/deep/0/npm.png)

## 機能

>- create and share unlimited private modules for $7/month
>- Host as many private packages as you want
>- Give read access or read-write access for those packages to any other paid user
>- Install and use any packages that other paid users have given you read access to
>- Collaborate on any packages that other paid users have given you write access to

要旨としては

- 月額7ドルでモジュールをプライベート公開に出来る
- 任意の課金済みユーザーに読み書き権限の発行ができる

ということのようです。

>Currently, private packages are only available for individual users, but support for organization accounts is coming soon. Feel free to create a user for your organization in the meantime, and we can upgrade it to an organization when support is here.

企業(組織)アカウント向けの対応はこれからのようです。


## 導入の仕方

1. [npmjs](https://www.npmjs.com)にサインアップして課金しておきます。
※クレジットカード払いのみのようです。

2. 任意のnpmの`package.json`の`name`プロパティに自分のユーザーアカウント名を入れます
```json
  public
  "name": "modulename"

  private
  "name": "@username/modulename"
```

3. `npm publish`で公開します。


説明が必要なほどの手順はないですね。。。
既にnpmを公開したことがある人なら、ほぼやることはありません。
課金して、モジュールの名前にユーザー名を入れるだけ。

ロードマップにあるという、企業アカウント向けの対応が始まったらnpmの使い方が変わってくるかも知れません。
