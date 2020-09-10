# Go TLI Pic
## 概要
端末上で画像を表示するためのプログラム。

## Install

```sh
go get github.com/gdamore/tcell
go get -u github.com/mattn/go-runewidth
go build
```

24bitで描画するためには、環境変数に次を追加する必要があります。

```sh
export TERM='xterm-256color'
```
