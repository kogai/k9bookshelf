- k9bookshelfの開発レポジトリはbazelでビルドを管理している
- 仕事でremote-cacheを使ったのでこちらにも展開した
- Google Cloud Storageがファーストクラスでサポートされている
- Cloud Storageの`storage.objects`の get, create, updateを持つroleを作成
- Service Accountを作成して紐付ける
- Keyを発行し、githubのsecretesに設定する
- `.bazelrc`にremote cacheの設定を追加する
- user.babelrc でcredentialファイルを指定も出来るようにする

```yml
# GitHub Actionsの設定
  - name: Set up Cloud SDK
    uses: google-github-actions/setup-gcloud@master
    with:
      project_id: "sandbox-bazel-cache"
      service_account_key: ${{ secrets.GCP_BAZEL_CACHE_KEY }}
      export_default_credentials: true

```

```bazel
# .bazelrcに設定して、すべてのbazel build実行時にremote cacheを使う(test, runも同様)
build --remote_cache=https://storage.googleapis.com/your-bucket-name --google_default_credentials

```

[月ごとの費用感](https://cloud.google.com/products/calculator/#id=a10aa66c-0e7d-4ee6-bfb4-d3f5a785018e) は数ドルくらい
