---
title: "RxJSのカスタムオペレータを作った"
published_at: "2016-09-10T17:09:40+09:00"
tags:
  - "javascript"
  - "typescript"
  - "作ったもの"
  - "npm"
---

最近RxJSを使って状態の遷移を表現するような設計でアプリケーションを作っていた。

ReduxをはじめとするFlux由来のライブラリをいくつか試してみたが、どうもreducerの役割をうまく扱いきれなかったり(凝集を疎にしてしまっているように思えてならなかったり...理解が浅いだけかも知れませんが)、非同期処理間に複数の依存関係が生じるような要件が提示されていたので(例えばある要素をクリックするとAPIリクエストが走り、そのレスポンスから別の要素のクリック・APIリクエストが可能になる、といった具合)、ReactiveProgramming的な考え方で作るのが最もフィットすると考えたからだ。

最初は入力の起点毎にSubjectを作成して、各DOM要素にsubject.nextを紐付けるようなやり方をしていたのだが、ちょっと冗長というか無駄だと思う。
Reduxっぽく、単一のdispatcher関数に、ユーザの関心事の種類を表す情報(ActionType)を載せて分岐するようにした方が、スマートだろう。
[redux-observable](https://github.com/redux-observable/redux-observable)なんかのソースコードを読んでいると[createEpicMiddleware](https://github.com/redux-observable/redux-observable/blob/master/src/createEpicMiddleware.js#L17)あたりでそういうアプローチを取っているように読み取った。

※Twitterで教えて頂いたが、[boajs](https://github.com/bouzuya/b-o-a)も同様のアプローチを取っているそうで、[このあたり](https://github.com/bouzuya/boa-core/blob/master/src/run.ts#L7)のことかなと読み取った(間違っていたらすいません)

で、このアプローチを取ると、共有しているsubject由来のストリームにフィルタをかけるような処理を頻繁に書くようになる。

```javascript
const subject = new Subject<any>();
const action$ = subject.asObservable();
  
action$ // subject由来のストリーム
  .filter(action => action.type === "SOME_ACTION_TYPE");
```

できればいわゆるアクションにはこういう型を期待して、

```javascript
interface Action<T> {
  type: string;
  payload?: T;
}
```

こんな風にアクションのペイロードを引き出せるとうれしい

```javascript
action$ // subject由来のストリーム
  .ofType("SOME_ACTION_TYPE");
// => Observable<Payload>
```

RxJS(というかたぶんRx系ライブラリはみんなそうだろうが)には、ライブラリのユーザがオペレータを拡張できるようになっていて、[ガイドライン](https://github.com/ReactiveX/rxjs/blob/master/doc/operator-creation.md)も用意されている。

いたずらにオペレータを増やすのはきな臭い予感しかしないが、ここぞという場面で使えばコードのリーダビリティと凝集度に大いに貢献してくれそうだ。
[Functional Reactive Programming](https://www.manning.com/books/functional-reactive-programming)には、そもそもRP系ライブラリで提供されるオペレータは、10個のプリミティブなオペレータの組み合わせによる拡張に過ぎない、みたいな話もあった。

さて、実はこのofTypeオペレータは[redux-observableの内部](https://github.com/redux-observable/redux-observable/blob/master/src/ActionsObservable.js)で実装されているものなのだった。
ただ、上述のようにことさらredux-observableを用いていないアプリケーションであっても、RxJSにFlux的なエッセンスを加えて構成する場合にとても有用なオペレータだ。
そこで、このオペレータに着想を得て(少しインターフェースを変えている)似たようなオペレータを作成・npmに公開した。

[of-type-operator](https://github.com/kogai/of-type-operator)

[型定義](https://github.com/kogai/of-type-operator/blob/master/dist/index.d.ts)にあるように、

```javascript
<T>(actionType: string, pickBy?: PickByFunction): Observable<T>
```

という感じだ。
redux-observableとは、複数のアクションタイプを受け取れない代わりに、project関数を受け取って、返す`Observable<Result>`を操作できるようにしている。
例えばこういう感じのアクションを定義している時もあるだろうという想定だ

```javascript
interface Action {
  type: string;
  userName: string;
  userMail: string;
}
```

作ってみてObservableへの理解は多少深まったような気がする。
特にオペレータの、Subscriberの子クラスをsubscribeのチェーンにつなげてliftする、みたいな基本的な構造を読むことになったのが良かった。

ところで、RxJSのコードを眺めた後に『すごいHaskellたのしく学ぼう！』を再読してて、liftについて触れている箇所があることに気付いた時は興奮した。
関数の持ち上げ以外にも、Haskellに入門した時に触れた概念でRxへの理解がだいぶ進んだところがあるように思う。IOアクションとか。

というわけで、よかったら使ってみてください。
