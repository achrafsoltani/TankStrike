package system

import (
	"math"

	"github.com/AchrafSoltani/TankStrike/config"
	"github.com/AchrafSoltani/TankStrike/entity"
	"github.com/AchrafSoltani/TankStrike/world"
)

// MoveTank attempts to move a tank in its facing direction.
// Returns true if movement occurred.
func MoveTank(t *entity.Tank, grid *world.Grid, dt float64, otherTanks []BBox) bool {
	if !t.Alive || !t.Moving {
		return false
	}

	dx := t.Dir.DX() * t.Speed * dt
	dy := t.Dir.DY() * t.Speed * dt
	newX := t.X + dx
	newY := t.Y + dy

	// Clamp to play area boundaries
	newX = math.Max(0, math.Min(newX, float64(config.PlayAreaWidth-config.TankSize)))
	newY = math.Max(0, math.Min(newY, float64(config.PlayAreaHeight-config.TankSize)))

	// Snap to sub-block grid on the axis perpendicular to movement
	// This makes tanks align to grid lanes when turning
	switch t.Dir {
	case entity.DirUp, entity.DirDown:
		snapped := math.Round(t.X/float64(config.SubBlock)) * float64(config.SubBlock)
		newX = snapped
	case entity.DirLeft, entity.DirRight:
		snapped := math.Round(t.Y/float64(config.SubBlock)) * float64(config.SubBlock)
		newY = snapped
	}

	// Check grid collision
	if !checkGridPassable(newX, newY, grid) {
		return false
	}

	// Check tank-tank collision
	myBox := BBox{X: newX, Y: newY, W: config.TankSize, H: config.TankSize}
	for _, other := range otherTanks {
		if boxOverlap(myBox, other) {
			return false
		}
	}

	t.X = newX
	t.Y = newY
	return true
}

// MovePlayerTank handles player movement including ice sliding.
func MovePlayerTank(p *entity.PlayerTank, grid *world.Grid, dt float64, otherTanks []BBox) {
	if !p.Alive || p.Respawning {
		return
	}

	// Check if on ice
	cx := int(p.X+24) / config.SubBlock
	cy := int(p.Y+24) / config.SubBlock
	onIce := false
	for dy := 0; dy < 2; dy++ {
		for dx := 0; dx < 2; dx++ {
			if grid.Get(cx+dx, cy+dy) == world.TileIce {
				onIce = true
			}
		}
	}
	p.OnIce = onIce

	if p.Moving {
		if onIce {
			p.Tank.Speed = config.PlayerSpeed * config.IceSlideMultiplier
		} else {
			p.Tank.Speed = config.PlayerSpeed
		}
		MoveTank(&p.Tank, grid, dt, otherTanks)
	}
}

// checkGridPassable checks if a tank-sized rectangle at pixel position (x,y) fits in passable tiles.
func checkGridPassable(x, y float64, grid *world.Grid) bool {
	sb := float64(config.SubBlock)

	// Check all sub-blocks the tank overlaps
	x0 := int(math.Floor(x / sb))
	y0 := int(math.Floor(y / sb))
	x1 := int(math.Floor((x + float64(config.TankSize) - 1) / sb))
	y1 := int(math.Floor((y + float64(config.TankSize) - 1) / sb))

	for sy := y0; sy <= y1; sy++ {
		for sx := x0; sx <= x1; sx++ {
			if !grid.Get(sx, sy).IsPassable() {
				return false
			}
		}
	}
	return true
}

// BBox is an axis-aligned bounding box.
type BBox struct {
	X, Y float64
	W, H float64
}

func boxOverlap(a, b BBox) bool {
	return a.X < b.X+b.W && a.X+a.W > b.X && a.Y < b.Y+b.H && a.Y+a.H > b.Y
}

// TankBBox returns the bounding box for a tank.
func TankBBox(t *entity.Tank) BBox {
	return BBox{X: t.X, Y: t.Y, W: config.TankSize, H: config.TankSize}
}
