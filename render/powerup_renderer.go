package render

import (
	"github.com/AchrafSoltani/TankStrike/entity"
	"github.com/AchrafSoltani/glow"
)

// DrawPowerUp renders a power-up.
func DrawPowerUp(canvas *glow.Canvas, p *entity.PowerUp, offsetX, offsetY int) {
	if !p.Active || !p.IsVisible() {
		return
	}

	px := int(p.X) + offsetX
	py := int(p.Y) + offsetY
	size := 24

	// Background box
	canvas.DrawRect(px, py, size, size, ColorBlack)
	canvas.DrawRectOutline(px, py, size, size, ColorWhite)

	// Icon depends on type
	switch p.Type {
	case entity.PowerUpStar:
		drawStarIcon(canvas, px, py, size)
	case entity.PowerUpTank:
		drawTankIcon(canvas, px, py, size)
	case entity.PowerUpHelmet:
		drawHelmetIcon(canvas, px, py, size)
	case entity.PowerUpShovel:
		drawShovelIcon(canvas, px, py, size)
	case entity.PowerUpBomb:
		drawBombIcon(canvas, px, py, size)
	case entity.PowerUpClock:
		drawClockIcon(canvas, px, py, size)
	}
}

func drawStarIcon(canvas *glow.Canvas, px, py, size int) {
	cx := px + size/2
	cy := py + size/2
	// Simple star shape
	canvas.DrawRect(cx-2, cy-6, 4, 12, ColorPowerUpStar)
	canvas.DrawRect(cx-6, cy-2, 12, 4, ColorPowerUpStar)
	canvas.DrawRect(cx-4, cy-4, 8, 8, ColorPowerUpStar)
}

func drawTankIcon(canvas *glow.Canvas, px, py, _ int) {
	canvas.DrawRect(px+4, py+6, 16, 14, ColorPowerUpTank)
	canvas.DrawRect(px+9, py+2, 6, 6, ColorPowerUpTank)
}

func drawHelmetIcon(canvas *glow.Canvas, px, py, size int) {
	cx := px + size/2
	canvas.FillCircle(cx, py+10, 8, ColorPowerUpHelmet)
	canvas.DrawRect(px+4, py+12, 16, 6, ColorPowerUpHelmet)
}

func drawShovelIcon(canvas *glow.Canvas, px, py, _ int) {
	canvas.DrawRect(px+10, py+2, 4, 14, ColorPowerUpShovel)
	canvas.DrawRect(px+6, py+16, 12, 4, ColorPowerUpShovel)
}

func drawBombIcon(canvas *glow.Canvas, px, py, size int) {
	cx := px + size/2
	canvas.FillCircle(cx, py+14, 7, ColorPowerUpBomb)
	canvas.DrawRect(cx-1, py+3, 2, 6, ColorPowerUpBomb)
}

func drawClockIcon(canvas *glow.Canvas, px, py, size int) {
	cx := px + size/2
	cy := py + size/2
	canvas.DrawCircle(cx, cy, 9, ColorPowerUpClock)
	canvas.DrawLine(cx, cy, cx, cy-6, ColorPowerUpClock) // hour hand
	canvas.DrawLine(cx, cy, cx+4, cy, ColorPowerUpClock) // minute hand
}
