---
title: "Webデザイナーでも出来たGoogle Mapsを使わずに地図を表示させる方法"
published_at: "2014-08-22T14:01:04+09:00"
tags:
  - "javascript"
  - "openlayers"
---

## 前置き

1年ほど前に、政府がGoogle Mapsを使わないように要請したという事が話題になりました。
[ハフィントン・ポスト「グーグルマップの「利用禁止令」 竹島や北方領土が「日本名でない表記」 政府が自治体などに要請」](http://www.huffingtonpost.jp/2013/09/29/google-map-forbidden_n_4011180.html)

クライアントが政府筋の案件で地図を使わなくてはいけないなんて時に、手軽なGoogle Mapsが利用出来ないのは困りますね。
そんな時は[国土交通省-国土地理院](http://www.gsi.go.jp/)が提供している「地理院地図（電子国土Web.NEXT）」を使う事になるでしょう。

以前はサイト内で詳細な作成例を提供していたのですが、2014年7月にバージョンアップが行われた際に、
そこら辺のドキュメントもアーカイブされてしまったようです。
アーカイブ化前に実装した経験をもとに、「地理院地図APIから地図タイルを取得して、任意の座標に画像を表示させる」方法をポストしたいと思います。

------

## 手順

1. OpenLayersライブラリを読み込む
2. 「DenshiKokudo Web API for OpenLayers」を読み込む
3. 背景地図レイヤーを作成
4. 任意の座標に画像レイヤーを作成

------

## コード

```html
<script type="text/javascript" src="http://portal.cyberjapan.jp/sys/OpenLayers-2.11/OpenLayers.js" charset="UTF-8"></script>
<script type="text/javascript" src="http://portal.cyberjapan.jp/sys/v4/webtis/webtis_v4.js"></script>
<script>
  function init() {
    var map = null; //地図インスタンス
    initCX = 139; //初期の経度
    initCY = 38.5; //初期の緯度
    initZoomLv = 10; //初期のズームレベル
    projection900913 = new OpenLayers.Projection("EPSG:900913"); //真球メルカトル投影を定義
    projection4326 = new OpenLayers.Projection("EPSG:4326"); //等経緯度投影を定義
    selectControl; //SelectFeature
    maxExtent = new OpenLayers.Bounds(-20037508, -20037508, 20037508, 20037508); //真球メルカトル投影のときの最大範囲(単位はm)
    restrictedExtent = maxExtent.clone(); //真球メルカトル投影のときの最大範囲に範囲を制限
    maxResolution = 156543.0339; //真球メルカトル投影のときの最大解像度
    icon = '表示させる画像要素のパス';

    //地図表示画面のオプション設定
    var options = {
      controls: [ //「controls」を設定することで、デフォルトのコントロールを破棄してコントロールを再設定
        new OpenLayers.Control.Navigation({
          mouseWheelOptions: {
            interval: 100
          }
        }), //地図マウスイベントのハンドル設定。
        new OpenLayers.Control.PanZoomBar(), //左上のパンズームバーを設定
        new OpenLayers.Control.KeyboardDefaults(), //キーボードをデフォルトに設定
        new OpenLayers.Control.Attribution() //国土地理院著作表示
      ],
      projection: projection900913, //背景地図の地理座標系
      displayProjection: projection4326, //表示の地理座標系
      units: "m", //背景地図の単位
      maxResolution: maxResolution, //背景地図の最大解像度
      maxExtent: maxExtent, //背景地図の最大範囲
      restrictedExtent: restrictedExtent //背景地図の表示制限範囲
    };

    //OpenLayers APIのMapクラスからインスタンスを作成
    map = new OpenLayers.Map('map', options);

    //スケールバーコントロール表示(最大ピクセル150、下段単位無、EPSG:900913)
    map.addControl(new OpenLayers.Control.ScaleLine({
      maxWidth: 150,
      bottomOutUnits: "",
      bottomInUnits: "",
      geodesic: true
    }));

    //電子国土WebシステムVer.4背景地図レイヤーインスタンスを作成。データセットは未指定で、デフォルトデータセットを利用
    webtisMap = new webtis.Layer.BaseMap("webtismap");

    //背景地図レイヤーをMapに追加
    map.addLayer(webtisMap);


    //***ここからベクターレイヤーの記述***

    //ベクターレイヤーの作成
    myVectorLayer = new OpenLayers.Layer.Vector("layername");

    //ベクターレイヤーの座標を定義 (for文を回すと、複数の座標にアイコンを定義できます)
    myVectorLayer.addFeatures([new OpenLayers.Feature.Vector(
      new OpenLayers.Geometry.Point(initCX, initCY).transform(projection4326, projection900913)
      // new OpenLayers.Geometry.Point(経度,緯度).transform(projection4326,projection900913)
      // 緯度と経度を任意の数値に変えると、初期表示位置とは別の場所にアイコン画像を表示できる
    )]);

    vector_style = new OpenLayers.Style({
      'externalGraphic': icon,
      'pointRadius': 16
    });

    myVectorLayer.styleMap = new OpenLayers.StyleMap({
      'default': vector_style
    });

    //ベクターレイヤーをMapに追加
    map.addLayer(myVectorLayer);

    //***ここまでベクターレイヤーの記述***


    // MapにSelectFeatureコントロールを追加
    selectControl = new OpenLayers.Control.SelectFeature(myVectorLayer);
    map.addControl(selectControl);
    selectControl.activate();

    //初期の中心座標を指定（経緯度で入力して、内部的に真球メルカトル座標に変換して表示）
    map.setCenter(new OpenLayers.LonLat(initCX, initCY).transform(projection4326, projection900913), initZoomLv);
  }

  init();
</script>

<div id="map"></div>
```
