# Go TLI Pic
## 概要
端末上で画像を表示するためのプログラム。

## Install

```sh
go install github.com/TKMAX777/showpic-go/cmd/show@latest
```

24bitで描画するためには、環境変数に次を追加する必要があります。

```sh
export TERM='xterm-256color'
```


### キー操作
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