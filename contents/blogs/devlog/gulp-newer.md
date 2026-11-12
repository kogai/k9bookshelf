---
title: "gulp:開発/本番環境でタスクの内容を調整して幸せになる"
published_at: "2015-04-12T20:00:00+09:00"
tags:
  - "gulp"
  - "browser"
  - "javascript"
  - "coffeescript"
  - "npm"
  - "node.js"
---

### 更新したファイルだけをタスクの対象に

gulpでファイルの自動生成をする時、対象ファイルの数が増えてくると生成に時間がかかってツラい感じになってきます。
`gulp-newer`を使えば更新した差分ファイルのみをタスクの対象にしてくれるので幸せな感じになれます。

``` coffee
gulp = require 'gulp'
newer = require 'gulp-newer'
jade = require 'gulp-jade'

gulp.task 'sampleTask', ->
  # 更新差分の基準になるファイル
  criterionFiles = './dest/html'

	gulp.src([
		'./src/jade/!(_)*.jade'
	])
	.pipe(newer(criterionFiles))
	.pipe(jade())
	.pipe(gulp.dest(dest))

gulp.task 'default', [
  'sampleTask'
]
```

### 本番環境用ファイルを生成する時

本番環境用のファイルを作る時は、生成漏れがあったら困るので更新差分にはしたくありません。
そこで`production`環境の変数を作ってあげて、更新差分の基準ファイルを振り分けするようにしました。

``` coffee
config =
	production: false

gulp.task 'setProduction', ->
  config.production = true

gulp.task 'sampleTask', ->

  if( config.production )
	   criterionFiles = ''
  else
	   criterionFiles = './dest/html'

	gulp.src([
		'./src/jade/!(_)*.jade'
	])
	.pipe(newer(criterionFiles))
	.pipe(jade())
	.pipe(gulp.dest(dest))

gulp.task 'default', [
	'sampleTask'
]

gulp.task 'build', [
	'setProduction'
	'default'
]

```

### まとめ

静的サイトのジェネレータとしてgulpを使う事が多いので、100枚以上あるhtmlファイルの生成とか画像ファイルの生成にかかる時間が劇的に減りました。
※ちなみに、gruntにも同様のプラグイン(grunt-newer)があります。

#### 参考

[gulp-newer](https://www.npmjs.com/package/gulp-newer)
[How to handle development and production assets with gulp?](http://laravel.io/forum/04-03-2014-how-to-handle-development-and-production-assets-with-gulp)
