# Go TLI Pic
## 概要
端末上で画像を表示するためのプログラム。

![例](https://i.gyazo.com/7e182aaebdbcbd459f134c0e21c99947.gif)

## 目次
<!-- TOC -->

- [Go TLI Pic](#go-tli-pic)
    - [概要](#概要)
    - [目次](#目次)
    - [Install](#install)
    - [キー操作](#キー操作)
    - [うまく行かないとき](#うまく行かないとき)
        - [コマンドが見つからない](#コマンドが見つからない)
        - [発色がおかしい](#発色がおかしい)

<!-- /TOC -->

## Install

```sh
go install github.com/TKMAX777/showpic-go/cmd/show@latest
show IMAGE.png
```

## キー操作
- `Enter`
  - 次の画像へ進みます
- `BackSpace`
  - 前の画像に戻ります
- `Ctrl + I`
  - 設定を初期化し、再読み込みをします。
- `Ctrl + L`
  - 再読み込みします。
- `+` / `-`
  - 拡大/縮小します。
- `↓↑←→`
  - 画像を移動します。
- `q`/`Esc`
  - ビューワを閉じます

## うまく行かないとき
### コマンドが見つからない
GOBINにパスが通っている必要があります。

- GOBINを設定していない場合

```sh
export PATH=$GOPATH/bin:$PATH
```

- GOBINを設定している場合

```sh
export PATH=$GOBIN/bin:$PATH
```

### 発色がおかしい

24bitで描画するためには、環境変数に次を追加する必要がある場合があります。

```sh
export TERM='xterm-256color'
```
