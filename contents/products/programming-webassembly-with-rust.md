## 出版社より

WebAssembly fulfills the long-awaited promise of web technologies: fast code, type-safe at compile time, execution in the browser, on embedded devices, or anywhere else. Rust delivers the power of C in a language that strictly enforces type safety. Combine both languages and you can write for the web like never before! Learn how to integrate with JavaScript, run code on platforms other than the browser, and take a step into IoT. Discover the easy way to build cross-platform applications without sacrificing power, and change the way you write code for the web.

WebAssembly is more than just a revolutionary new technology. It’s reshaping how we build applications for the web and beyond. Where technologies like ActiveX and Flash have failed, you can now write code in whatever language you prefer and compile to WebAssembly for fast, type-safe code that runs in the browser, on mobile devices, embedded devices, and more. Combining WebAssembly’s portable, high-performance modules with Rust’s safety and power is a perfect development combination.

Learn how WebAssembly’s stack machine architecture works, install low-level wasm tools, and discover the dark art of writing raw wast code. Build on that foundation and learn how to compile WebAssembly modules from Rust by implementing the logic for a checkers game. Create wasm modules in Rust to interoperate with JavaScript in many compelling ways. Apply your new skills to the world of non-web hosts, and create everything from an app running on a Raspberry Pi that controls a lighting system, to a fully-functioning online multiplayer game engine where developers upload their own arena-bound WebAssembly combat modules. Get started with WebAssembly today, and change the way you think about the web.

**What You Need:**

You’ll need a Linux, Mac, or Windows workstation with an Internet connection. You’ll need an up-to-date web browser that supports WebAssembly. To work with the sample code, you can use your favorite text editor or IDE. The book will guide you through installing the Rust and WebAssembly tools needed for each chapter.


---

2020年末時点でもまだ珍しい、WebAssemblyをテーマに据えた本です。 WebAssemblyの背景・基礎から解説を始めて、Web環境での実行へと進めます。 ステップ・バイ・ステップのチュートリアル的な内容なので、初学者にオススメできます。ただしRustは普通に書ける程度が期待されているのでご注意を（特に難しい表現は出てきませんが）。

個人的には三章 Working with Non-Web Hosts が興味深かったです。 元々Web環境で実行されるものとして出てきたWebAssemblyですが、実はその仕様は実行環境を前提してません（と言うかWasmコードを実行するホスト側の仕様が定められている）。 そのため、Webブラウザ以外の実行環境が[いくつも開発されています](https://github.com/appcypher/awesome-wasm-runtimes)。

本書で採用されているのは、欧州を基盤にするEtherium系企業 [Parity Technologies ](https://www.parity.io/)が開発を主導する [wasmi](https://github.com/paritytech/wasmi) です。 EtheriumはSmart Contractの次期実行環境としてWebAssemblyをベースにしたものが検討されているらしく、その関連で開発されたものなのでしょう。[Edge computingの実行環境としても検討される](https://www.publickey1.jp/blog/19/fastly_ctowebassemblylucet.html)など、WebAssemblyの広がりはむしろWeb環境でないところに出てくるのかも知れません。

特にWebアプリケーションの文脈では、高負荷な計算をクライアントで行う必然性は、プライバシーへの配慮などのケースを除けばあまり強くないとも言えます。 そういう意味では、Non-Web HostsはWebAssemblyの実行環境として傍流と言うことにはならず、そこに紙幅を割いた本書の内容は一読に値するものと思います。

なお、前述の awsome-wasm-runtimeには[このサイトの運営者作の実行環境もある](https://github.com/appcypher/awesome-wasm-runtimes#wasmvm-top-1)ので、良ければ見てみて下さい。


