---
title: "Node.jsでReact.jsのサーバーサイドレンダリング"
published_at: "2015-04-10T00:00:00+09:00"
tags:
  - "node.js"
  - "server"
  - "browser"
  - "react.js"
  - "javascript"
---

Reactjsでwebアプリケーションを作る時に良いことの一つには、サーバーサイドレンダリングが自然な感じで出来るということがあるかと思います。

基本的に違いは2つだけみたいです。
- renderメソッドの代わりにrenderToString(特殊な状況ではrenderToStaticMarkUp)を使う
- 描画対象のDOMノードは指定しない

こんな感じで試してみました。
```javascript
var React = require('react');

var ServerSide = React.createClass({
  render: function () {
    'use strict';
    return (
      <h1>Hello! React-Server-Side-Rendering.</h1>
    );
  }
});

module.exports = React.renderToString(<ServerSide />);

// クライアントサイドだとこんな感じでDOMノードを指定してrenderしている
// React.render(<ServerSide />, document.body);
```

###gulp
jsxで書いていたら、jsにコンパイルしてあげる必要があります。
上のファイルを、例えばgulpでこんな感じにコンパイルしてあげます。

```javascript
gulp.task('default', function(){
  return gulp.src('エントリーポイントへのパス')
  .pipe(react())
  .pipe(gulp.dest('生成先のディレクトリ'));
});

```

###サーバー側
[Hapi](http://hapijs.com/)を使ってサーバーを立てています。

```javascript
var Hapi = require('hapi');
var server = new Hapi.Server();
var ServerSide = require('gulpで生成したファイルへのパス');

server.connection({
  host: 'localhost',
  port: 3000
});

server.route({
  method: 'GET',
  path: '/',
  handler: function ( req, reply ) {
    // React.renderToStringで生成したDOMノードをレスポンスとして返します。
    reply(ServerSide);
  }
});

server.start();
```

これならクライアントサイドの延長で記述できそうで素敵です。
