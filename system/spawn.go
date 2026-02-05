package system

import (
	"math/rand"

	"github.com/AchrafSoltani/TankStrike/config"
	"github.com/AchrafSoltani/TankStrike/entity"
)

// Spawn points (sub-block coordinates, top row)
var SpawnPoints = [][2]int{
	{0, 0},   // top-left
	{12, 0},  // top-centre
	{24, 0},  // top-right
}

// Spawner manages enemy spawning.
type Spawner struct {
	Queue         []entity.EnemyType // remaining enemies to spawn
	Timer         float64
	NextSpawnIdx  int // cycles through spawn points
	TotalSpawned  int
	TotalForLevel int
}

// NewSpawner creates a new spawner for a level.
func NewSpawner(level int) *Spawner {
	s := &Spawner{
		Timer:         2.0, // initial delay before first spawn
		TotalForLevel: config.EnemiesPerLevel,
	}
	s.buildQueue(level)
	return s
}

func (s *Spawner) buildQueue(level int) {
	total := config.EnemiesPerLevel
	s.Queue = make([]entity.EnemyType, 0, total)

	// Mix of enemy types depends on level
	for i := 0; i < total; i++ {
		var typ entity.EnemyType
		roll := rand.Float64()
		switch {
		case level < 3:
			if roll < 0.6 {
				typ = entity.EnemyBasic
			} else if roll < 0.85 {
				typ = entity.EnemyFast
			} else {
				typ = entity.EnemyPower
			}
		case level < 6:
			if roll < 0.35 {
				typ = entity.EnemyBasic
			} else if roll < 0.6 {
				typ = entity.EnemyFast
			} else if roll < 0.85 {
				typ = entity.EnemyPower
			} else {
				typ = entity.EnemyArmour
			}
		default:
			if roll < 0.2 {
				typ = entity.EnemyBasic
			} else if roll < 0.45 {
				typ = entity.EnemyFast
			} else if roll < 0.7 {
				typ = entity.EnemyPower
			} else {
				typ = entity.EnemyArmour
			}
		}
		s.Queue = append(s.Queue, typ)
	}
}

// Update checks if it's time to spawn a new enemy.
// Returns a new enemy tank if one should spawn, nil otherwise.
func (s *Spawner) Update(dt float64, activeEnemies int) *entity.EnemyTank {
	if len(s.Queue) == 0 {
		return nil
	}
	if activeEnemies >= config.MaxActiveEnemies {
		return nil
	}

	s.Timer -= dt
	if s.Timer > 0 {
		return nil
	}

	s.Timer = config.SpawnInterval

	// Pick spawn point
	sp := SpawnPoints[s.NextSpawnIdx%len(SpawnPoints)]
	s.NextSpawnIdx++

	typ := s.Queue[0]
	s.Queue = s.Queue[1:]

	// Every 4th enemy carries a power-up
	hasPowerUp := s.TotalSpawned%4 == 3
	s.TotalSpawned++

	x := float64(sp[0] * config.SubBlock)
	y := float64(sp[1] * config.SubBlock)

	return entity.NewEnemyTank(x, y, typ, hasPowerUp)
}

// Remaining returns the number of enemies still to spawn.
func (s *Spawner) Remaining() int {
	return len(s.Queue)
}

// Done returns true when all enemies have been spawned.
func (s *Spawner) Done() bool {
	return len(s.Queue) == 0
}
