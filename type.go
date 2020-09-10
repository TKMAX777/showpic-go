package main

import "io"

//ImageReader wrap io.Reader
type ImageReader struct {
	reader io.Reader
	title  string
}

func (imgr *ImageReader) Read(p []byte) (n int, err error) {
	return imgr.reader.Read(p)
}

// New initialize struct
func (imgr *ImageReader) New(r io.Reader) {
	imgr.reader = r
}
