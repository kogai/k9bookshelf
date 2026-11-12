---
title: "なぜタスクランナーを使ったほうがいいのかをデザイナーに伝えるためにまとめた"
published_at: "2015-05-13T19:54:03+09:00"
tags:
  - "gulp.js"
  - "grunt.js"
  - "sass"
  - "jade"
  - "javascript"
---

近年のフロントエンド制作現場では grunt, gulp などのタスクランナーや、
browserify, webpack などのビルドツールを使用するのが半ば常識となっています。
なぜこういったツールを使うのか、社内のデザイナーに伝えるためにまとめてみました。

## タスクランナーとは?

現在要求されるレベルで成果物のクオリティを担保するためには、以下のような多岐にわたる処理が必要になります。

_必要になる処理の例_

- テンプレートエンジンによる HTML ファイルの生成
- CSS プリプロセッサによる CSS の生成&圧縮
- CSS スプライトの作成
- JavaScript ファイルの依存関係解決&圧縮
- Alt-JS ファイルの変換
- 各種画像の最適化
- 他多数...

これら全てを手作業で行うのは、非現実的な選択肢と言えます。
そこで、こういったタスクを全て自動化するためのツールが`grunt`/`gulp`に代表される、タスクランナーと呼ばれるツールです。

### ビルドツール

極々小規模なプロジェクト(数十行で完結する程度のもの)を除いて、`browserify`, `webpack`などのビルドツールを用いて
`JavaScript`ファイルを生成する事が、スタンダードになってきています。

これには以下のような理由があります。

- `CoffeeScript`や`TypeScript`など、`AltJS`の普及によって、実行用の`JavaScript`ファイルへの変換が必要になった
- 既存の`JavaScript`にはモジュール機構が存在しないため、全ての関数がグローバルになってしまう問題があった(※2015 年夏頃より勧告される ES6 では`import`文によるモジュール機構が実装される)
- モジュール機構が存在しないため、サードパーティ製ライブラリに対する依存状態が管理できず、運用が進むにつれてライブラリの要・不要が判別できなくなることが多かった
- モジュール機構が存在しないため、手作業によるライブラリの読み込みの設定を余儀なくされた。そのため読み込み順によっては動作しないコードが出てくる問題があった

こういった問題を解決するために、`browserify`, `webpack`に代表される、ビルドツールが用いられるようになりました。
ビルドツールによるビルド処理も、タスクランナーを用いて行われます。

## 導入
`grunt`, `gulp`は共に、`node.js`というサーバーサイド`JavaScript`実行環境(PC で直接`JavaScript`を実行出来るようになるソフトウェア、と考えれば大体正しいです)を用いて実装されています。
`grunt`, `gulp`のどちらも、使い方や導入方法、出来ることに大きな差はありません。
各タスクの処理の仕方の差で、`gulp`の方が動作が速いので、新規のプロジェクトであれば`gulp`を使うのをお勧めします。

詳しい導入方法を説明した記事は無数にあるので、`grunt 導入`とか`grunt how to use`などでググって下さい。

以下なんかはわりと詳しく載っていてお勧めです。
[Web デザイナーさん向け Grunt を使った コーディング作業の効率化、はじめの一歩](http://hyper-text.org/archives/2014/01/grunt_quick_start_for_web_designer.shtml)

### 大体の手順

1. `node.js`をインストール
2. `grunt-cli`(Grunt コマンドラインクライアント)をインストール
3. 各プロジェクトで Grunt タスクを実行可能に

一応これで伝わるかな。。。
ついでに普段使ってる`gulpfile.js`の雛形も置いておきます

```coffeescript
gulp = require 'gulp'
browserSync = require 'browser-sync'
reload = browserSync.reload
data = require 'gulp-data'
util = require('gulp-util')
debug = require 'gulp-debug'
newer = require 'gulp-newer'
compass = require 'gulp-compass'
sourcemaps = require 'gulp-sourcemaps'
minify = require 'gulp-minify-css'
del = require('del')

# browserify
browserify = require 'browserify'
debowerify = require 'debowerify'
licensify = require 'licensify'
source = require 'vinyl-source-stream'
streamify = require 'gulp-streamify'
uglify = require 'gulp-uglify'
buffer = require 'vinyl-buffer'

# jade
jade = require 'gulp-jade'

# images
pngmin = require 'gulp-pngmin'
imagemin = require 'gulp-imagemin'
jpegtran = require 'imagemin-jpegtran'

# watch
watch = require 'gulp-watch'

# deploy
scp = require 'scp'

src =
  js : './src/javascript'
  css : './src/sass'
  html : './src/jade'
  image : './src/image'
  copy : './src/copy'

build =
  js : './build/public'
  css : './build/public'
  html : './build'
  image : './build/public/image'

env =
  isProduct : false
  deployToProduct : false
  rootDev: '/path/to/staging'
  rootPro: '/path/to/production'

gulp.task 'deploy', ->
  dest = env.rootDev

  if env.deployToProduct
    dest = env.rootPro
  else
    dest = env.rootDev

  scp.send(
    file : './deploy/*'
    host : 'hostname'
    user : 'username'
    path : dest
  , (e) ->
    util.beep()
    console.log e
  )

gulp.task 'changeEnv', ->
  env.isProduct = true
  build =
    js : './deploy/public'
    css : './deploy/public'
    html : './deploy'
    image : './deploy/public/image'

gulp.task 'changeDeployEnv', ->
  env.deployToProduct = true

gulp.task 'css', ->
  if env.isProduct
    gulp.src(src.css + '/*.sass')
    .pipe compass(
      css : build.css
      sass : src.css
    )
    .pipe minify()
    .pipe gulp.dest(build.css)
  else
    gulp.src(src.css + '/*.sass')
    .pipe sourcemaps.init()
    .pipe compass(
      css : build.css
      sass : src.css
    )
    .pipe sourcemaps.write('./')
    .pipe gulp.dest(build.css)

gulp.task 'html', ->
  if env.isProduct
    opt = pretty : false
  else
    opt = pretty : true

  gulp.src(src.html + '/*.jade')
  .pipe data (file) ->
    return require src.html + '/index.json'
  .pipe(jade(opt))
  .pipe gulp.dest(build.html)
  .on 'error', (e) ->
    util.beep();
    console.log e
  return

gulp.task 'js', ->
  if env.isProduct
    browserify
      entries : [src.js + '/index.coffee']
      extensions : ['.coffee', '.js']
    .plugin licensify
    .transform 'coffeeify'
    .transform 'debowerify'
    .bundle()
    .pipe source('bundle.min.js')
    .pipe buffer()
    .pipe streamify uglify()
    .pipe gulp.dest(build.js)
  else
    browserify
      entries : [src.js + '/index.coffee']
      extensions : ['.coffee', '.js']
    .plugin licensify
    .transform 'coffeeify'
    .transform 'debowerify'
    .bundle()
    .pipe source('bundle.min.js')
    .pipe buffer()
    .pipe sourcemaps.init
      loadMaps : true
    .pipe streamify uglify()
    .pipe sourcemaps.write('./')
    .pipe gulp.dest(build.js)

gulp.task 'pngmin', ->
  gulp.src([
    src.image + '/*.png'
    src.image + '/**/*.png'
  ])
  .pipe newer build.image
  .pipe pngmin()
  .pipe gulp.dest build.image

gulp.task 'jpgmin', ->
  gulp.src([
    src.image + '/*.jpg'
    src.image + '/**/*.jpg'
  ])
  .pipe imagemin(
    prpgressive : true
    svgoPlugins : [{removeViewBox : false}]
    use : [jpegtran()]
  )
  .pipe gulp.dest build.image

gulp.task 'copy-dir', ->
  gulp.src([
    src.copy + "/*"
    src.copy + "/**/*"
    src.copy + "/**/**/*"
  ])
  .pipe(gulp.dest(build.html))

gulp.task 'copy-img', ->
  gulp.src([
    src.image + '/*.gif'
    src.image + '/**/*.gif'
    src.image + '/**/**/*.gif'
  ])
  .pipe(gulp.dest(build.image))

gulp.task 'img', [
  'pngmin'
  'jpgmin'
  'copy-img'
]

gulp.task 'browser', ->
  browserSync server :
    baseDir : 'build'
    directory : true
  return

gulp.task 'clean', (cb) ->
  del([
    'deploy/'
  ], cb)

gulp.task 'watch', ->
  gulp.watch([
      src.js + '/*.coffee'
      src.js + '/**/*.coffee'
    ], [
      'js'
      reload
  ])
  gulp.watch([
      src.html + '/*.jade'
      src.html + '/**/*.jade'
    ], [
      'html'
      reload
  ])
  gulp.watch([
      src.css + '/*.sass'
      src.css + '/**/*.sass'
    ], [
      'css'
      reload
  ])
  gulp.watch([
      src.image + '/*'
      src.image + '/**/*'
    ], [
      'img'
      reload
  ])
  return

gulp.task 'default', [
  'css'
  'js'
  'html'
  'img'
  'copy-dir'
]

gulp.task 'server', [
  'default'
  'browser'
  'watch'
]

gulp.task 'build', [
  'clean'
  'changeEnv'
], ->
  gulp.run 'default'

gulp.task 'staging', [
  'deploy'
]

gulp.task 'production', [
  'changeDeployEnv'
  'deploy'
]

```
