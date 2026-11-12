---
title: "circleCIのテストがnodev0.12で失敗する"
published_at: "2015-03-19T00:00:00+09:00"
tags:
  - "circleCI"
  - "node.js"
  - "npm"
---

`nodejs`の`v0.12`を使ってるプロジェクトでcircleCIのテストが通りませんでした。
手元のテストは成功しているのですが。。。

#### こんな感じのエラー

```sh
   Uncaught Error: Module did not self-register.
      at Error (native)
      at Module.load (module.js:355:32)
      at Function.Module._load (module.js:310:12)
      at Module.require (module.js:365:17)
      at require (module.js:384:17)
```

#### 困ったときのStackOverflow

[Uncaught Error: Module did not self-register](http://stackoverflow.com/questions/28486891/uncaught-error-module-did-not-self-register)

事前に
```sh
npm rebuild
```
すれば良いみたい。

`circle.yml`をこんな感じにして解決。
```yml
dependencies:
  override:
    - npm install
    - npm rebuild
```

ローカルマシンでも同様のエラーを吐く時があるけど、
同じように
```sh
npm rebuild
```
してあげれば解決するっぽいです。
