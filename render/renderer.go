package render

import (
	"github.com/AchrafSoltani/TankStrike/config"
	"github.com/AchrafSoltani/TankStrike/world"
	"github.com/AchrafSoltani/glow"
)

// Renderer handles all drawing operations.
type Renderer struct {
	OffsetX int // play area X offset (padding)
	OffsetY int // play area Y offset (padding)
	Time    float64
}

// NewRenderer creates a new renderer.
func NewRenderer() *Renderer {
	return &Renderer{
		OffsetX: config.Padding,
		OffsetY: config.Padding,
	}
}

// DrawGrid draws all non-overlay tiles.
func (r *Renderer) DrawGrid(canvas *glow.Canvas, g *world.Grid) {
	for y := 0; y < config.GridHeight; y++ {
		for x := 0; x < config.GridWidth; x++ {
			t := g.Get(x, y)
			if t != world.TileEmpty && t != world.TileForest {
				DrawTile(canvas, t, x, y, r.OffsetX, r.OffsetY, r.Time)
			}
		}
	}
}

// DrawForest draws forest overlay tiles (above tanks).
func (r *Renderer) DrawForest(canvas *glow.Canvas, g *world.Grid) {
	DrawForestOverlay(canvas, g, r.OffsetX, r.OffsetY)
}

// DrawPlayAreaBorder draws the border around the play area.
func (r *Renderer) DrawPlayAreaBorder(canvas *glow.Canvas) {
	hudX := r.OffsetX + config.PlayAreaWidth + config.Padding
	canvas.DrawRect(hudX, 0, config.HUDWidth, config.WindowHeight, ColorHUDBG)
}
