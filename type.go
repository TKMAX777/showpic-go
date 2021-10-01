package pic

import (
	"image"
	"io"
)

//ImageReader wrap io.Reader
type ImageReader struct {
	Title  string
	imgSrc image.Image
	imgDst image.Image
	rctSrc image.Rectangle
	zoom   int
	rate   float64
}

// Pos put position settings
type Pos struct {
	X int
	Y int
}

// New initialize struct
func (imgr *ImageReader) New(r io.Reader) error {
	// 画像を解析
	imgSrc, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	imgr.imgSrc = imgSrc
	imgr.rctSrc = imgSrc.Bounds()

	return err
}
