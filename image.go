package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/gdamore/tcell"
	"golang.org/x/image/draw"
)

// Set set screen pixcels
func (imgr ImageReader) Set(delta Pos) {
	// 画面の大きさを取得
	width, height := Screen.Size()
	var bounds = imgr.imgDst.Bounds()

	var minX, maxX int = bounds.Min.X, bounds.Max.X
	var minY, maxY int = bounds.Min.Y, bounds.Max.Y

	var titleDel int
	if imgr.Title != "" {
		titleDel = 1
	}

	// 配置
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			// 色をつけるオブジェクトの作成
			var st = tcell.StyleDefault

			// 16bitから24bitに変換
			const rate = float64(256) / float64(65536)

			// 色を取得
			r, g, b, _ := imgr.imgDst.RGBAAt(x, y).RGBA()

			st = st.Background(tcell.NewRGBColor(
				int32(float64(r)*rate),
				int32(float64(g)*rate),
				int32(float64(b)*rate),
			))

			var X = (width-maxX+minX)/2 + x + delta.X
			var Y = (height-maxY+minY)/2 + y + delta.Y + titleDel

			Screen.SetContent(X, Y, ' ', nil, st)
		}
	}

	return
}

// Zoom zoom picture at put rate
func (imgr *ImageReader) Zoom(rate float64) {

	var w, h = float64(imgr.rctSrc.Dx()), float64(imgr.rctSrc.Dy())

	// 画像のサイズを変更
	var imgDst *image.RGBA
	imgDst = image.NewRGBA(image.Rect(0, 0, int(w*rate*imgr.rate*2), int(h*rate*imgr.rate)))
	draw.CatmullRom.Scale(imgDst, imgDst.Bounds(), imgr.imgSrc, imgr.rctSrc, draw.Over, nil)

	imgr.imgDst = imgDst

	return
}

func (imgr *ImageReader) getSuitRate(sW, sH int) float64 {
	var rctSrc = imgr.rctSrc
	var rate float64

	// 画面の大きさを取得
	width, height := Screen.Size()
	// 100以上の大きさを受け付けない
	if sW > 100 || sH > 100 {
		sW, sH = 0, 0
	}

	// 画面の外枠に対する大きさを指定する
	switch sW {
	case 0:
		switch sH {
		case 0:
			// 指定のない場合
			// 画像の大きさで適宜大きさを画面に合わせる
			switch {
			case rctSrc.Dx() <= rctSrc.Dy():
				switch imgr.Title {
				case "":
					// タイトル無し
					// 縦幅を基準に設定
					sH = int(height + 1)
					rate = float64(sH) / float64(rctSrc.Dy())
					sW = int(float64(rctSrc.Dx()) * rate)
					// 横幅がはみ出てしまう場合
					if width+1 < sW*2 {
						sW = int(width + 1)
						rate = float64(sW) / float64(rctSrc.Dx()) / 2
					}

				default:
					// タイトル有り
					// 縦幅を基準に設定
					sH = int(height)
					rate = float64(sH) / float64(rctSrc.Dy())
					sW = int(float64(rctSrc.Dx()) * rate)

					// 横幅がはみ出てしまう場合
					if width < sW*2 {
						sW = int(width)
						rate = float64(sW) / float64(rctSrc.Dx()) / 2
					}
				}
			default:
				// 横幅を基準に設定
				sW = int(width + 1)
				rate = float64(sW) / float64(rctSrc.Dx()) / 2

				sH = int(float64(rctSrc.Dy()) * rate)

				switch imgr.Title {
				case "":
					// タイトル無し
					// 縦幅がはみ出てしまう場合
					if height+1 < sH {
						sH = int(height + 1)

						rate = float64(sH) / float64(rctSrc.Dy())
					}
				default:
					// タイトル有り
					// 縦幅がはみ出てしまう場合
					if height < sH {
						sH = int(height)

						rate = float64(sH) / float64(rctSrc.Dy())
					}

				}
			}
		default:
			// 縦幅のみ指定のある場合
			switch imgr.Title {
			case "":
				sH = int(float64(height+1) * float64(sH) / 100.0)
			default:
				sH = int(float64(height) * float64(sH) / 100.0)
			}
			rate = float64(sH) / float64(rctSrc.Dy())
		}
	default:
		// 横幅のみ指定のある場合
		sW = int(float64(width+1) * float64(sW) / 100.0)
		rate = float64(sW) / float64(rctSrc.Dx()) / 2
	}

	imgr.rate = rate

	return rate
}

// SetTitle show image title on top of the screen
func (imgr *ImageReader) SetTitle(ZoomRate float64) {

	// スタイル設定生成
	var st = tcell.StyleDefault

	st = st.Background(tcell.ColorWhiteSmoke).Foreground(tcell.ColorBlack)

	PutAln(Screen, st, 1, 0, fmt.Sprintf("%s 拡大率:%3.1f%%", imgr.Title, ZoomRate*100))

	return
}
