package game

import (
	"github.com/AchrafSoltani/TankStrike/config"
	"github.com/AchrafSoltani/TankStrike/entity"
	"github.com/AchrafSoltani/TankStrike/render"
	"github.com/AchrafSoltani/TankStrike/system"
	"github.com/AchrafSoltani/TankStrike/world"
	"github.com/AchrafSoltani/glow"
)

// Game is the top-level game orchestrator.
type Game struct {
	State     GameState
	Grid      *world.Grid
	Renderer  *render.Renderer
	Input     *system.Input
	Player    *entity.PlayerTank
	Enemies   []*entity.EnemyTank
	Bullets   []*entity.Bullet
	Particles *render.ParticlePool
	Spawner   *system.Spawner
	Level     int
	Time      float64

	// Eagle position (pixel centre)
	EagleX, EagleY float64
}

// NewGame creates a new game instance.
func NewGame() *Game {
	g := &Game{
		State:     StatePlaying,
		Grid:      world.NewGrid(),
		Renderer:  render.NewRenderer(),
		Input:     system.NewInput(),
		Player:    entity.NewPlayerTank(),
		Particles: render.NewParticlePool(),
		Level:     0,
	}
	g.LoadLevel(0)
	return g
}

// LoadLevel loads a level by index.
func (g *Game) LoadLevel(index int) {
	if index >= 0 && index < len(world.Levels) {
		g.Level = index
		world.LoadLevel(g.Grid, world.Levels[index])
		g.Bullets = g.Bullets[:0]
		g.Enemies = g.Enemies[:0]
		g.Spawner = system.NewSpawner(index)
		g.Player.Respawn()
		g.findEagle()
	}
}

func (g *Game) findEagle() {
	for y := 0; y < config.GridHeight; y++ {
		for x := 0; x < config.GridWidth; x++ {
			if g.Grid.Get(x, y) == world.TileEagle {
				g.EagleX = float64(x*config.SubBlock) + float64(config.SubBlock)/2
				g.EagleY = float64(y*config.SubBlock) + float64(config.SubBlock)/2
				return
			}
		}
	}
}

// KeyDown handles key press events.
func (g *Game) KeyDown(key glow.Key) {
	g.Input.KeyDown(key)
}

// KeyUp handles key release events.
func (g *Game) KeyUp(key glow.Key) {
	g.Input.KeyUp(key)
}

// Update advances game state by dt seconds.
func (g *Game) Update(dt float64) {
	g.Time += dt
	g.Renderer.Time = g.Time
	g.Input.Update()

	switch g.State {
	case StatePlaying:
		g.updatePlaying(dt)
	}
}

func (g *Game) updatePlaying(dt float64) {
	// Player input and movement
	g.Player.HandleInput(g.Input.Keys)
	g.Player.UpdatePlayer(dt)

	// Build list of other tank bounding boxes for collision
	otherTanks := g.enemyBBoxes()
	system.MovePlayerTank(g.Player, g.Grid, dt, otherTanks)

	// Player shooting
	if g.Player.WantsToShoot(g.Input.Keys) && g.Player.CanShoot() {
		if system.CountPlayerBullets(g.Bullets) < config.MaxPlayerBullets {
			bx, by := g.Player.Shoot()
			bullet := entity.NewBullet(bx, by, g.Player.Dir, g.Player.BulletSpeed, g.Player.PowerLevel, true)
			g.Bullets = append(g.Bullets, bullet)
		}
	}

	// Spawn enemies
	if enemy := g.Spawner.Update(dt, len(g.Enemies)); enemy != nil {
		g.Enemies = append(g.Enemies, enemy)
	}

	// Update enemies
	for _, e := range g.Enemies {
		if !e.Alive {
			continue
		}
		// Build collision boxes excluding self
		others := g.tankBBoxesExcluding(&e.Tank)
		system.UpdateEnemyAI(e, g.Grid, dt,
			g.Player.CenterX(), g.Player.CenterY(),
			g.EagleX, g.EagleY, others)

		// Enemy shooting
		if system.ShouldShoot(e, dt) {
			bx, by := e.Shoot()
			bullet := entity.NewBullet(bx, by, e.Dir, e.BulletSpeed, 0, false)
			g.Bullets = append(g.Bullets, bullet)
		}
	}

	// Update bullets
	for _, b := range g.Bullets {
		b.Update(dt)
	}

	// Bullet-grid collisions
	for _, b := range g.Bullets {
		system.BulletGridCollision(b, g.Grid, g.Particles)
	}

	// Bullet-tank collisions: player bullets hit enemies
	for _, b := range g.Bullets {
		if !b.Active || !b.IsPlayer {
			continue
		}
		for _, e := range g.Enemies {
			if !e.Alive {
				continue
			}
			if system.BulletTankCollision(b, &e.Tank) {
				b.Active = false
				destroyed := e.Hit(1)
				if destroyed {
					g.Player.Score += e.ScoreValue
					g.Particles.SpawnExplosion(e.CenterX(), e.CenterY(), 35)
				}
				break
			}
		}
	}

	// Bullet-tank collisions: enemy bullets hit player
	for _, b := range g.Bullets {
		if !b.Active || b.IsPlayer {
			continue
		}
		if g.Player.Alive && !g.Player.IsInvulnerable() {
			if system.BulletTankCollision(b, &g.Player.Tank) {
				b.Active = false
				g.Particles.SpawnExplosion(g.Player.CenterX(), g.Player.CenterY(), 30)
				g.Player.Die()
			}
		}
	}

	// Check eagle destroyed
	eagleDestroyed := false
	for y := 0; y < config.GridHeight; y++ {
		for x := 0; x < config.GridWidth; x++ {
			if g.Grid.Get(x, y) == world.TileEagleDead {
				eagleDestroyed = true
			}
		}
	}

	// Game over conditions
	if eagleDestroyed || (g.Player.Lives <= 0 && !g.Player.Alive) {
		g.State = StateGameOver
	}

	// Level complete: all enemies spawned and destroyed
	if g.Spawner.Done() && g.countAliveEnemies() == 0 {
		g.State = StateLevelComplete
	}

	// Clean up
	g.cleanBullets()
	g.cleanEnemies()
	g.Particles.Update(dt)

	// Level switching (debug)
	if g.Input.IsJustPressed(glow.KeyN) {
		next := g.Level + 1
		if next < len(world.Levels) {
			g.LoadLevel(next)
		}
	}
	if g.Input.IsJustPressed(glow.KeyP) {
		prev := g.Level - 1
		if prev >= 0 {
			g.LoadLevel(prev)
		}
	}
}

func (g *Game) enemyBBoxes() []system.BBox {
	boxes := make([]system.BBox, 0, len(g.Enemies))
	for _, e := range g.Enemies {
		if e.Alive {
			boxes = append(boxes, system.TankBBox(&e.Tank))
		}
	}
	return boxes
}

func (g *Game) tankBBoxesExcluding(self *entity.Tank) []system.BBox {
	boxes := make([]system.BBox, 0, len(g.Enemies)+1)
	if g.Player.Alive {
		boxes = append(boxes, system.TankBBox(&g.Player.Tank))
	}
	for _, e := range g.Enemies {
		if e.Alive && &e.Tank != self {
			boxes = append(boxes, system.TankBBox(&e.Tank))
		}
	}
	return boxes
}

func (g *Game) countAliveEnemies() int {
	count := 0
	for _, e := range g.Enemies {
		if e.Alive {
			count++
		}
	}
	return count
}

func (g *Game) cleanBullets() {
	n := 0
	for _, b := range g.Bullets {
		if b.Active {
			g.Bullets[n] = b
			n++
		}
	}
	g.Bullets = g.Bullets[:n]
}

func (g *Game) cleanEnemies() {
	n := 0
	for _, e := range g.Enemies {
		if e.Alive {
			g.Enemies[n] = e
			n++
		}
	}
	g.Enemies = g.Enemies[:n]
}

func enemyColors(typ entity.EnemyType) render.TankColors {
	switch typ {
	case entity.EnemyFast:
		return render.EnemyFastColors
	case entity.EnemyPower:
		return render.EnemyPowerColors
	case entity.EnemyArmour:
		return render.EnemyArmourColors
	default:
		return render.EnemyBasicColors
	}
}

// Draw renders the current game state.
func (g *Game) Draw(canvas *glow.Canvas) {
	g.Renderer.DrawPlayAreaBorder(canvas)
	g.Renderer.DrawGrid(canvas, g.Grid)

	// Draw enemies
	for _, e := range g.Enemies {
		colors := enemyColors(e.Type)
		if e.IsFlashing() {
			colors = render.TankColors{Body: render.ColorWhite, Tread: render.ColorYellow, Dark: render.ColorGray}
		}
		render.DrawTank(canvas, &e.Tank, colors, config.Padding, config.Padding)
	}

	// Draw player
	if g.Player.Alive {
		render.DrawTank(canvas, &g.Player.Tank, render.PlayerColors, config.Padding, config.Padding)
		if g.Player.IsInvulnerable() {
			render.DrawShield(canvas, &g.Player.Tank, config.Padding, config.Padding, g.Time)
		}
	}

	// Draw bullets
	for _, b := range g.Bullets {
		render.DrawBullet(canvas, b, config.Padding, config.Padding)
	}

	// Particles
	g.Particles.Draw(canvas, config.Padding, config.Padding)

	// Forest overlay (on top of everything)
	g.Renderer.DrawForest(canvas, g.Grid)
}
