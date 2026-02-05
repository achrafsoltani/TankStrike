package render

import (
	"github.com/AchrafSoltani/TankStrike/config"
	"github.com/AchrafSoltani/TankStrike/entity"
	"github.com/AchrafSoltani/glow"
)

// TankColors holds the colour scheme for a tank type.
type TankColors struct {
	Body  glow.Color
	Tread glow.Color
	Dark  glow.Color
}

var (
	PlayerColors = TankColors{ColorPlayerBody, ColorPlayerTread, ColorPlayerDark}

	EnemyBasicColors  = TankColors{ColorEnemyBasicBody, ColorEnemyBasicTread, glow.RGB(120, 120, 120)}
	EnemyFastColors   = TankColors{ColorEnemyFastBody, ColorEnemyFastTread, glow.RGB(180, 150, 0)}
	EnemyPowerColors  = TankColors{ColorEnemyPowerBody, ColorEnemyPowerTread, glow.RGB(180, 30, 60)}
	EnemyArmourColors = TankColors{ColorEnemyArmourBody, ColorEnemyArmourTread, glow.RGB(0, 120, 60)}
)

// DrawTank draws a tank at its position with the given colour scheme.
func DrawTank(canvas *glow.Canvas, t *entity.Tank, colors TankColors, offsetX, offsetY int) {
	if !t.Alive {
		return
	}

	px := int(t.X) + offsetX
	py := int(t.Y) + offsetY
	size := config.TankSize

	// Centre of tank
	cx := px + size/2
	cy := py + size/2
	bodyHalf := config.TankBodySize / 2

	switch t.Dir {
	case entity.DirUp, entity.DirDown:
		drawTankVertical(canvas, cx, cy, bodyHalf, t, colors)
	case entity.DirLeft, entity.DirRight:
		drawTankHorizontal(canvas, cx, cy, bodyHalf, t, colors)
	}
}

func drawTankVertical(canvas *glow.Canvas, cx, cy, bodyHalf int, t *entity.Tank, colors TankColors) {
	bw := config.TreadWidth
	bl := config.TreadLength

	// Left tread
	lx := cx - bodyHalf - bw
	ly := cy - bl/2
	drawTread(canvas, lx, ly, bw, bl, true, t.TreadFrame, colors)

	// Right tread
	rx := cx + bodyHalf
	ry := cy - bl/2
	drawTread(canvas, rx, ry, bw, bl, true, t.TreadFrame, colors)

	// Body
	canvas.DrawRect(cx-bodyHalf, cy-bodyHalf, config.TankBodySize, config.TankBodySize, colors.Body)

	// Inner detail
	canvas.DrawRect(cx-bodyHalf+4, cy-bodyHalf+4, config.TankBodySize-8, config.TankBodySize-8, colors.Dark)

	// Turret base (circle-ish)
	canvas.FillCircle(cx, cy, 8, colors.Body)

	// Barrel
	bx := cx - config.BarrelWidth/2
	if t.Dir == entity.DirUp {
		canvas.DrawRect(bx, cy-bodyHalf-config.BarrelLength+4, config.BarrelWidth, config.BarrelLength, colors.Tread)
	} else {
		canvas.DrawRect(bx, cy+bodyHalf-4, config.BarrelWidth, config.BarrelLength, colors.Tread)
	}
}

func drawTankHorizontal(canvas *glow.Canvas, cx, cy, bodyHalf int, t *entity.Tank, colors TankColors) {
	bw := config.TreadWidth
	bl := config.TreadLength

	// Top tread
	tx := cx - bl/2
	ty := cy - bodyHalf - bw
	drawTread(canvas, tx, ty, bl, bw, false, t.TreadFrame, colors)

	// Bottom tread
	bx := cx - bl/2
	by := cy + bodyHalf
	drawTread(canvas, bx, by, bl, bw, false, t.TreadFrame, colors)

	// Body
	canvas.DrawRect(cx-bodyHalf, cy-bodyHalf, config.TankBodySize, config.TankBodySize, colors.Body)

	// Inner detail
	canvas.DrawRect(cx-bodyHalf+4, cy-bodyHalf+4, config.TankBodySize-8, config.TankBodySize-8, colors.Dark)

	// Turret base
	canvas.FillCircle(cx, cy, 8, colors.Body)

	// Barrel
	barrelY := cy - config.BarrelWidth/2
	if t.Dir == entity.DirLeft {
		canvas.DrawRect(cx-bodyHalf-config.BarrelLength+4, barrelY, config.BarrelLength, config.BarrelWidth, colors.Tread)
	} else {
		canvas.DrawRect(cx+bodyHalf-4, barrelY, config.BarrelLength, config.BarrelWidth, colors.Tread)
	}
}

func drawTread(canvas *glow.Canvas, x, y, w, h int, vertical bool, frame int, colors TankColors) {
	canvas.DrawRect(x, y, w, h, colors.Tread)

	// Alternating dark/light strips for animation
	if vertical {
		stripH := 4
		for i := 0; i < h; i += stripH * 2 {
			offset := (frame * 2) % (stripH * 2)
			sy := y + i + offset
			if sy+stripH <= y+h {
				canvas.DrawRect(x, sy, w, stripH, colors.Dark)
			}
		}
	} else {
		stripW := 4
		for i := 0; i < w; i += stripW * 2 {
			offset := (frame * 2) % (stripW * 2)
			sx := x + i + offset
			if sx+stripW <= x+w {
				canvas.DrawRect(sx, y, stripW, h, colors.Dark)
			}
		}
	}
}

// DrawShield draws the invulnerability shield around a tank.
func DrawShield(canvas *glow.Canvas, t *entity.Tank, offsetX, offsetY int, time float64) {
	cx := int(t.X) + offsetX + config.TankSize/2
	cy := int(t.Y) + offsetY + config.TankSize/2

	// Pulsing shield
	alpha := int(time*10) % 2
	if alpha == 0 {
		canvas.DrawCircle(cx, cy, config.TankSize/2+2, ColorWhite)
		canvas.DrawCircle(cx, cy, config.TankSize/2+3, ColorCyan)
	}
}
