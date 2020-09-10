package main

import (
	"flag"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"time"

	"github.com/gdamore/tcell"
)

func init() {
	Init()
	flag.ErrHelp = fmt.Errorf("pic PICTURE")

}

func main() {

	// チャネルの生成
	var quit = make(chan struct{})
	var next = make(chan bool)
	var resize = make(chan bool)
	var back = make(chan bool)

	go func() {
		for {
			ev := Screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					close(quit)
					return
				case tcell.KeyEnter:
					next <- true
				case tcell.KeyCtrlL:
					Screen.Sync()
					resize <- true
				case tcell.KeyBackspace, tcell.KeyBackspace2:
					back <- true
				}
			case *tcell.EventResize:
				Screen.Sync()
			}
		}
	}()

	// 引数の読み込み
	flag.Parse()
	var args = flag.Args()
	var W, H int

Arg:
	for i, arg := range args {
		switch arg {
		case "--help", "-help", "help":
			fmt.Printf("Usage:\n" +
				"pic PICTURENAME(S)\n" +
				"Esc to escape view mode\n" +
				"Enter to view next picture\n" +
				"Options:\n" +
				"-h	Set height of the displayed picture[%%]\n" +
				"-w	Set width of the displayed picture[%%]\n",
			)
			os.Exit(0)
		case "-w":
			if len(args) > i {
				W, _ = strconv.Atoi(args[i+1])
			}
			break Arg
		case "-h":
			if len(args) > i {
				H, _ = strconv.Atoi(args[i+1])
			}
			break Arg
		}
	}

mainloop:
	for i := 0; i < len(args); i++ {
		var arg = args[i]
	set:
		f, err := os.Open(arg)
		if err != nil {
			continue
		}
		Screen.Clear()

		var imgr ImageReader
		imgr.New(f)

		imgr.SetTitle(arg)

		// 描画
		err = imgr.Set(W, H)
		if err != nil {
			continue
		}

		Screen.Show()
	subloop:
		for {
			select {
			case <-quit:
				break mainloop
			case <-next:
				PutRow = 0
				break subloop
			case <-resize:
				goto set
			case <-back:
				if i > 0 {
					i -= 2
					break subloop
				}
			case <-time.After(time.Millisecond * 50):
			}
		}

	}

	Screen.Fini()
}
