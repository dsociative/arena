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
	spawnCount = flag.Int("spawn", 1, "spawn iteration count")
	pprof      = flag.Bool("pprof", false, "enables pprof http handler")
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
				case tcell.KeyCtrlP:
					gmap.Debug.Toggle()
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
			http.ListenAndServe("localhost:8787", nil)
		}()
	}

	arena := gmap.NewArena(screen, 150, 50)
	game := game.NewGame(arena, *spawnCount)
	game.Run(ctx, keyChan)
}
