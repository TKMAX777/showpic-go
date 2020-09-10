package main

import (
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
func (imgr ImageReader) Set(sW, sH int) (err error) {
	img, _, err := image.Decode(&imgr)
	if err != nil {
		return err
	}

	var rctSrc = img.Bounds()

	{
		// 100以上の大きさを受け付けない
		if sW > 100 || sH > 100 {
			sW, sH = 0, 0
		}

		sH *= 2

		// 画面の外枠に対する大きさを指定する
		switch sW {
		case 0:
			switch sH {
			case 0:
				// 指定のない場合
				// 画像の大きさで適宜大きさを画面に合わせる
				switch {
				case rctSrc.Dx() <= rctSrc.Dy():
					_, height := Screen.Size()
					switch imgr.title {
					case "":
						sH = int(height+1) * 2
					default:
						sH = int(height) * 2
					}
					sW = int(float64(rctSrc.Dx()) / float64(rctSrc.Dy()) * float64(sH))
				default:
					width, height := Screen.Size()
					sW = int(width + 1)
					sH = int(float64(rctSrc.Dy())/float64(rctSrc.Dx())*float64(sW)) * 2
					switch imgr.title {
					case "":
						if height+1 < sH {
							sH = int(height+1) * 2
							sW = int(float64(rctSrc.Dx()) / float64(rctSrc.Dy()) * float64(sH))
						}
					default:
						if height < sH {
							sH = int(height) * 2
							sW = int(float64(rctSrc.Dx()) / float64(rctSrc.Dy()) * float64(sH))
						}

					}

				}

			default:
				// 縦幅のみ指定のある場合
				_, height := Screen.Size()
				switch imgr.title {
				case "":
					sH = int(float64(height+1) * float64(sH) / 100.0)
				default:
					sH = int(float64(height) * float64(sH) / 100.0)
				}
				sW = int(float64(rctSrc.Dx()) / float64(rctSrc.Dy()) * float64(sH))
			}
		default:
			// 横幅のみ指定のある場合
			width, _ := Screen.Size()
			sW = int(float64(width+1) * float64(sW) / 100.0)
			sH = int(float64(rctSrc.Dy()) / float64(rctSrc.Dx()) * float64(sW))
		}

		sH /= 2

	}

	// 画像のサイズを変更
	var imgDst *image.RGBA
	imgDst = image.NewRGBA(image.Rect(0, 0, sW, sH))
	draw.CatmullRom.Scale(imgDst, imgDst.Bounds(), img, img.Bounds(), draw.Over, nil)

	// 画像の解析
	var bounds = imgDst.Bounds()

	width, height := Screen.Size()

	var minX, maxX int = bounds.Min.X, bounds.Max.X
	var minY, maxY int = bounds.Min.Y, bounds.Max.Y

	// タイトル行の用意
	if imgr.title != "" {
		minY++
		maxY++
	}

	// 配置
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			// 色をつけるオブジェクトの作成
			var st = tcell.StyleDefault

			// 16bitから24bitに変換
			const rate = float64(256) / float64(65536)

			// 色を取得
			r, g, b, _ := imgDst.RGBAAt(x, y).RGBA()

			st = st.Background(tcell.NewRGBColor(
				int32(float64(r)*rate),
				int32(float64(g)*rate),
				int32(float64(b)*rate),
			))

			var X = (width-bounds.Max.X+bounds.Min.X)/2 + x
			var Y = (height-bounds.Max.Y+bounds.Min.Y)/2 + y

			Screen.SetContent(X, Y, ' ', nil, st)
		}
	}

	return
}

// SetTitle show image title on top of the screen
func (imgr *ImageReader) SetTitle(title string) {
	imgr.title = title

	// スタイル設定生成
	var st = tcell.StyleDefault

	st = st.Background(tcell.ColorWhiteSmoke).Foreground(tcell.ColorBlack)

	PutAln(Screen, st, 1, 0, title)

	return
}
