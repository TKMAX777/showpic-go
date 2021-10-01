package pic

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

// Screen put screen
var Screen tcell.Screen

//PutRow put current head of the page
var PutRow = 0

// Init initialize screen
func Init() {
	os.Setenv("TERM", "xterm-256color")

	var err error
	// TLI画面の生成
	Screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err = Screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	Screen.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite),
	)
	Screen.Clear()

	return
}

// Putln output put data to screen
func Putln(s tcell.Screen, style tcell.Style, d ...interface{}) {
	Puts(s, style, 1, PutRow, fmt.Sprint(d...))
	PutRow++
}

// Puts output data to put point and return the last coordinate
func Puts(s tcell.Screen, style tcell.Style, x, y int, d ...interface{}) (int, int) {
	var i int = 0

	var str string = fmt.Sprint(d...)
	var deferred []rune

	var dwidth = 0
	var zwj = false

	for _, r := range str {
		if r == '\u200d' {
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
			deferred = append(deferred, r)
			zwj = true
			continue
		}
		if zwj {
			deferred = append(deferred, r)
			zwj = false
			continue
		}
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}

	return x + i, y
}

// PutAln show text on put coordinate and apply the style to whole line
func PutAln(s tcell.Screen, style tcell.Style, X, Y int, d ...interface{}) {
	width, _ := s.Size()

	for x := 0; x < width; x++ {
		s.SetContent(x, Y, ' ', nil, style)
	}

	Puts(s, style, X, Y, d...)

	return

}
