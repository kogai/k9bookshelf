---
title: "vue.jsでjqueryのpluginを使ってみる"
published_at: "2015-03-09T13:51:56+09:00"
tags:
  - "vue.js"
  - "javascript"
  - "jQuery"
---

業務でIE9〜の案件が降ってきたので、vue.jsを使ってみた。
よくあるスライド切替をvue.js + jQuery + jQuery.bxSliderで。

## スライド画像が繰り返される部分に`v-repeat`を設定

```jade
ul(id='slide')
  li(v-repeat='photos')
    img(src="{{photo}}", alt="{{title}}")
```

## `vue.js`で要素を描画してからjQueryプラグインを適用

~~~js
var slide = new Vue({
  el: '#slide',
  data: {
    photos: [
      {
        "photo": "画像1のパス",
        "title": "画像1の名前"
      },
      {
        "photo": "画像2のパス",
        "title": "画像2の名前"
      }
    ]
  }
});
jQuery("#slide").bxSlider();
~~~

これで行けました。
