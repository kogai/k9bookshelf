このサイトの開発レポジトリではBazelでビルドを管理しています。

ビルド最適化の一環としてBazelのremote cacheを設定したので、ログを残しておきます。

## キャッシュサーバの選択

[https://docs.bazel.build/versions/master/remote-caching.html](%E3%83%89%E3%82%AD%E3%83%A5%E3%83%A1%E3%83%B3%E3%83%88) によれば、自前でキャッシュサーバを用意する方法と、GCPなどのマネージドサービスを利用する方法があるようです。

BazelはGoogleのプロジェクトだけに、デフォルトでGCPのcredentialsを用いて認証するオプションが付いていますので、これを用いるのが一番手っ取り早そうです。 (月数百円程度のコストは掛かってしまいますが)

自宅サーバが遊んでいればそっちを流用してみても良かったかも。

## GCPの設定

キャッシュ用のストレージへのアクセス権限と鍵発行を担うアカウントを作成して、ストレージへ紐付けます。

### Projectを作成

既存の関連Projectがなければ作成します。

![](https://cdn.shopify.com/s/files/1/0512/0091/7703/files/create-project_480x480.png?v=1609060704)

### Roleを作ってService Accountに紐付け

Remote cacheにはCloud Storageを使います。

その前にCloud StorageにアクセスするためのService Accountと、Accountに紐付けるRoleを用意します。

![](https://cdn.shopify.com/s/files/1/0512/0091/7703/files/create-role_480x480.png?v=1609060715)

Remote cacheのためには、以下の3つの権限が付いていれば良いようです。

- `storage.objects.get`
- `storage.objects.create`
- `storage.objects.update`

次にService Accountを作成して、先程作成したRoleを紐付けます。

![](https://cdn.shopify.com/s/files/1/0512/0091/7703/files/create-service-account_480x480.png?v=1609060723)

### Keyを作成

Service AccountからKeyを発行します。

### Cloud StorageにBucketを作成

Bucketを作成します。objecgtにはBazelからしかアクセスしないので、Access ControlにはUniformedを設定します。

### PERMISSIONSを追加

先程作成したService Accountを追加します。

![](https://cdn.shopify.com/s/files/1/0512/0091/7703/files/add-role_480x480.png?v=1609060696)

### LifyCycleを設定

キャッシュなのでいつまでも残しておく必要はありません。

Actionに `Delete Object` を選択して、conditionsの `Age` に適当な期間を入力します。 (私は30日で設定しています)

## ローカルに設定

### .bazelrcを設定

`.bazelrc` に、設定したBucketをremote cacheとして紐付けます。

```sh
run --remote_cache=https://storage.googleapis.com/your-bucket-name
run --google_default_credentials
build --remote_cache=https://storage.googleapis.com/your-bucket-name
build --google_default_credentials
test --verbose_failures --remote_cache=https://storage.googleapis.com/your-bucket-name
test --google_default_credentials

```

これでremote cacheの設定は完了です。

`google_default_credentials` はGoogle Cloudへアクセスする際の、 [実行環境のデフォルトのService Accountをクレデンシャルとして用いる](https://cloud.google.com/iam/docs/service-accounts#default) オプションです。 先程発行したKeyで`gcloud`コマンドから認証しておくことで、BazelがBucketへアクセス出来るようになります。

ただしこのためだけに`gcloud`コマンドをインストールするのも手間なので、以下のようにKeyを直接クレデンシャルとして設定することも可能です。

```diff
--- .bazelrc
+++ .bazelrc
@@ -1,5 +1,3 @@
+try-import %workspace%/dev.bazelrc

```

```sh
# user.bazelrc
run --nogoogle_default_credentials
build --nogoogle_default_credentials
test --nogoogle_default_credentials

run --google_credentials=secrets.json
build --google_credentials=secrets.json
test --google_credentials=secrets.json

```

### ビルドしてみる

設定が正しいか確認してみましょう。

```sh
$ npx bazelisk build //...

```

設定に問題があれば、404系のエラーがログに表示されるはずです。 ビルドが完了したら、Bucketにobjectが出来ているかも確認できます。

## GitHub Actionsに設定

最後にCI環境でremote cacheを設定します。 このサイトの [開発レポジトリではGitHub Actionsを使って](https://github.com/kogai/k9bookshelf/blob/main/.github/workflows/test_go.yml) いますのでそれに応じた作業ログになります。

### Secretsに設定

[gcloudの設定用actions](https://github.com/google-github-actions/setup-gcloud/tree/master/setup-gcloud) が提供されているので、こちらを用います。

先程発行したKeyを適当なsecretsに設定して、actionからgcloudが認証出来るようにします。

```yml
# GitHub Actionsの設定
  - name: Set up Cloud SDK
    uses: google-github-actions/setup-gcloud@master
    with:
      project_id: "your-project-id"
      service_account_key: ${{ secrets.GCP_BAZEL_CACHE_KEY }}
      export_default_credentials: true

```

## まとめ

以上で設定はすべて完了です。

[月ごとの費用感](https://cloud.google.com/products/calculator/#id=a10aa66c-0e7d-4ee6-bfb4-d3f5a785018e) は数ドル~くらいのようです。

実はこのサイトの開発レポジトリくらいの規模感だと費用対効果はいまいちなのですが、別のレポジトリで各7~8分くらいのaction(がcommit毎に複数)が2~3分くらいにまで短縮されました。
