---
title: "Promiseで非同期処理をするメモ"
published_at: "2015-04-13T19:20:20+09:00"
tags:
  - "node.js"
  - "promise"
  - "javascript"
  - "es6"
---

非同期処理の書き方色々

myAsync.js
```javascript
  'use strict';

  // Promiseオブジェクトをそのまま返す
  let readAsync = new Promise(function (resolve, reject) {
      fs.readFile('./README.md', function (err, data) {
        resolve(data);
      })
    });

  // Promiseオブジェクトを返す関数
  let writeAsync = function (data) {
    return new Promise(function (resolve, reject) {
      fs.writeFile('./README-copy.md', data, function (err) {
        resolve();
      });
    });
  }

  // Defferedオブジェクトのpromiseを返す
  let renameAsync = function () {
    var d = Promise.defer();
    fs.rename('./README-copy.md', './README-copy-modified.md', function (err) {
      d.resolve();
    });
    return d.promise;
  };

  // Generator
  let AppendAsync = function* (){
    fs.appendFile('./README-copy-modified.md', 'this string is appended.');
    yield;
  };

  readAsync
  .then(writeAsync)
  .then(renameAsync)
  .then(function () {
    AppendAsync().next();
  });
```

実行してみる
```sh
babel-node myAsync.js
```
