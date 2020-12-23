## 出版社より

Learn to build configuration file readers, data readers, model-driven code generators, source-to-source translators, source analyzers, and interpreters. You don’t need a background in computer science–ANTLR creator Terence Parr demystifies language implementation by breaking it down into the most common design patterns. Pattern by pattern, you’ll learn the key skills you need to implement your own computer languages.

Knowing how to create domain-specific languages (DSLs) can give you a huge productivity boost. Instead of writing code in a general-purpose programming language, you can first build a custom language tailored to make you efficient in a particular domain.

The key is understanding the common patterns found across language implementations. Language Design Patterns identifies and condenses the most common design patterns, providing sample implementations of each.

The pattern implementations use Java, but the patterns themselves are completely general. Some of the implementations use the well-known ANTLR parser generator, so readers will find this book an excellent source of ANTLR examples as well. But this book will benefit anyone interested in implementing languages, regardless of their tool of choice. Other language implementation books focus on compilers, which you rarely need in your daily life. Instead, Language Design Patterns shows you patterns you can use for all kinds of language applications.

You’ll learn to create configuration file readers, data readers, model-driven code generators, source-to-source translators, source analyzers, and interpreters. Each chapter groups related design patterns and, in each pattern, you’ll get hands-on experience by building a complete sample implementation. By the time you finish the book, you’ll know how to solve most common language implementation problems.

## DeepL粗訳

設定ファイルリーダ、データリーダ、モデル駆動型コードジェネレータ、ソースからソースへのトランスレータ、ソースアナライザ、インタプリタの構築を学ぶことができます。コンピュータサイエンスのバックグラウンドは必要ありません。ANTLRの作成者であるTerence Parr氏は、言語の実装を最も一般的な設計パターンに分解することで、言語の実装を理解します。パターンごとに、独自のコンピュータ言語を実装するために必要な主要なスキルを学ぶことができます。

ドメイン固有言語(DSL)の作成方法を知ることで、生産性を大幅に向上させることができる。汎用プログラミング言語でコードを書くのではなく、特定のドメインで効率的に動作するようにカスタマイズされたカスタム言語を最初に構築することができます。

重要なのは、言語の実装で見られる共通のパターンを理解することです。Language Design Patterns』では、最も一般的なデザイン・パターンを特定して凝縮し、それぞれの実装サンプルを提供しています。

パターンの実装はJavaを使用していますが、パターン自体は完全に一般的なものです。いくつかの実装では、よく知られているANTLRパーサジェネレータを使用しているので、読者は本書がANTLRのサンプルの優れた情報源であることに気づくでしょう。しかし、この本は、言語実装に興味がある人であれば、どんなツールを使っているかに関わらず、誰にでも役立つでしょう。他の言語実装の本では、日常生活ではほとんど必要のないコンパイラに焦点を当てています。その代わりに、この本では、あらゆる種類の言語アプリケーションに使えるパターンを紹介しています。

設定ファイル・リーダ、データ・リーダ、モデル駆動型コード・ジェネレータ、ソース・ツー・ソース・トランスレータ、ソース・アナライザ、インタプリタの作成方法を学ぶことができます。各章では、関連する設計パターンをグループ化し、各パターンでは、完全なサンプル実装を構築することで実践的な経験を得ることができます。この本を読み終える頃には、ほとんどの一般的な言語実装の問題を解決する方法がわかるようになっています。

---

表題の通り、言語実装のデザインパターンを色々紹介していく本です。 いくつかのパートに分けて言語解析器、インタープリタ、コンパイラ(Translator&Generatorとして紹介)の内部に現れるパターンを紹介していきます。

言語実装と言うとドラゴンブックを始めとする分厚い書籍を紐解いて、ウィザード的なプログラマーが取り組むものというイメージがあるかも知れません。 確かに、いわゆる汎用プログラミング言語の実装に直接関わるプログラマーは少ないですが、そうでなくとも言語実装のパターンを掴んでおくことは十分に意味のあることだとこの本は主張します。

LinterやFormatter、設定ファイルの解析器、あるいはマークダウンのような軽量マークアップ言語からHTMLを生成するなど、言語実装におけるパターンの適用範囲は実は一般的なプログラマーにとっても身近なものであるというのです。 (あるい単に教養として、ということでも良いかも知れません)

私はこの本にあたった後に、i18nの辞書となるjsonファイルからTypeScriptなどの型定義ファイルを生成するという [ツール](https://github.com/kogai/typed_i18n) を書いたことがあるのですが、 この本に出てくるコンセプトが溶け込んでいたと言って良いと思います。 例えば入力であるjsonを何らかの中間表現に落とし込み、IRをGeneratorが型定義として書き出すといったような、この本で解説されているパターンです。

LinterやFormatterを常用してレポジトリを管理することは一般的になっていますし、チーム毎に適用したいルールというものも出てくるでしょう。 そういった時の道案内にも良いと思います。
