package entity

import (
	"github.com/AchrafSoltani/TankStrike/config"
	"github.com/AchrafSoltani/glow"
)

// PlayerTank extends Tank with player-specific features.
type PlayerTank struct {
	Tank
	Lives         int
	Score         int
	Stars         int // upgrade level (0-3)
	ShieldTimer   float64
	RespawnTimer  float64
	Respawning    bool

	// Ice sliding
	SlideVX float64
	SlideVY float64
	OnIce   bool
}

// NewPlayerTank creates a new player tank at the default spawn position.
func NewPlayerTank() *PlayerTank {
	// Player spawns at bottom centre-left (sub-block 8,24 → pixel 192, 576)
	spawnX := float64(8 * config.SubBlock)
	spawnY := float64(24 * config.SubBlock)
	p := &PlayerTank{
		Tank:  NewTank(spawnX, spawnY, config.PlayerSpeed, 1),
		Lives: config.StartLives,
	}
	p.Tank.BulletSpeed = config.PlayerBulletSpd
	p.Tank.CooldownRate = 0.3
	return p
}

// HandleInput reads key state and updates movement direction.
func (p *PlayerTank) HandleInput(keys map[glow.Key]bool) {
	if !p.Alive || p.Respawning {
		p.Moving = false
		return
	}

	p.Moving = false
	if keys[glow.KeyW] || keys[glow.KeyUp] {
		p.Dir = DirUp
		p.Moving = true
	} else if keys[glow.KeyS] || keys[glow.KeyDown] {
		p.Dir = DirDown
		p.Moving = true
	} else if keys[glow.KeyA] || keys[glow.KeyLeft] {
		p.Dir = DirLeft
		p.Moving = true
	} else if keys[glow.KeyD] || keys[glow.KeyRight] {
		p.Dir = DirRight
		p.Moving = true
	}
}

// WantsToShoot returns true if the player is pressing fire.
func (p *PlayerTank) WantsToShoot(keys map[glow.Key]bool) bool {
	return keys[glow.KeySpace]
}

// Respawn resets the player tank to the spawn point.
func (p *PlayerTank) Respawn() {
	p.X = float64(8 * config.SubBlock)
	p.Y = float64(24 * config.SubBlock)
	p.Dir = DirUp
	p.HP = 1
	p.Alive = true
	p.Respawning = false
	p.RespawnTimer = 0
	p.ShieldTimer = 3.0 // brief invulnerability on respawn
	p.Moving = false
	p.ShootCooldown = 0
	p.SlideVX = 0
	p.SlideVY = 0
}

// UpdatePlayer handles player-specific update logic.
func (p *PlayerTank) UpdatePlayer(dt float64) {
	p.Tank.Update(dt)

	if p.ShieldTimer > 0 {
		p.ShieldTimer -= dt
	}

	if p.Respawning {
		p.RespawnTimer -= dt
		if p.RespawnTimer <= 0 {
			p.Respawn()
		}
	}
}

// Die handles player death.
func (p *PlayerTank) Die() {
	p.Alive = false
	p.Lives--
	if p.Lives > 0 {
		p.Respawning = true
		p.RespawnTimer = config.RespawnDelay
	}
}

// ApplyStar gives the player a star upgrade.
func (p *PlayerTank) ApplyStar() {
	if p.Stars < 3 {
		p.Stars++
	}
	p.PowerLevel = p.Stars
	switch p.Stars {
	case 1:
		p.CooldownRate = 0.2
		p.BulletSpeed = config.PlayerBulletSpd * 1.2
	case 2:
		p.CooldownRate = 0.15
		p.BulletSpeed = config.PlayerBulletSpd * 1.4
	case 3:
		p.CooldownRate = 0.1
		p.BulletSpeed = config.PlayerBulletSpd * 1.6
	}
}

// IsInvulnerable returns whether the player currently has a shield.
func (p *PlayerTank) IsInvulnerable() bool {
	return p.ShieldTimer > 0
}
