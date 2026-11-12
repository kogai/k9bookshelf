---
title: "object-assignの挙動について勘違いしていたこと"
published_at: "2015-06-04T19:48:55+09:00"
tags:
  - "node.js"
  - "javascript"
  - "npm"
---

仕様をきちんと理解していれば自明のことですが、ハマりかけたのでメモしておきます。

## この記事の概要
複数のオブジェクトを合成して新しいオブジェクトを返してくれる便利機能`Object.assign`。
現在は`harmony`オプション付きのNode.jsでも使えないので、`object-assign`などのモジュールを経由して使っていると思います。

私が勘違いしていたのですが、これは
>引数として渡したオブジェクトを合成したオブジェクトを返す

という機能ではなく、
>第二引数以降に渡したオブジェクトを、第一引数に渡したオブジェクトに合成して返す

ものなのですね。

MDNに記載されているリファレンスにも
>１つ以上のソースオブジェクトの保有する全ての列挙プロパティの値を、ターゲットのオブジェクトへコピーします。
>戻り値はターゲットオブジェクトになります。

と、明記されています。 [Object.assign()](https://developer.mozilla.org/ja/docs/Web/JavaScript/Reference/Global_Objects/Object/assign)

コードで表すと
```javascript
"use strict";

var objectAssign = require('object-assign');

var obj = {
  foo: "bar"
};

var obj2 = {
  baz: "foobar"
};

console.log(obj); // { foo: 'bar' }
console.log(obj2); // { baz: 'foobar' }

var obj3 = objectAssign(obj, obj2);

console.log(obj3); // { foo: 'bar', baz: 'foobar' }
console.log(obj); // { foo: 'bar', baz: 'foobar' } objの中身も変わっている
```
`obj3`に返されているのは`obj2`を合成した`obj`ですので、当然の挙動です。
`obj`の中身を変更したくないのであれば、次のような書き方をする必要があります。

```javascript
var target = {};
var obj4 = objectAssign(target, obj, obj2);

console.log(obj4); // { foo: 'bar', baz: 'foobar' }
console.log(obj); // { foo: 'bar' } objの中身は変わっていない
```

`React.js`で作っているサイトで、`Inline Styles`を導入してみて、初めてこの挙動を知った次第です。
きちんと仕様を理解するのが大事、というお話でした。
