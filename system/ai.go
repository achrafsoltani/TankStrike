package system

import (
	"math"
	"math/rand"

	"github.com/AchrafSoltani/TankStrike/config"
	"github.com/AchrafSoltani/TankStrike/entity"
	"github.com/AchrafSoltani/TankStrike/world"
)

// UpdateEnemyAI updates the AI for a single enemy tank.
func UpdateEnemyAI(e *entity.EnemyTank, grid *world.Grid, dt float64,
	playerX, playerY float64, eagleX, eagleY float64,
	otherTanks []BBox) {

	if !e.Alive {
		return
	}

	e.UpdateEnemy(dt)

	// Direction timer
	e.DirTimer -= dt
	moved := MoveTank(&e.Tank, grid, dt, otherTanks)

	// If blocked or timer expired, pick new direction
	if !moved || e.DirTimer <= 0 {
		pickNewDirection(e, playerX, playerY, eagleX, eagleY)
		e.DirTimer = config.AIDirectionMinTime +
			rand.Float64()*(config.AIDirectionMaxTime-config.AIDirectionMinTime)
	}
}

// ShouldShoot returns whether the enemy should fire this frame.
func ShouldShoot(e *entity.EnemyTank, dt float64) bool {
	if !e.CanShoot() {
		return false
	}
	return rand.Float64() < e.ShootChance*dt
}

func pickNewDirection(e *entity.EnemyTank, playerX, playerY, eagleX, eagleY float64) {
	roll := rand.Float64()
	if roll < 0.4 {
		// Random direction
		dirs := []entity.Direction{entity.DirUp, entity.DirDown, entity.DirLeft, entity.DirRight}
		e.Dir = dirs[rand.Intn(4)]
	} else if roll < 0.7 {
		// Toward player
		e.Dir = directionToward(e.CenterX(), e.CenterY(), playerX, playerY)
	} else {
		// Toward eagle
		e.Dir = directionToward(e.CenterX(), e.CenterY(), eagleX, eagleY)
	}
	e.Moving = true
}

func directionToward(fromX, fromY, toX, toY float64) entity.Direction {
	dx := toX - fromX
	dy := toY - fromY

	if math.Abs(dx) > math.Abs(dy) {
		if dx > 0 {
			return entity.DirRight
		}
		return entity.DirLeft
	}
	if dy > 0 {
		return entity.DirDown
	}
	return entity.DirUp
}
