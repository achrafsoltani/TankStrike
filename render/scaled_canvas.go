package render

import (
	"github.com/AchrafSoltani/TankStrike/config"
	"github.com/AchrafSoltani/glow"
)

// ScaledCanvas wraps a glow.Canvas and transparently scales all draw
// operations according to a config.Layout. Game code continues to use
// the original 864x672 coordinate system.
type ScaledCanvas struct {
	canvas *glow.Canvas
	layout config.Layout
}

// NewScaledCanvas creates a ScaledCanvas from a raw canvas and layout.
func NewScaledCanvas(canvas *glow.Canvas, layout config.Layout) *ScaledCanvas {
	return &ScaledCanvas{canvas: canvas, layout: layout}
}

func (sc *ScaledCanvas) tx(v int) int {
	return sc.layout.OffsetX + int(float64(v)*sc.layout.Scale)
}

func (sc *ScaledCanvas) ty(v int) int {
	return sc.layout.OffsetY + int(float64(v)*sc.layout.Scale)
}

func (sc *ScaledCanvas) ts(v int) int {
	return int(float64(v) * sc.layout.Scale)
}

// Clear fills the entire underlying canvas.
func (sc *ScaledCanvas) Clear(color glow.Color) {
	sc.canvas.Clear(color)
}

// SetPixel draws a scaled block at the given logical pixel.
func (sc *ScaledCanvas) SetPixel(x, y int, color glow.Color) {
	s := sc.ts(1)
	if s <= 1 {
		sc.canvas.SetPixel(sc.tx(x), sc.ty(y), color)
	} else {
		sc.canvas.DrawRect(sc.tx(x), sc.ty(y), s, s, color)
	}
}

// DrawRect draws a filled rectangle in scaled coordinates.
func (sc *ScaledCanvas) DrawRect(x, y, width, height int, color glow.Color) {
	sc.canvas.DrawRect(sc.tx(x), sc.ty(y), sc.ts(width), sc.ts(height), color)
}

// DrawRectOutline draws a rectangle outline in scaled coordinates.
func (sc *ScaledCanvas) DrawRectOutline(x, y, width, height int, color glow.Color) {
	sc.canvas.DrawRectOutline(sc.tx(x), sc.ty(y), sc.ts(width), sc.ts(height), color)
}

// FillCircle draws a filled circle in scaled coordinates.
func (sc *ScaledCanvas) FillCircle(x, y, radius int, color glow.Color) {
	sc.canvas.FillCircle(sc.tx(x), sc.ty(y), sc.ts(radius), color)
}

// DrawCircle draws a circle outline in scaled coordinates.
func (sc *ScaledCanvas) DrawCircle(x, y, radius int, color glow.Color) {
	sc.canvas.DrawCircle(sc.tx(x), sc.ty(y), sc.ts(radius), color)
}

// DrawLine draws a line in scaled coordinates.
func (sc *ScaledCanvas) DrawLine(x0, y0, x1, y1 int, color glow.Color) {
	sc.canvas.DrawLine(sc.tx(x0), sc.ty(y0), sc.tx(x1), sc.ty(y1), color)
}
