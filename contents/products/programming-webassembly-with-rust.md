## 出版社より

WebAssembly fulfills the long-awaited promise of web technologies: fast code, type-safe at compile time, execution in the browser, on embedded devices, or anywhere else. Rust delivers the power of C in a language that strictly enforces type safety. Combine both languages and you can write for the web like never before! Learn how to integrate with JavaScript, run code on platforms other than the browser, and take a step into IoT. Discover the easy way to build cross-platform applications without sacrificing power, and change the way you write code for the web.

WebAssembly is more than just a revolutionary new technology. It’s reshaping how we build applications for the web and beyond. Where technologies like ActiveX and Flash have failed, you can now write code in whatever language you prefer and compile to WebAssembly for fast, type-safe code that runs in the browser, on mobile devices, embedded devices, and more. Combining WebAssembly’s portable, high-performance modules with Rust’s safety and power is a perfect development combination.

Learn how WebAssembly’s stack machine architecture works, install low-level wasm tools, and discover the dark art of writing raw wast code. Build on that foundation and learn how to compile WebAssembly modules from Rust by implementing the logic for a checkers game. Create wasm modules in Rust to interoperate with JavaScript in many compelling ways. Apply your new skills to the world of non-web hosts, and create everything from an app running on a Raspberry Pi that controls a lighting system, to a fully-functioning online multiplayer game engine where developers upload their own arena-bound WebAssembly combat modules. Get started with WebAssembly today, and change the way you think about the web.

**What You Need:**

You’ll need a Linux, Mac, or Windows workstation with an Internet connection. You’ll need an up-to-date web browser that supports WebAssembly. To work with the sample code, you can use your favorite text editor or IDE. The book will guide you through installing the Rust and WebAssembly tools needed for each chapter.

## DeepL粗訳

WebAssemblyは、高速なコード、コンパイル時の型安全性、ブラウザでの実行、組み込みデバイスでの実行、その他の場所での実行など、Web技術の待望の約束を実現します。Rust は、型の安全性を厳密に守る言語で C のパワーを提供します。両方の言語を組み合わせることで、これまでにないほどのWebのための記述が可能になります。JavaScript との統合方法を学び、ブラウザ以外のプラットフォームでコードを実行し、IoT への一歩を踏み出しましょう。パワーを犠牲にすることなくクロスプラットフォームのアプリケーションを簡単に構築する方法を発見し、Web用のコードの書き方を変えましょう。

WebAssembly は単なる革命的な新技術ではありません。WebAssembly は、私たちが Web 用のアプリケーションを構築する方法を再形成しています。ActiveXやFlashのような技術が失敗してきたところに、好きな言語でコードを書き、WebAssemblyにコンパイルすることで、ブラウザやモバイルデバイス、組み込みデバイスなどで実行される高速でタイプセーフなコードを作成することができるようになりました。WebAssembly のポータブルで高性能なモジュールと Rust の安全性とパワーを組み合わせることは、完璧な開発の組み合わせです。

WebAssembly のスタックマシンアーキテクチャがどのように機能するかを学び、低レベルの wasm ツールをインストールし、生の無駄なコードを書くというダークアートを発見してください。その基礎を築き、チェッカーゲーム用のロジックを実装することで、Rust から WebAssembly モジュールをコンパイルする方法を学びます。Rustでasmモジュールを作成して、多くの説得力のある方法でJavaScriptと相互運用することができます。新しいスキルを非ウェブホストの世界に適用し、照明システムを制御するRaspberry Pi上で動作するアプリから、完全に機能するオンラインマルチプレイヤーゲームエンジンまで、あらゆるものを作成します。WebAssemblyを今すぐ使い始めて、Webについての考え方を変えてみませんか？

**必要なもの:**

インターネットに接続されたLinux、Mac、またはWindowsワークステーションが必要です。WebAssemblyをサポートする最新のWebブラウザが必要です。サンプルコードを使って作業するには、お好きなテキストエディタやIDEをお使いください。この本では、各章に必要なRustとWebAssemblyツールのインストール方法を説明しています。

---

2020年末時点でもまだ珍しい、WebAssemblyをテーマに据えた本です。 WebAssemblyの背景・基礎から解説を始めて、Web環境での実行へと進めます。 ステップ・バイ・ステップのチュートリアル的な内容なので、初学者にオススメできます。ただしRustは普通に書ける程度が期待されているのでご注意を（特に難しい表現は出てきませんが）。

個人的には三章 Working with Non-Web Hosts が興味深かったです。 元々Web環境で実行されるものとして出てきたWebAssemblyですが、実はその仕様は実行環境を前提してません（と言うかWasmコードを実行するホスト側の仕様が定められている）。 そのため、Webブラウザ以外の実行環境が [いくつも開発されています](https://github.com/appcypher/awesome-wasm-runtimes)。

本書で採用されているのは、欧州を基盤にするEtherium系企業 [Parity Technologies](https://www.parity.io/) が開発を主導する [wasmi](https://github.com/paritytech/wasmi) です。 EtheriumはSmart Contractの次期実行環境としてWebAssemblyをベースにしたものが検討されているらしく、その関連で開発されたものなのでしょう。 [Edge computingの実行環境としても検討される](https://www.publickey1.jp/blog/19/fastly_ctowebassemblylucet.html) など、WebAssemblyの広がりはむしろWeb環境でないところに出てくるのかも知れません。

特にWebアプリケーションの文脈では、高負荷な計算をクライアントで行う必然性は、プライバシーへの配慮などのケースを除けばあまり強くないとも言えます。 そういう意味では、Non-Web HostsはWebAssemblyの実行環境としては決して傍流とは言えず、そこに紙幅を割いた本書の内容は一読に値するものと思います。
