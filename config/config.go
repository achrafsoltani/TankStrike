package config

// Grid dimensions
const (
	GridWidth  = 26
	GridHeight = 26
	SubBlock   = 24 // pixels per sub-block
)

// Window dimensions
const (
	PlayAreaWidth  = GridWidth * SubBlock  // 624
	PlayAreaHeight = GridHeight * SubBlock // 624
	Padding        = 24
	HUDWidth       = 192
	WindowWidth    = PlayAreaWidth + Padding*2 + HUDWidth // 864
	WindowHeight   = PlayAreaHeight + Padding*2           // 672
)

// Tank dimensions (2x2 sub-blocks)
const (
	TankSize     = SubBlock * 2 // 48 pixels
	TankBodySize = 36
	TreadWidth   = 6
	TreadLength  = 40
	BarrelWidth  = 6
	BarrelLength = 16
)

// Gameplay timing
const (
	PlayerSpeed      = 120.0 // pixels per second
	PlayerBulletSpd  = 300.0
	MaxPlayerBullets = 2

	EnemySpeedBasic  = 60.0
	EnemySpeedFast   = 120.0
	EnemySpeedPower  = 80.0
	EnemySpeedArmour = 50.0

	EnemyBulletSpd = 200.0

	BulletSize = 4

	SpawnInterval    = 3.0 // seconds between enemy spawns
	MaxActiveEnemies = 4
	EnemiesPerLevel  = 20

	RespawnDelay = 2.0 // seconds before player respawns
	StartLives   = 3

	PowerUpDuration = 15.0 // seconds for timed power-ups (helmet, clock, shovel)

	IceSlideMultiplier = 1.6
	IceFriction        = 0.92

	AIDirectionMinTime = 0.5
	AIDirectionMaxTime = 2.5
)

// Scoring
const (
	ScoreBasic  = 100
	ScoreFast   = 200
	ScorePower  = 300
	ScoreArmour = 400
)

// Particle system
const (
	MaxParticles = 512
)
