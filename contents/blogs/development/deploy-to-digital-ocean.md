---
title: "DigitalOceanでNode.jsのサーバーを立てる初期設定メモ"
published_at: "2014-12-23T11:58:02+09:00"
tags:
  - "server"
  - "digitalocean"
  - "git"
  - "node.js"
---

## ターミナルからログイン
```shell
ssh root@IP
```

----------

## Ubuntu版

```shell
sudo apt-get update
sudo apt-get install nodejs
sudo apt-get install npm
sudo apt-get install git
npm install -g bower
```

----------

## CentOS版

```shell
yum update
curl -sL https://rpm.nodesource.com/setup | bash -
yum install -y nodejs
yum install git
npm install -g bower
```

### gitデプロイつなぎ込み
[解説記事](https://www.digitalocean.com/community/tutorials/how-to-set-up-automatic-deployment-with-git-with-a-vps)


1. `/var/repo/site.git/` を作成して `git init —bare`
2. `/var/repo/site.git/hooks/` へ移動
3. `cat > post-receive` で`post-receive`ファイルを作成
4. `vim`モードにするなどして以下記入`:x`で保存
```bash
#!/bin/sh
git --work-tree=/var/www/ディレクトリ名 --git-dir=/var/repo/site.git checkout -f
```

5. `chmod +x post-receive`

6. ローカルレポジトリにgitパスをつなぎ込み
```shell
git remote add product ssh://root@IP/var/repo/site.git
```

7. `SSHキー`の登録（新規アカウント作成時のみ）

8. プッシュしてデプロイ

### サーバー永続化
`npm`の`pm2`をインストール
`pm2 start app.js` でサーバー永続化
