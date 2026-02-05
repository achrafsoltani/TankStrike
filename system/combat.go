package system

import (
	"math"

	"github.com/AchrafSoltani/TankStrike/config"
	"github.com/AchrafSoltani/TankStrike/entity"
	"github.com/AchrafSoltani/TankStrike/render"
	"github.com/AchrafSoltani/TankStrike/world"
)

// BulletGridCollision checks bullet-to-grid collisions and destroys tiles.
// Returns true if the bullet was consumed.
func BulletGridCollision(b *entity.Bullet, grid *world.Grid, particles *render.ParticlePool) bool {
	if !b.Active {
		return false
	}

	// Find the sub-block the bullet centre is in
	bx := int(math.Floor(b.X+float64(config.BulletSize)/2)) / config.SubBlock
	by := int(math.Floor(b.Y+float64(config.BulletSize)/2)) / config.SubBlock

	tile := grid.Get(bx, by)
	if !tile.BlocksBullets() {
		return false
	}

	switch tile {
	case world.TileBrick:
		grid.Destroy(bx, by, b.Power)
		cx := float64(bx*config.SubBlock) + float64(config.SubBlock)/2
		cy := float64(by*config.SubBlock) + float64(config.SubBlock)/2
		particles.SpawnDebris(cx, cy)
		b.Active = false
		return true

	case world.TileSteel:
		if b.Power >= 3 {
			grid.Destroy(bx, by, b.Power)
			cx := float64(bx*config.SubBlock) + float64(config.SubBlock)/2
			cy := float64(by*config.SubBlock) + float64(config.SubBlock)/2
			particles.SpawnDebris(cx, cy)
		} else {
			cx := float64(bx*config.SubBlock) + float64(config.SubBlock)/2
			cy := float64(by*config.SubBlock) + float64(config.SubBlock)/2
			particles.SpawnSpark(cx, cy)
		}
		b.Active = false
		return true

	case world.TileEagle:
		grid.Destroy(bx, by, b.Power)
		cx := float64(bx*config.SubBlock) + float64(config.SubBlock)/2
		cy := float64(by*config.SubBlock) + float64(config.SubBlock)/2
		particles.SpawnExplosion(cx, cy, 40)
		b.Active = false
		return true
	}

	return false
}

// BulletTankCollision checks if a bullet hits a tank. Returns true if hit.
func BulletTankCollision(b *entity.Bullet, t *entity.Tank) bool {
	if !b.Active || !t.Alive {
		return false
	}

	// AABB overlap
	bBox := BBox{X: b.X, Y: b.Y, W: float64(config.BulletSize), H: float64(config.BulletSize)}
	tBox := TankBBox(t)
	return boxOverlap(bBox, tBox)
}

// CountPlayerBullets counts active player bullets.
func CountPlayerBullets(bullets []*entity.Bullet) int {
	count := 0
	for _, b := range bullets {
		if b.Active && b.IsPlayer {
			count++
		}
	}
	return count
}
