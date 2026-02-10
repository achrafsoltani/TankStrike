package main

import (
	"log"
	"time"

	"github.com/AchrafSoltani/TankStrike/config"
	"github.com/AchrafSoltani/TankStrike/game"
	"github.com/AchrafSoltani/glow"
)

func main() {
	win, err := glow.NewWindow("TankStrike", config.WindowWidth, config.WindowHeight)
	if err != nil {
		log.Fatal(err)
	}
	defer win.Close()

	g := game.NewGame()
	canvas := win.Canvas()
	running := true
	lastTime := time.Now()

	for running {
		now := time.Now()
		dt := now.Sub(lastTime).Seconds()
		lastTime = now

		if dt > 0.05 {
			dt = 0.05
		}

		for {
			event := win.PollEvent()
			if event == nil {
				break
			}
			switch event.Type {
			case glow.EventQuit:
				running = false
			case glow.EventKeyDown:
				if event.Key == glow.KeyF11 {
					win.SetFullscreen(!win.IsFullscreen())
				}
				g.KeyDown(event.Key)
			case glow.EventKeyUp:
				g.KeyUp(event.Key)
			case glow.EventWindowResize:
				g.OnResize(event.Width, event.Height)
			}
		}

		g.Update(dt)

		canvas.Clear(glow.Black)
		g.Draw(canvas)
		win.Present()

		elapsed := time.Since(now)
		target := time.Second / 60
		if elapsed < target {
			time.Sleep(target - elapsed)
		}
	}
}
