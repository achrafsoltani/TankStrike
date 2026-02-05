package entity

import "github.com/AchrafSoltani/TankStrike/config"

// EnemyType represents the type of enemy tank.
type EnemyType int

const (
	EnemyBasic  EnemyType = iota
	EnemyFast
	EnemyPower
	EnemyArmour
)

// EnemyTank extends Tank with enemy-specific AI state.
type EnemyTank struct {
	Tank
	Type           EnemyType
	DirTimer       float64 // time until next direction change
	DirInterval    float64 // how often to change direction
	ShootChance    float64 // probability of shooting per second
	ScoreValue     int
	HasPowerUp     bool // drops a power-up when destroyed
	FlashTimer     float64
	FlashForPowerUp bool
}

// NewEnemyTank creates a new enemy tank of the given type at the given position.
func NewEnemyTank(x, y float64, typ EnemyType, hasPowerUp bool) *EnemyTank {
	e := &EnemyTank{
		Type:            typ,
		HasPowerUp:      hasPowerUp,
		FlashForPowerUp: hasPowerUp,
	}

	switch typ {
	case EnemyBasic:
		e.Tank = NewTank(x, y, config.EnemySpeedBasic, 1)
		e.ShootChance = 0.8
		e.DirInterval = 2.0
		e.ScoreValue = config.ScoreBasic
	case EnemyFast:
		e.Tank = NewTank(x, y, config.EnemySpeedFast, 1)
		e.ShootChance = 1.2
		e.DirInterval = 1.0
		e.ScoreValue = config.ScoreFast
	case EnemyPower:
		e.Tank = NewTank(x, y, config.EnemySpeedPower, 1)
		e.ShootChance = 2.0
		e.DirInterval = 1.5
		e.ScoreValue = config.ScorePower
		e.BulletSpeed = config.EnemyBulletSpd * 1.3
	case EnemyArmour:
		e.Tank = NewTank(x, y, config.EnemySpeedArmour, 4)
		e.ShootChance = 0.6
		e.DirInterval = 2.5
		e.ScoreValue = config.ScoreArmour
	}

	if e.BulletSpeed == 0 {
		e.BulletSpeed = config.EnemyBulletSpd
	}
	e.CooldownRate = 1.0
	e.Dir = DirDown
	e.DirTimer = e.DirInterval
	e.Moving = true

	return e
}

// UpdateEnemy handles enemy-specific update logic.
func (e *EnemyTank) UpdateEnemy(dt float64) {
	e.Tank.Update(dt)

	if e.FlashForPowerUp {
		e.FlashTimer += dt
	}
}

// IsFlashing returns true during the bright phase of the power-up flash.
func (e *EnemyTank) IsFlashing() bool {
	if !e.FlashForPowerUp {
		return false
	}
	return int(e.FlashTimer*6)%2 == 0
}
