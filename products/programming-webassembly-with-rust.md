2020年末時点でもまだ珍しい、WebAssemblyをテーマに据えた本です。
WebAssemblyの背景・基礎から解説を始めて、Web環境での実行へと進めます。
ステップ・バイ・ステップのチュートリアル的な内容なので、初学者にオススメできます。ただしRustは普通に書ける程度が期待されているのでご注意を（特に難しい表現は出てきませんが）。

個人的には三章 Working with Non-Web Hosts が興味深かったです。
元々Web環境で実行されるものとして出てきたWebAssemblyですが、実はその仕様は実行環境を前提してません（と言うかWasmコードを実行するホスト側の仕様が定められている）。
そのため、Webブラウザ以外の実行環境が[いくつも開発されています](https://github.com/appcypher/awesome-wasm-runtimes)。

本書で採用されているのは、欧州を基盤にするEtherium系企業 [Parity Technologies ](https://www.parity.io/)が開発を主導する [wasmi](https://github.com/paritytech/wasmi) です。
EtheriumはSmart Contractの次期実行環境としてWebAssemblyをベースにしたものが検討されているらしく、その関連で開発されたものなのでしょう。
[Edge computingの実行環境としても検討される](https://www.publickey1.jp/blog/19/fastly_ctowebassemblylucet.html)など、WebAssemblyの広がりはむしろWeb環境でないところに出てくるのかも知れません。

特にWebアプリケーションの文脈では、高負荷な計算をクライアントで行う必然性は、プライバシーへの配慮などのケースを除けばあまり強くないとも言えます。
そういう意味では、Non-Web HostsはWebAssemblyの実行環境として傍流と言うことにはならず、そこに紙幅を割いた本書の内容は一読に値するものと思います。

なお、前述の awsome-wasm-runtimeには[このサイトの運営者作の実行環境もある](https://github.com/appcypher/awesome-wasm-runtimes#wasmvm-top-1)ので、良ければ見てみて下さい。
