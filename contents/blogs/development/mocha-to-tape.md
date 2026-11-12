---
title: "TypeScript製ライブラリのテストフレームワークをmochaからtapeに切り替えた"
published_at: "2016-12-15T12:35:04+09:00"
tags:
  - "TypeScript"
  - "JavaScript"
  - "mocha"
  - "tape"
---

普段テストは、Node.jsのコードはmocha・ブラウザ向けのコードはkarma(jasmine)で書くことが多いんだけど、たまには別のライブラリも試してみたいな〜と思ってたところにこの記事を読んだ

[Why I use Tape Instead of Mocha & So Should You](https://medium.com/javascript-scene/why-i-use-tape-instead-of-mocha-so-should-you-6aa105d8eaf4#.yaaii9wua)

1年半くらい前の記事で、要約すると設定ファイル不要でグローバルな変数やテスト間で共有される変数を持たないシンプルなテストフレームワークなのでMochaよりtapeに乗り換えた、テストは極力複雑にならないようにして(その方がデバッグも容易、みたいな文脈もある気がする)もっとプロダクトの質に注力しよう、みたいな感じか。
テストフレームワークやらランナーやらの設定でごちゃごちゃなってるなというのは、前々から感じていたことではあったので、自分で公開している小さなライブラリのmochaを使って書いていたテストコードをtapeに切り替えてみた。
前に書いた[このポスト](/blogs/development/create-custom-operator)で作ったやつ。

元々のコードはこんな感じ
RxJSのmarbleを使ってテストしてるのでちょっと論旨がズレそうな気もする...

```typescript
import { deepStrictEqual } from "assert";
import { TestScheduler } from "rxjs/testing/TestScheduler";
import { Observable } from "rxjs/Observable";
import "./";
import {
  createTestScheduler,
  TestSchedulers,
  createColdObservable
} from "./test-helpers";
  
describe("ofTypeObservable", () => {
  let action$: Observable<any>;
  let testSchedulers: TestSchedulers;
  let testScheduler: TestScheduler;
  let cold: createColdObservable;
  
  beforeEach(() => {
    // createTestSchedulerの中でassertionをしている
    testSchedulers = createTestScheduler();
    testScheduler = testSchedulers.testScheduler;
    cold = testSchedulers.cold;
    action$ = cold("abc", {
      a: { type: "foo", payload: "fooPayload" },
      b: { type: "bar", payload: "barPayload" },
      c: { type: "buzz", payload: "buzzPayload" },
    });
  });
  
  it("should filtering specific type", () => {
    const expect$ = action$.ofType("bar");
    testScheduler.expectObservable(expect$).toBe("-a", { a: "barPayload" });
    testScheduler.flush();
  });
  
  it("can pick by own payload", () => {
    action$ = cold("ab", {
      a: { type: "foo", payload: "fooPayload" },
      b: { type: "bar", customPayload: "barPayload" },
    });
    const expect$ = action$.ofType("bar", (value: any) => value.customPayload);
    testScheduler.expectObservable(expect$).toBe("-a", { a: "barPayload" });
    testScheduler.flush();
  });
});
```

変更後はこんな感じ

```typescript
import * as test from "tape";
import { deepStrictEqual } from "assert";
import { TestScheduler } from "rxjs/testing/TestScheduler";
import { Observable } from "rxjs/Observable";
  
import "./";
import { createTestScheduler, TestSchedulers, createColdObservable } from "./test-helpers";
  
const setup = (assert: test.Test) => {
  // createTestSchedulerの中でassertionをしているので、tapeの提供しているassertオブジェクトを渡してる
  const testSchedulers = createTestScheduler(assert);
  const testScheduler = testSchedulers.testScheduler;
  const cold = testSchedulers.cold;
  const action$ = cold("abc", {
    a: { type: "foo", payload: "fooPayload" },
    b: { type: "bar", payload: "barPayload" },
    c: { type: "buzz", payload: "buzzPayload" },
  });
  return {
    testSchedulers, testScheduler, cold, action$,
  }
}
  
test("should filtering specific type", assert => {
  const {action$, testScheduler} = setup(assert)
  const expect$ = action$.ofType("bar");
  testScheduler.expectObservable(expect$).toBe("-a", { a: "barPayload" });
  testScheduler.flush();
  assert.end()
})
  
test("can pick by own payload", assert => {
  const {cold, testScheduler} = setup(assert)
  const action$ = cold("ab", {
    a: { type: "foo", payload: "fooPayload" },
    b: { type: "bar", customPayload: "barPayload" },
  });
  const expect$ = action$.ofType("bar", (value: any) => value.customPayload);
  testScheduler.expectObservable(expect$).toBe("-a", { a: "barPayload" });
  testScheduler.flush();
  assert.end()
})
```

コード全体は[ココ](https://github.com/kogai/of-type-operator)にある

<iframe style="width:120px;height:240px;" marginwidth="0" marginheight="0" scrolling="no" frameborder="0" src="//rcm-fe.amazon-adsystem.com/e/cm?lt1=_blank&bc1=000000&IS2=1&bg1=FFFFFF&fc1=000000&lc1=0000FF&t=kogai-22&o=9&p=8&l=as4&m=amazon&f=ifr&ref=as_ss_li_til&asins=B00ESXY9MA&linkId=b69fd1b67eac33ba21c756c87b0e0380"></iframe>

たしかこの本で読んだんだったと思うけど、テストフレームワークの目的は入力に対して期待した出力が得られているか確認する単なる関数に過ぎない、みたいな話があって、自前のassertion関数を実装するように促されたりしていた。
`tape`はそういう思想にわりと近い感じがしていて、登場人物がとても少なく且つグローバルな変数を使っていたりしないので、テストファイルの中で全てが完結しているのが良い感じだと思った。

テストリポーターを任意のものにスイッチするのが

```bash
$(npm bin)/tape test.js | $(npm bin)/tape-custom-reporter
```

みたいにパイプで流せるのも良い。

プロジェクト毎に色々変わったことをしようとするとまた話が変わってくるのかも知れないが、小さなライブラリのために使うのであれば、軽く導入できて良いように思った。
あとブラウザでテストを実行するのは試してみてないので、それは別途試してみたい。
