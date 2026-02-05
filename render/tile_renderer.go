package render

import (
	"math"

	"github.com/AchrafSoltani/TankStrike/config"
	"github.com/AchrafSoltani/TankStrike/world"
	"github.com/AchrafSoltani/glow"
)

// DrawTile draws a single tile at the given sub-block position.
func DrawTile(canvas *glow.Canvas, t world.TileType, sx, sy int, offsetX, offsetY int, time float64) {
	px := offsetX + sx*config.SubBlock
	py := offsetY + sy*config.SubBlock
	s := config.SubBlock

	switch t {
	case world.TileBrick:
		drawBrick(canvas, px, py, s)
	case world.TileSteel:
		drawSteel(canvas, px, py, s)
	case world.TileWater:
		drawWater(canvas, px, py, s, time)
	case world.TileIce:
		drawIce(canvas, px, py, s)
	case world.TileEagle:
		drawEagle(canvas, px, py, s)
	case world.TileEagleDead:
		drawEagleDead(canvas, px, py, s)
	}
}

// DrawForestOverlay draws forest tiles as an overlay (after tanks).
func DrawForestOverlay(canvas *glow.Canvas, g *world.Grid, offsetX, offsetY int) {
	for y := 0; y < config.GridHeight; y++ {
		for x := 0; x < config.GridWidth; x++ {
			if g.Get(x, y) == world.TileForest {
				drawForest(canvas, offsetX+x*config.SubBlock, offsetY+y*config.SubBlock, config.SubBlock)
			}
		}
	}
}

func drawBrick(canvas *glow.Canvas, px, py, s int) {
	canvas.DrawRect(px, py, s, s, ColorBrick)
	for row := 0; row < 4; row++ {
		my := py + row*6
		canvas.DrawRect(px, my, s, 1, ColorMortar)
	}
	for row := 0; row < 4; row++ {
		offset := 0
		if row%2 == 1 {
			offset = s / 2
		}
		my := py + row*6
		for col := 0; col < 3; col++ {
			mx := px + offset + col*12
			if mx >= px && mx < px+s {
				canvas.DrawRect(mx, my, 1, 6, ColorMortar)
			}
		}
	}
	canvas.DrawRect(px+1, py+1, 3, 2, ColorBrickLight)
	canvas.DrawRect(px+13, py+7, 3, 2, ColorBrickLight)
}

func drawSteel(canvas *glow.Canvas, px, py, s int) {
	canvas.DrawRect(px, py, s, s, ColorSteel)
	canvas.DrawRect(px, py, s, 2, ColorSteelLight)
	canvas.DrawRect(px, py, 2, s, ColorSteelLight)
	canvas.DrawRect(px, py+s-2, s, 2, ColorSteelDark)
	canvas.DrawRect(px+s-2, py, 2, s, ColorSteelDark)
	rivetSize := 3
	canvas.DrawRect(px+2, py+2, rivetSize, rivetSize, ColorSteelRivet)
	canvas.DrawRect(px+s-5, py+2, rivetSize, rivetSize, ColorSteelRivet)
	canvas.DrawRect(px+2, py+s-5, rivetSize, rivetSize, ColorSteelRivet)
	canvas.DrawRect(px+s-5, py+s-5, rivetSize, rivetSize, ColorSteelRivet)
}

func drawWater(canvas *glow.Canvas, px, py, s int, time float64) {
	canvas.DrawRect(px, py, s, s, ColorWater)
	for row := 0; row < s; row += 4 {
		offset := int(math.Sin(float64(row)/4.0+time*3.0) * 3)
		wx := px + offset
		if wx < px {
			wx = px
		}
		w := s - 2
		if wx+w > px+s {
			w = px + s - wx
		}
		if w > 0 {
			canvas.DrawRect(wx, py+row, w, 2, ColorWaterWave)
		}
	}
}

func drawIce(canvas *glow.Canvas, px, py, s int) {
	canvas.DrawRect(px, py, s, s, ColorIce)
	for i := 0; i < s; i += 8 {
		gx := px + i
		gy := py + i
		if gx+4 <= px+s && gy+2 <= py+s {
			canvas.DrawRect(gx, gy, 4, 2, ColorIceGlint)
		}
	}
	for i := 4; i < s; i += 8 {
		gx := px + s - i - 4
		gy := py + i
		if gx >= px && gy+2 <= py+s {
			canvas.DrawRect(gx, gy, 4, 2, ColorIceGlint)
		}
	}
}

func drawForest(canvas *glow.Canvas, px, py, s int) {
	canvas.DrawRect(px, py+2, s, s-4, ColorForest1)
	canvas.DrawRect(px+2, py, s-4, s, ColorForest2)
	canvas.DrawRect(px+4, py+4, s-8, s-8, ColorForest3)
	canvas.DrawRect(px+1, py+1, 6, 6, ColorForest2)
	canvas.DrawRect(px+s-8, py+s-8, 6, 6, ColorForest1)
}

func drawEagle(canvas *glow.Canvas, px, py, s int) {
	canvas.DrawRect(px+4, py+4, s-8, s-8, ColorEagleBody)
	canvas.DrawRect(px+1, py+8, 4, s-16, ColorEagleWing)
	canvas.DrawRect(px+s-5, py+8, 4, s-16, ColorEagleWing)
	canvas.DrawRect(px+8, py+2, s-16, 6, ColorEagleWing)
	canvas.DrawRect(px+s/2-1, py+5, 2, 2, ColorBlack)
}

func drawEagleDead(canvas *glow.Canvas, px, py, s int) {
	canvas.DrawRect(px+2, py+s-8, 6, 6, ColorEagleDead)
	canvas.DrawRect(px+10, py+s-6, 5, 4, ColorEagleDead)
	canvas.DrawRect(px+s-10, py+s-10, 7, 8, ColorEagleDead)
	canvas.DrawRect(px+4, py+s-14, 4, 5, ColorEagleDead)
}
