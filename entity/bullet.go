package entity

import "github.com/AchrafSoltani/TankStrike/config"

// Bullet represents a projectile.
type Bullet struct {
	X, Y       float64
	Dir        Direction
	Speed      float64
	Power      int  // 0=normal, 3=can destroy steel
	IsPlayer   bool // true if fired by player
	Active     bool
	TrailX     [3]float64
	TrailY     [3]float64
	TrailCount int
}

// NewBullet creates a new bullet.
func NewBullet(x, y float64, dir Direction, speed float64, power int, isPlayer bool) *Bullet {
	return &Bullet{
		X:        x,
		Y:        y,
		Dir:      dir,
		Speed:    speed,
		Power:    power,
		IsPlayer: isPlayer,
		Active:   true,
	}
}

// Update moves the bullet and records trail positions.
func (b *Bullet) Update(dt float64) {
	if !b.Active {
		return
	}

	// Record trail
	if b.TrailCount < 3 {
		b.TrailX[b.TrailCount] = b.X
		b.TrailY[b.TrailCount] = b.Y
		b.TrailCount++
	} else {
		b.TrailX[0] = b.TrailX[1]
		b.TrailY[0] = b.TrailY[1]
		b.TrailX[1] = b.TrailX[2]
		b.TrailY[1] = b.TrailY[2]
		b.TrailX[2] = b.X
		b.TrailY[2] = b.Y
	}

	b.X += b.Dir.DX() * b.Speed * dt
	b.Y += b.Dir.DY() * b.Speed * dt

	// Deactivate if out of play area
	if b.X < -config.BulletSize || b.X > float64(config.PlayAreaWidth) ||
		b.Y < -config.BulletSize || b.Y > float64(config.PlayAreaHeight) {
		b.Active = false
	}
}
