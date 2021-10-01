package pic

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

func Do() {
	defer Screen.Fini()

	// チャネルの生成
	var KeyEvent = make(chan *tcell.EventKey)

	go func() {
		for {
			ev := Screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					// 終了
					Screen.Fini()
					os.Exit(0)
				default:
					switch ev.Rune() {
					case 'q':
						// 終了
						Screen.Fini()
						os.Exit(0)
					default:
						KeyEvent <- ev
					}
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
				"Ctrl + I initialize settings\n" +
				"Ctrl + L Reload picture\n" +
				"Options:\n" +
				"-h	Set height of the displayed picture[%%]\n" +
				"-w	Set width of the displayed picture[%%]\n",
			)
			Screen.Fini()
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

	for i := 0; i < len(args); i++ {
		var arg = args[i]

		var delta = Pos{}

		var ZoomRate float64 = 1
		var imgr ImageReader

		var init = true
		var zoom = true

		f, err := os.Open(arg)
		if err != nil {
			continue
		}

		err = imgr.New(f)
		if err != nil {
			f.Close()
			continue
		}
		f.Close()

	set:
		Screen.Clear()

		if init {
			// get suit rate
			imgr.getSuitRate(W, H)
			init = false
		}

		imgr.Title = arg

		if zoom {
			imgr.Zoom(ZoomRate)
			zoom = false
		}

		// 描画
		imgr.SetTitle(ZoomRate)
		imgr.Set(delta)

		Screen.Show()

	subloop:
		for {

			select {
			case Key := <-KeyEvent:
				switch Key.Key() {

				case tcell.KeyEnter:
					// 次の写真へ
					PutRow = 0
					break subloop
				case tcell.KeyCtrlL:
					// 再読み込み
					Screen.Sync()
					goto set
				case tcell.KeyCtrlI:
					// 初期化
					init = true
					zoom = true
					ZoomRate = 1
					delta = Pos{}
					goto set
				case tcell.KeyBackspace, tcell.KeyBackspace2:
					// 一つ前の写真へ
					if i > 0 {
						i -= 2
						break subloop
					}
				case tcell.KeyUp:
					delta.Y += 3
					goto set
				case tcell.KeyDown:
					delta.Y -= 3
					goto set
				case tcell.KeyLeft:
					delta.X += 3
					goto set
				case tcell.KeyRight:
					delta.X -= 3
					goto set
				}
				switch Key.Rune() {
				case '+':
					ZoomRate += 0.5
					zoom = true
					goto set
				case '-':
					ZoomRate -= 0.5
					zoom = true
					goto set
				}

			case <-time.After(time.Millisecond * 50):
			}
		}

	}

	Screen.Fini()
}
