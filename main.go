package main

import (
	"context"
	"flag"
	game "github.com/dsociative/arena/game"
	"github.com/dsociative/arena/gmap"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var (
	spawnCount = flag.Int("spawn", 10, "spawn iteration count")
	pprof      = flag.Bool("pprof", false, "enables pprof http handler")
	pprofAddr  = flag.String("pprofAddr", "localhost:8787", "addr for pprof handler")
	fullscreen = flag.Bool("fullscreen", false, "fullscreen")
	tick       = flag.Duration("tick", time.Millisecond*50, "world tick duration")
	x          = flag.Int("x", 80, "width")
	y          = flag.Int("y", 40, "height")
)

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	encoding.Register()

	if err = screen.Init(); err != nil {
		log.Fatal(err)
	}

	width, height := *x, *y
	if *fullscreen {
		width, height = screen.Size()
	}

	if width < 5 || height < 5 {
		log.Fatal("minimum screen size 5x5")
	}

	if *spawnCount < 1 {
		log.Fatal("minimum spawn count 1")
	}

	screen.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	screen.Clear()
	screen.Show()

	keyChan := make(chan *tcell.EventKey)
	ctx, closeFun := context.WithCancel(context.Background())
	go func() {
		for {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				select {
				case keyChan <- ev:
				default:
				}

				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					closeFun()
					return
				case tcell.KeyCtrlL:
					screen.Sync()
				}
			case *tcell.EventResize:
				screen.Sync()
			}
		}
	}()

	if *pprof {
		go func() {
			http.ListenAndServe(*pprofAddr, nil)
		}()
	}

	arena := gmap.NewArena(screen, width, height)
	game := game.NewGame(arena, *spawnCount, *tick)
	game.Run(ctx, keyChan)
}
