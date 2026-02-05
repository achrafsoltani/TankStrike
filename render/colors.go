package render

import "github.com/AchrafSoltani/glow"

// Game colour palette
var (
	// Background
	ColorBackground = glow.RGB(0, 0, 0)
	ColorPlayArea   = glow.RGB(0, 0, 0)
	ColorHUDBG      = glow.RGB(99, 99, 99)

	// Brick tile
	ColorBrick      = glow.RGB(165, 81, 33)
	ColorBrickDark  = glow.RGB(107, 53, 20)
	ColorBrickLight = glow.RGB(198, 108, 58)
	ColorMortar     = glow.RGB(74, 37, 16)

	// Steel tile
	ColorSteel      = glow.RGB(180, 180, 180)
	ColorSteelLight = glow.RGB(220, 220, 220)
	ColorSteelDark  = glow.RGB(120, 120, 120)
	ColorSteelRivet = glow.RGB(80, 80, 80)

	// Water tile
	ColorWater     = glow.RGB(0, 51, 153)
	ColorWaterWave = glow.RGB(51, 102, 204)

	// Ice tile
	ColorIce      = glow.RGB(180, 220, 240)
	ColorIceGlint = glow.RGB(230, 245, 255)

	// Forest tile
	ColorForest1 = glow.RGB(0, 100, 0)
	ColorForest2 = glow.RGB(0, 130, 0)
	ColorForest3 = glow.RGB(34, 139, 34)

	// Eagle
	ColorEagleBody = glow.RGB(220, 180, 40)
	ColorEagleWing = glow.RGB(180, 140, 30)
	ColorEagleDead = glow.RGB(100, 100, 100)

	// Player tank
	ColorPlayerBody  = glow.RGB(0, 160, 0)
	ColorPlayerTread = glow.RGB(0, 100, 0)
	ColorPlayerDark  = glow.RGB(0, 80, 0)

	// Enemy tanks
	ColorEnemyBasicBody  = glow.RGB(190, 190, 190)
	ColorEnemyBasicTread = glow.RGB(130, 130, 130)

	ColorEnemyFastBody  = glow.RGB(255, 220, 0)
	ColorEnemyFastTread = glow.RGB(200, 170, 0)

	ColorEnemyPowerBody  = glow.RGB(255, 60, 100)
	ColorEnemyPowerTread = glow.RGB(200, 40, 70)

	ColorEnemyArmourBody  = glow.RGB(0, 180, 100)
	ColorEnemyArmourTread = glow.RGB(0, 130, 70)

	// Bullets
	ColorBulletPlayer = glow.RGB(255, 255, 200)
	ColorBulletEnemy  = glow.RGB(255, 150, 150)
	ColorBulletTrail  = glow.RGB(150, 150, 100)

	// UI colours
	ColorWhite     = glow.RGB(255, 255, 255)
	ColorBlack     = glow.RGB(0, 0, 0)
	ColorRed       = glow.RGB(255, 0, 0)
	ColorYellow    = glow.RGB(255, 255, 0)
	ColorGray      = glow.RGB(128, 128, 128)
	ColorDarkGray  = glow.RGB(64, 64, 64)
	ColorOrange    = glow.RGB(255, 165, 0)
	ColorCyan      = glow.RGB(0, 255, 255)

	// Particles
	ColorExplosion1 = glow.RGB(255, 200, 50)
	ColorExplosion2 = glow.RGB(255, 140, 20)
	ColorExplosion3 = glow.RGB(255, 80, 0)
	ColorExplosion4 = glow.RGB(200, 40, 0)
	ColorSpark      = glow.RGB(255, 255, 200)
	ColorDebris1    = glow.RGB(165, 81, 33)
	ColorDebris2    = glow.RGB(120, 60, 25)

	// Power-ups
	ColorPowerUpStar   = glow.RGB(255, 255, 0)
	ColorPowerUpTank   = glow.RGB(0, 200, 0)
	ColorPowerUpHelmet = glow.RGB(200, 200, 200)
	ColorPowerUpShovel = glow.RGB(139, 119, 101)
	ColorPowerUpBomb   = glow.RGB(255, 50, 50)
	ColorPowerUpClock  = glow.RGB(100, 150, 255)

	// HUD
	ColorHUDEnemyIcon = glow.RGB(200, 60, 60)
	ColorHUDText      = glow.RGB(255, 255, 255)
	ColorHUDLevelBG   = glow.RGB(60, 60, 60)
)
