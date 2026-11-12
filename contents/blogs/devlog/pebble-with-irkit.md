---
title: "PebbleからIRKitを操作するアプリを作った"
published_at: "2014-08-28T10:57:43+09:00"
tags:
  - "pebble"
  - "pebble.js"
  - "node.js"
  - "javascript"
---

スマートウォッチのPebbleからWi-FiスマートリモコンのIRKitに家中のリモコン家電をON/OFFにする命令を送れるようにしました。

### 手順

1. 「 IRKitを使ってJavaScriptから家電を操作。外出先でも！」を参考に、clienttokenを取得してdeviceidとclientkey を取得する
2. Pebbleアプリの「Simply.js」 をインストール
3. 実行するjavascriptファイルを書く

[gist url](https://gist.github.com/kogai/4ab760b8f7fd2da66423)
※ globalDefineとremoteCommandの中の数値はユーザー毎に適宜設定する

赤外線送信を一度に行うと、干渉するのか配列内のいずれかの命令しか実行されないっぽかったので、setTimeoutでディレイをかけています。
手首を一振りしてテレビやエアコンが一斉にOFFになっていく光景はちょっとSFチックで素敵です。

### 動作サンプル

[instagram url](http://instagram.com/p/rhJ2B8hBy8/)
