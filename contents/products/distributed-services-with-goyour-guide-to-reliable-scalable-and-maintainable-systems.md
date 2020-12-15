分散システムにおける一つのServiceをステップ・バイ・ステップで構築していく本です。基本的にはアプリケーションレイヤーの話がメインですが、最後に少しだけk8sの話が出てきます。

2020年現在では単一のデプロイ単位からなるアプリケーションで成立するサービスってほぼないと思うのですが、そういう意味で新しい世代のWeb開発入門本として読めるかも知れません。
(とは言えいわゆるモノリシックなWebアプリケーションは普通に作れる知識レベルを前提とした本だと思います)

この手の書籍は大抵そうですが、コード例にはテストも付いているので読み進めやすいです。

ところでボブ、アリス、マロリーって公開鍵暗号の説明をする時の定形キャラクターだったんですね。知らなかった。

## 出版社より

You know the basics of Go and are eager to put your knowledge to work. This book is just what you need to apply Go to real-world situations. You’ll build a distributed service that’s highly available, resilient, and scalable. Along the way you’ll master the techniques, tools, and tricks that skilled Go programmers use every day to build quality applications. Level up your Go skills today.

Take your Go skills to the next level by learning how to design, develop, and deploy a distributed service. Start from the bare essentials of storage handling, then work your way through networking a client and server, and finally to distributing server instances, deployment, and testing. All this will make coding in your day job or side projects easier, faster, and more fun.

Lay out your applications and libraries to be modular and easy to maintain. Build networked, secure clients and servers with gRPC. Monitor your applications with metrics, logs, and traces to make them debuggable and reliable. Test and benchmark your applications to ensure they’re correct and fast. Build your own distributed services with service discovery and consensus. Write CLIs to configure your applications. Deploy applications to the cloud with Kubernetes and manage them with your own Kubernetes Operator.

Dive into writing Go and join the hundreds of thousands who are using it to build software for the real world.

**What You Need:**

Go 1.11 and Kubernetes 1.12.

---

## DeepL粗訳

あなたは囲碁の基本を知っていて、その知識を実践したいと思っています。この本は、Goを実際の状況に適用するために必要なものです。高度に利用可能で、回復力があり、スケーラブルな分散サービスを構築することができます。その過程で、熟練したGoプログラマーが日々使用している、高品質のアプリケーションを構築するためのテクニック、ツール、およびトリックを習得します。今すぐ囲碁スキルをレベルアップさせましょう。

分散サービスの設計、開発、配備の方法を学ぶことで、Goのスキルを次のレベルに引き上げます。ストレージ処理の基本的な要素から始め、クライアントとサーバのネットワーク構築、サーバインスタンスの分散、デプロイメント、テストまでを学習します。これらすべてが、本業やサイドプロジェクトでのコーディングをより簡単に、より速く、より楽しくしてくれます。

アプリケーションとライブラリをモジュール化してメンテナンスが簡単になるようにレイアウトします。gRPC を使用して、ネットワーク化された安全なクライアントとサーバを構築します。メトリクス、ログ、トレースでアプリケーションを監視し、デバッガブルで信頼性の高いものにします。アプリケーションのテストとベンチマークを行い、アプリケーションが正しく高速であることを確認します。サービスディスカバリとコンセンサスを使って独自の分散サービスを構築します。CLI を書いてアプリケーションを設定します。Kubernetesを使ってアプリケーションをクラウドにデプロイし、独自のKubernetes Operatorで管理します。

Goの記述に飛び込んで、実際の世界でGoを使用してソフトウェアを構築している何十万人もの人たちの仲間入りをしましょう。

**必要なもの:** 必要なもの

Go 1.11とKubernetes 1.12。
