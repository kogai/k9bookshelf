## 出版社より

Learn proven patterns, techniques, and tricks to take full advantage of the Node.js platform. Master well-known design principles to create applications that are readable, extensible, and that can grow big.

**Key Features**

- Learn how to create solid server-side applications by leveraging the full power of Node.js 14
- Understand how Node.js works and learn how to take full advantage of its core components as well as the solutions offered by its ecosystem
- Avoid common mistakes and use proven patterns to create production grade Node.js applications

**Book Description**

In this book, we will show you how to implement a series of best practices and design patterns to help you create efficient and robust Node.js applications with ease.

We kick off by exploring the basics of Node.js, analyzing its asynchronous event driven architecture and its fundamental design patterns. We then show you how to build asynchronous control flow patterns with callbacks, promises and async/await. Next, we dive into Node.js streams, unveiling their power and showing you how to use them at their full capacity. Following streams is an analysis of different creational, structural, and behavioral design patterns that take full advantage of JavaScript and Node.js. Lastly, the book dives into more advanced concepts such as Universal JavaScript, scalability and messaging patterns to help you build enterprise-grade distributed applications.

Throughout the book, you’ll see Node.js in action with the help of several real-life examples leveraging technologies such as LevelDB, Redis, RabbitMQ, ZeroMQ, and many others. They will be used to demonstrate a pattern or technique, but they will also give you a great introduction to the Node.js ecosystem and its set of solutions.

**What you will learn**

- Become comfortable with writing asynchronous code by leveraging callbacks, promises, and the async/await syntax
- Leverage Node.js streams to create data-driven asynchronous processing pipelines
- Implement well-known software design patterns to create production grade applications
- Share code between Node.js and the browser and take advantage of full-stack JavaScript
- Build and scale microservices and distributed systems powered by Node.js
- Use Node.js in conjunction with other powerful technologies such as Redis, RabbitMQ, ZeroMQ, and LevelDB

**Who this book is for**

This book is for developers and software architects who have some prior basic knowledge of JavaScript and Node.js and now want to get the most out of these technologies in terms of productivity, design quality, and scalability. Software professionals with intermediate experience in Node.js and JavaScript will also find valuable the more advanced patterns and techniques presented in this book.

This book assumes that you have an intermediate understanding of web application development, databases, and software design principles.

## DeepL粗訳

Node.js プラットフォームを最大限に活用するための実証済みのパターン、テクニック、トリックを学びます。よく知られている設計原則をマスターして、読みやすく、拡張性があり、大きく成長できるアプリケーションを作成します。

**主な機能**

- Node.js 14のフルパワーを活用して、堅実なサーバーサイドアプリケーションを作成する方法を学びます。
- Node.js がどのように動作するかを理解し、そのコアコンポーネントとエコシステムによって提供されるソリューションをフルに活用する方法を学びます。
- よくある間違いを避け、実績のあるパターンを使用して実運用レベルのNode.jsアプリケーションを作成します。

**本の説明**

この本では、効率的で堅牢なNode.jsアプリケーションを簡単に作成するための一連のベストプラクティスとデザインパターンを実装する方法を紹介します。

私たちは、Node.js の基本を探り、その非同期イベント駆動アーキテクチャと基本的なデザインパターンを分析することから始めます。次に、コールバック、プロミス、非同期/待ち受けを使って非同期制御フローパターンを構築する方法を紹介します。次に、Node.jsのストリームに飛び込み、そのパワーを明らかにし、それらをフルに活用する方法を示します。ストリームに続くのは、JavaScriptとNode.jsをフルに活用するさまざまな生成、構造、および動作の設計パターンの分析です。最後に、Universal JavaScript、スケーラビリティ、メッセージングパターンなどのより高度な概念に飛び込んで、エンタープライズグレードの分散アプリケーションを構築するのに役立ちます。

この本を通して、LevelDB、Redis、RabbitMQ、ZeroMQ、その他多くの技術を活用した実例を用いて、Node.jsが実際に動作しているのを見ることができます。これらの例はパターンやテクニックを示すために使用されますが、Node.js のエコシステムとそのソリューションのセットを紹介するためにも使用されます。

**あなたが学ぶこと**

- コールバック、プロミス、async/await構文を利用して、非同期コードを書くことに慣れていきます。
- Node.js ストリームを活用して、データ駆動型の非同期処理パイプラインを作成します。
- よく知られたソフトウェアのデザインパターンを実装して、生産グレードのアプリケーションを作成します。
- Node.jsとブラウザ間でコードを共有し、フルスタックJavaScriptを活用する
- Node.jsを利用したマイクロサービスと分散システムの構築とスケーリング
- Redis、RabbitMQ、ZeroMQ、LevelDBのような他の強力なテクノロジーとNode.jsを組み合わせて使用します。

**この本は誰のために書かれたのでしょうか？**

この本は、JavaScript と Node.js の基本的な知識を持っていて、生産性、設計品質、スケーラビリティの面でこれらの技術を最大限に活用したいと考えている開発者やソフトウェアアーキテクトのための本です。Node.js と JavaScript の中級者の経験を持つソフトウェアの専門家は、本書で紹介されているより高度なパターンやテクニックにも価値を見いだすことができるでしょう。

本書は、Webアプリケーション開発、データベース、ソフトウェア設計の原則についての中級的な理解があることを前提としています。

---

Node.jsでコードを書くときのデザインプラクティスやそのプラクティスが重要である理由について解説した書籍です。 JavaScriptのイベントループによる実行モデルを、時に分かりやすい図を挟みながら詳細に解説が載っています。 業務でJavaScriptを書いていて、普通に書けるようにはなってきたけど、もう一歩深入りしたいという方にオススメです。（特にJavaScriptが最初の言語である方）

JavaScriptのイベントループモデルは、「WebブラウザというGUI環境で実行される」「そのため、ユーザーの操作をブロックしてはならない」「だから、イベントループによる実行モデルを採用している」という前提がありますが、多くの入門書ではそのことに紙面を割くことが少ない印象です。 （入門書であるから、読んだ時にそのことに気づかなかっただけかも知れませんが）

この本では、あまり語られることのない、JavaScriptが「どうやって」「なぜそのように」実行されているのかということが解説されています。 （表題のDesign Patternは、「そういう理由であるからこういったデザインを採用する」という意味と解釈しています）

なお、表題のイメージとは異なり、いわゆるデザパタ本的な性質は薄いです。
