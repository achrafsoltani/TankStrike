package render

import (
	"fmt"

	"github.com/AchrafSoltani/TankStrike/config"
	"github.com/AchrafSoltani/glow"
)

// HUDRenderer draws the sidebar HUD.
type HUDRenderer struct {
	X int // left edge of HUD area
}

// NewHUDRenderer creates a new HUD renderer.
func NewHUDRenderer() *HUDRenderer {
	return &HUDRenderer{
		X: config.Padding + config.PlayAreaWidth + config.Padding,
	}
}

// DrawHUD draws the complete HUD sidebar.
func (h *HUDRenderer) DrawHUD(canvas *glow.Canvas, enemiesRemaining int, lives int, level int, score int) {
	// Background
	canvas.DrawRect(h.X, 0, config.HUDWidth, config.WindowHeight, ColorHUDBG)

	x := h.X + 16
	y := 24

	// Enemy count icons (small red squares in a 2-column grid)
	DrawText(canvas, "ENEMY", x, y, ColorHUDText, 1)
	y += 16
	for i := 0; i < enemiesRemaining; i++ {
		col := i % 2
		row := i / 2
		ix := x + col*20
		iy := y + row*14
		drawEnemyIcon(canvas, ix, iy)
	}
	maxRows := (config.EnemiesPerLevel + 1) / 2
	y += maxRows*14 + 16

	// Separator
	canvas.DrawRect(h.X+8, y, config.HUDWidth-16, 2, ColorDarkGray)
	y += 12

	// Player info
	DrawText(canvas, "1P", x, y, ColorYellow, 1)
	y += 14

	// Lives
	drawPlayerIcon(canvas, x, y)
	DrawText(canvas, fmt.Sprintf("x%d", lives), x+20, y+2, ColorHUDText, 1)
	y += 20

	// Stars
	DrawText(canvas, "SCORE", x, y, ColorHUDText, 1)
	y += 12
	DrawText(canvas, fmt.Sprintf("%06d", score), x, y, ColorYellow, 1)
	y += 24

	// Separator
	canvas.DrawRect(h.X+8, y, config.HUDWidth-16, 2, ColorDarkGray)
	y += 12

	// Level
	canvas.DrawRect(x, y, config.HUDWidth-40, 28, ColorHUDLevelBG)
	DrawText(canvas, "STAGE", x+8, y+2, ColorHUDText, 1)
	DrawText(canvas, fmt.Sprintf("  %2d", level+1), x+8, y+14, ColorYellow, 1)
}

func drawEnemyIcon(canvas *glow.Canvas, x, y int) {
	canvas.DrawRect(x, y, 12, 10, ColorHUDEnemyIcon)
	canvas.DrawRect(x+4, y-2, 4, 3, ColorHUDEnemyIcon) // barrel
}

func drawPlayerIcon(canvas *glow.Canvas, x, y int) {
	canvas.DrawRect(x, y, 12, 14, ColorPlayerBody)
	canvas.DrawRect(x+4, y-2, 4, 4, ColorPlayerTread) // barrel
}
