---
title: "Reactとredux-Observableでショッピングカートアプリを作ってみる話"
published_at: "2016-11-30T23:16:37+09:00"
tags:
  - "Rx"
---

これは[React Advent Calendar 2016](http://qiita.com/advent-calendar/2016/react)の第 1 日目の記事です。

夏頃から`RxJS`で状態を管理して`React`で描画するようなアプリケーションを作っていました。
`RxJS`に限らず、(F)RP 系ライブラリの導入に際して、導入のメリットが分かりづらいという話があります。
難しい(そうに見える)新しいことを覚えてまで導入するほどのメリットがあるようには思えない、みたいな意見です。
特にライブラリのレポジトリに上がっている Example アプリのような(例えば TODO アプリとか)ものだと、単に複雑性が増しているように見えるということもあります。

これについて、[Functional Reactive Programming](https://www.manning.com/books/functional-reactive-programming)という本に、(F)RP は小さく単純なアプリケーションでサンプルを作ってもメリットが見えてこない、機能同士に複雑な依存関係があるような複雑なアプリケーションになって、はじめてその真価がわかってくる、というようなことが書いてありました。
(この本ではガソリン給油機アプリを想定して、ハイオク・レギュラーを給油しつつ、価格の最大値まで給油するみたいな仕組みを実装していました)

この意見は、実際にアプリケーションを作った上での実感にも近いように思います。
そこでこの記事では`redux`のレポジトリにあるサンプルアプリケーションを元にちょっと複雑な仕様を足してみて、`RxJS`の導入メリットが分かりやすくなることを目指したいと思います。

こんな手順で進めます。

1. `redux`のチュートリアルからサンプルアプリケーションを見繕ってコードと構造をザッと把握する
2. `1`で選んだサンプルアプリケーションの非同期処理部分を`RxJS`ベースの`redux`ミドルウェアである[redux-observable](https://github.com/redux-observable/redux-observable)で差し替える
3. 現実にあってもおかしくなさそうな仕様を足してみて、`RxJS`が複雑さを吸収できていそうか見てみる

## redux のチュートリアルにあるサンプルアプリケーションをクローンしてきてざっと見る

[redux レポジトリの examples ディレクトリ](https://github.com/reactjs/redux/tree/master/examples)にいくつかサンプルアプリケーションがあります。
この中で、非同期な通信の処理が入って現実のアプリケーションを想像しやすく、適度にシンプルな「ショッピングカート」アプリケーションを元にして進めてみたいと思います。

[最初はこんな状態です](]https://github.com/kogai/sample-redux-observable/tree/redux-thunk)
[https://www.evernote.com/l/AKiS90K33WFMOLWsEwDhEZSeLhagkl2xJkYB/image.png](アプリケーションの姿)

図のように、`Products`を在庫の許す限り`Cart`に投入できて、`Cart`の内容によって支払い総額が変わる、といったアプリケーションです。
非同期アクションを`redux-thunk`で`action creator`の中に dispatcher を流し込んで実現しているという、シンプルなアプリケーションとしてはわりとよくある感じの実装になっているようです。
※サーバ側の実装は煩雑さを避けるために、単なる json ファイルと setTimeout で、非同期にレスポンスが返ってくることだけを表現しています。

## サンプルアプリを、redux-observable で再実装してみる

では次に、このアプリケーションの`redux-thunk`を用いている部分を`redux-observable`に置き換えてみます。
まだこの時点では、`RxJS`を用いるメリットよりもデメリットの方が多いはずです。

[次はこういう状態にします](https://github.com/kogai/sample-redux-observable/tree/redux-observable)

```javascript
// actions/index.js
import * as types from "../constants/ActionTypes";

export const fetchAllProducts = () => ({
  type: types.FETCH_PRODUCTS
});

export const receiveAllProducts = products => ({
  type: types.RECEIVE_PRODUCTS,
  products
});

export const addToCart = productId => ({
  type: types.ADD_TO_CART_UNSAFE,
  productId
});

export const receiveInCart = productId => ({
  type: types.ADD_TO_CART,
  productId
});

export const checkout = products => ({
  type: types.CHECKOUT_REQUEST,
  products
});

export const checkoutSuccess = cart => ({
  type: types.CHECKOUT_SUCCESS,
  cart
});
```

```javascript
// epics/index.js
import { Observable } from "rxjs";
import shop from "../api/shop";
import * as types from "../constants/ActionTypes";
import { receiveAllProducts, receiveInCart, checkoutSuccess } from "../actions";

export const allProductsEpic = action$ =>
  action$
    .ofType(types.FETCH_PRODUCTS)
    .mergeMap(_ => Observable.bindCallback(shop.getProducts)())
    .map(receiveAllProducts);

export const addToCartEpic = (action$, store) =>
  action$
    .ofType(types.ADD_TO_CART_UNSAFE)
    .map(({ productId }) => productId)
    .filter(
      productId => store.getState().products.byId[productId].inventory > 0
    )
    .map(receiveInCart);

export const checkoutEpic = (action$, store) =>
  action$
    .ofType(types.CHECKOUT_REQUEST)
    .map(({ products }) => products)
    .mergeMap(products => Observable.bindCallback(shop.buyProducts)(products))
    .map(_ => store.getState().cart)
    .map(checkoutSuccess);
```

大きく変わっているのは、`redux-thunk`を用いて`action creator`の中で表現していた非同期処理を、`Epic`と呼ばれる別の構造の中で宣言することで、`action creator`が純粋で同期的な関数に戻っていることです。
`Epic`とは`redux-observable`が提供する構造で、私は`redux`の世界に`RxJS`の`Observable`を組み込むためのもの、みたいに解釈しています。
`redux-observable`のドキュメントには

> It is a function which takes a stream of actions and returns a stream of actions. Actions in, actions out.
> You can think of it of having roughly this type signature:
> `function (action$: Observable<Action>, store: Store): Observable<Action>;`

とあります。
`action`の`Observable`を受け取って、変性・フィルタリングなどを加えた`action`の`Observable`を返す関数である`Epic`を`Middleware`として`redux`の世界の中に組み込んでいる感じのようです。

## 仕様をいくつか足して、アプリ仕様を複雑にしてみて、Rx が複雑性を吸収する様を見る

さて、ここまでだとコードが減っているわけではないですし、あまりメリットが見えて来ないかも知れません。(`action`が純粋で同期的な関数からのみ生成されるようにはなりましたが。。。)
そこで(F)RP 的なメリットが見えやすいように、機能同士に複雑な依存関係が必要になるような仕様を追加してみたいと思います。

例えば今の仕様だと「カートの中身」は「商品の ID を投げる`action`」のみに依存しています。
これを別の`action`にも同時に依存しなければならないようにしてみたいと思います。

こんな仕様です。

- 3 個の商品を購入すると支払い総額から 1 割値引きされる
- 購入しようとしているユーザには通常会員・プレミアム会員の 2 種があり、属性によって割引率が変わる
- 購入した商品の総額によって、割引率が増える

こうなると、「カートの中身」とは「商品の ID を投げる Action」に加えて「ユーザの属性を投げる`action`」に依存したものになります。

実際にコードに起こしたものはこういう感じです。
[少しだけ複雑なアプリケーション](https://github.com/kogai/sample-redux-observable/tree/redux-obsrvable-complicated)

まず「カートの中身」を`Epic`にします。
`cart reducer`に定義されていた「カートの中身とは、ADD_TO_CART アクションで渡ってくる商品 ID をとりまとめたもの」という定義を`Epic`として再定義する感じです。
ですので、`Epic`の中に定義された関数は、`reducers/index.js`に定義されていたユーティリティ関数を少し改修したものが多いです。

```javascript
export const cartEpic = (action$, store) => {
  const productId$ = Observable.merge(
    action$.ofType(ADD_TO_CART).map(({ productId }) => productId),
    action$.ofType(CHECKOUT_REQUEST).mapTo(null)
  );

  const addedIds$ = productId$.scan((addedIds, id) => {
    if (id === null) {
      return [];
    }
    return addedIds.indexOf(id) !== -1 ? addedIds : [...addedIds, id];
  }, []);

  const quantityById$ = productId$.scan((quantityById, id) => {
    if (id === null) {
      return {};
    }
    return { ...quantityById, [id]: (quantityById[id] || 0) + 1 };
  }, {});

  return Observable.combineLatest(
    addedIds$,
    quantityById$,
    (addedIds, quantityById) => ({ addedIds, quantityById })
  ).map(({ addedIds, quantityById }) => {
    const { byId } = store.getState().products;

    const total = addedIds
      .reduce((acc, id) => acc + byId[id].price * (quantityById[id] || 0), 0)
      .toFixed(2);

    const products = addedIds.map(id => ({
      ...byId[id],
      quantity: quantityById[id] || 0
    }));

    return updateCart({
      total,
      products
    });
  });
};
```

次にアプリケーションにおける「ユーザー」を定義します。
まず既存の実装をベースに、偽の API レスポンス関数を用意します。

```json
// こんなレスポンスを返すAPIを想定しています
{
  "userType": 1, // 通常会員には0、プレミアム会員には1を返します
  "amount": 499 // そのユーザーの、今までの購買額を返します
}
```

```javascript
export const getUser = cb => setTimeout(() => cb(_user), TIMEOUT);
```

更にユーザー情報を`Epic`として定義して、ユーザー情報の取得(必要なら更新も)をアプリケーションに組み込みます。

```javascript
const userEpic = action$ =>
  action$
    .ofType(ON_LOAD)
    .mergeMap(_ => Observable.bindCallback(getUser)())
    .map(recieveUser);
```

最後に、割引率の算出ロジックを実装します。
まず`userEpic`がユーザー属性と購買総額に応じた割引率を取得できるように変更します。

```javascript
export const userEpic = action$ => {
  const user$ = action$
    .ofType(ON_LOAD)
    .mergeMap(_ => Observable.bindCallback(getUser)());

  const total$ = action$
    .ofType(CHECKOUT_REQUEST)
    .map(({ total }) => ({ amount: Number(total) }));

  return user$
    .merge(total$)
    .scan((acc, next) => ({
      amount: acc.amount + next.amount,
      userType: acc.userType
    }))
    .map(user => ({
      ...user,
      // プレミアムユーザー・支払い総額が1000ドルを超えるユーザーはそれぞれ割引率が追加される
      discountRate:
        (user.userType === UserTypes.PREMIUM ? 0.1 : 0) +
        (user.amount > 1000 ? 0.1 : 0)
    }))
    .map(recieveUser);
};
```

次に`userEpic`が生成する「割引率」を「カートの中身」に混ぜ込みます。
これで最終的な割引率を加味した購入額をカートのチェックアウトの際に取得できるようになります。

```javascript
export const cartEpic = (action$, store) => {
  // ...略
  const discountRate$ = action$.ofType(RECEIVE_USER).pluck("discountRate");

  return Observable.combineLatest(
    addedIds$,
    quantityById$,
    discountRate$,
    (addedIds, quantityById, userDiscountRate) => ({
      addedIds,
      quantityById,
      userDiscountRate
    })
  ).map(({ addedIds, quantityById, userDiscountRate }) => {
    // ...略
    const totalQuantity = addedIds.reduce(
      (acc, id) => acc + (quantityById[id] || 0),
      0
    );
    const discountRate = (totalQuantity >= 3 ? 0.3 : 0) + userDiscountRate;

    return updateCart({
      total,
      products,
      discountRate
    });
  });
};
```

という感じで、難しい書き方をせずに「複数の`action`を起点に別の`action`を発行する」ような機能が実装できるのが、`redux-observable`と`RxJS`を`redux`の世界で用いるメリットということになるでしょうか。
ここでは触れられませんでしたが、`Epic`は`Observable`を受け取って`Observable`を返す関数に過ぎないので、[RxJS の提供するテスト手法](http://blog.mmmcorp.co.jp/blog/2016/06/25/testing-rxjs-5/)に乗っかって簡単にテストできるのもうれしいところです。

## まとめ

以上、`redux-thunk`で非同期処理を取り扱っている`redux`の公式サンプルを`redux-observable`に置き換えて、少し仕様を複雑にするところまでを実装してみました。
何かの参考になればうれしいです。
明日の[React Advent Calendar 2016](http://qiita.com/advent-calendar/2016/react)もお楽しみに！(地図ライブラリの話めっちゃ気になります)

## 参考

- [RxJS のリードエンジニアによる redux-observable の紹介記事](https://medium.com/@benlesh/redux-observable-ec0b00d2eb52)
- [redux-observable のドキュメント兼チュートリアル](https://redux-observable.js.org/)
