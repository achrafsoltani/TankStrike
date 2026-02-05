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
	Eagle     *entity.Eagle
	Enemies   []*entity.EnemyTank
	Bullets   []*entity.Bullet
	Particles *render.ParticlePool
	Spawner   *system.Spawner
	Level     int
	Time      float64

	// Transition timers
	GameOverTimer    float64
	LevelComplTimer  float64
	LevelIntroTimer  float64
}

// NewGame creates a new game instance.
func NewGame() *Game {
	g := &Game{
		State:     StateMenu,
		Grid:      world.NewGrid(),
		Renderer:  render.NewRenderer(),
		Input:     system.NewInput(),
		Player:    entity.NewPlayerTank(),
		Particles: render.NewParticlePool(),
		Level:     0,
	}
	return g
}

// StartGame begins a new game from level 0.
func (g *Game) StartGame() {
	g.Player = entity.NewPlayerTank()
	g.Level = 0
	g.startLevel(0)
}

func (g *Game) startLevel(index int) {
	if index >= 0 && index < len(world.Levels) {
		g.Level = index
		world.LoadLevel(g.Grid, world.Levels[index])
		g.Bullets = g.Bullets[:0]
		g.Enemies = g.Enemies[:0]
		g.Spawner = system.NewSpawner(index)
		g.findEagle()
		g.Player.Respawn()
		g.State = StateLevelIntro
		g.LevelIntroTimer = 2.0
	}
}

func (g *Game) findEagle() {
	for y := 0; y < config.GridHeight; y++ {
		for x := 0; x < config.GridWidth; x++ {
			if g.Grid.Get(x, y) == world.TileEagle {
				g.Eagle = entity.NewEagle(x, y)
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
	case StateMenu:
		if g.Input.IsJustPressed(glow.KeyEnter) || g.Input.IsJustPressed(glow.KeySpace) {
			g.StartGame()
		}
	case StateLevelIntro:
		g.LevelIntroTimer -= dt
		if g.LevelIntroTimer <= 0 {
			g.State = StatePlaying
		}
	case StatePlaying:
		g.updatePlaying(dt)
		if g.Input.IsJustPressed(glow.KeyEscape) {
			g.State = StatePaused
		}
	case StatePaused:
		if g.Input.IsJustPressed(glow.KeyEscape) || g.Input.IsJustPressed(glow.KeyEnter) {
			g.State = StatePlaying
		}
	case StateGameOver:
		g.GameOverTimer -= dt
		g.Particles.Update(dt)
		if g.GameOverTimer <= 0 {
			if g.Input.IsJustPressed(glow.KeyEnter) || g.Input.IsJustPressed(glow.KeySpace) {
				g.State = StateMenu
			}
		}
	case StateLevelComplete:
		g.LevelComplTimer -= dt
		if g.LevelComplTimer <= 0 {
			if g.Input.IsJustPressed(glow.KeyEnter) || g.Input.IsJustPressed(glow.KeySpace) {
				next := g.Level + 1
				if next < len(world.Levels) {
					g.startLevel(next)
				} else {
					g.State = StateMenu // Beat the game
				}
			}
		}
	}
}

func (g *Game) updatePlaying(dt float64) {
	// Player input and movement
	g.Player.HandleInput(g.Input.Keys)
	g.Player.UpdatePlayer(dt)

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
	if enemy := g.Spawner.Update(dt, g.countAliveEnemies()); enemy != nil {
		g.Enemies = append(g.Enemies, enemy)
	}

	// Update enemies
	eagleCX, eagleCY := g.Eagle.CenterX(), g.Eagle.CenterY()
	for _, e := range g.Enemies {
		if !e.Alive {
			continue
		}
		others := g.tankBBoxesExcluding(&e.Tank)
		system.UpdateEnemyAI(e, g.Grid, dt,
			g.Player.CenterX(), g.Player.CenterY(),
			eagleCX, eagleCY, others)

		if system.ShouldShoot(e, dt) {
			bx, by := e.Shoot()
			bullet := entity.NewBullet(bx, by, e.Dir, e.BulletSpeed, 0, false)
			g.Bullets = append(g.Bullets, bullet)
		}
	}

	// Update eagle
	g.Eagle.Update(dt)

	// Update bullets
	for _, b := range g.Bullets {
		b.Update(dt)
	}

	// Bullet-grid collisions
	for _, b := range g.Bullets {
		system.BulletGridCollision(b, g.Grid, g.Particles)
	}

	// Player bullets hit enemies
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

	// Enemy bullets hit player
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
	if g.Eagle != nil {
		for y := 0; y < config.GridHeight; y++ {
			for x := 0; x < config.GridWidth; x++ {
				if g.Grid.Get(x, y) == world.TileEagleDead {
					g.Eagle.Alive = false
				}
			}
		}
	}

	// Game over conditions
	if !g.Eagle.Alive || (g.Player.Lives <= 0 && !g.Player.Alive) {
		g.State = StateGameOver
		g.GameOverTimer = 2.0
	}

	// Level complete
	if g.Spawner.Done() && g.countAliveEnemies() == 0 {
		g.State = StateLevelComplete
		g.LevelComplTimer = 1.5
	}

	// Clean up
	g.cleanBullets()
	g.cleanEnemies()
	g.Particles.Update(dt)

	// Debug level switching
	if g.Input.IsJustPressed(glow.KeyN) {
		next := g.Level + 1
		if next < len(world.Levels) {
			g.startLevel(next)
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
	switch g.State {
	case StateMenu:
		g.drawMenu(canvas)
	case StateLevelIntro:
		g.drawLevelIntro(canvas)
	case StatePlaying, StatePaused:
		g.drawPlayField(canvas)
		if g.State == StatePaused {
			g.drawPauseOverlay(canvas)
		}
	case StateGameOver:
		g.drawPlayField(canvas)
		g.drawGameOver(canvas)
	case StateLevelComplete:
		g.drawPlayField(canvas)
		g.drawLevelComplete(canvas)
	}
}

func (g *Game) drawPlayField(canvas *glow.Canvas) {
	g.Renderer.DrawPlayAreaBorder(canvas)
	g.Renderer.DrawGrid(canvas, g.Grid)

	// Enemies
	for _, e := range g.Enemies {
		colors := enemyColors(e.Type)
		if e.IsFlashing() {
			colors = render.TankColors{Body: render.ColorWhite, Tread: render.ColorYellow, Dark: render.ColorGray}
		}
		render.DrawTank(canvas, &e.Tank, colors, config.Padding, config.Padding)
	}

	// Player
	if g.Player.Alive {
		render.DrawTank(canvas, &g.Player.Tank, render.PlayerColors, config.Padding, config.Padding)
		if g.Player.IsInvulnerable() {
			render.DrawShield(canvas, &g.Player.Tank, config.Padding, config.Padding, g.Time)
		}
	}

	// Bullets
	for _, b := range g.Bullets {
		render.DrawBullet(canvas, b, config.Padding, config.Padding)
	}

	// Particles
	g.Particles.Draw(canvas, config.Padding, config.Padding)

	// Forest overlay
	g.Renderer.DrawForest(canvas, g.Grid)
}

// Minimal text rendering using rectangles (replaced by bitmap font in Phase 6)
func drawSimpleText(canvas *glow.Canvas, text string, x, y int, color glow.Color, scale int) {
	// Placeholder: draw a coloured rectangle as text indicator
	w := len(text) * 6 * scale
	canvas.DrawRect(x, y, w, 8*scale, color)
}

func (g *Game) drawMenu(canvas *glow.Canvas) {
	// Simple placeholder menu
	cx := config.WindowWidth / 2
	cy := config.WindowHeight / 2

	// Title box
	canvas.DrawRect(cx-120, cy-80, 240, 40, render.ColorYellow)

	// "Press Enter" indicator
	if int(g.Time*2)%2 == 0 {
		canvas.DrawRect(cx-80, cy+20, 160, 20, render.ColorWhite)
	}
}

func (g *Game) drawLevelIntro(canvas *glow.Canvas) {
	cx := config.WindowWidth / 2
	cy := config.WindowHeight / 2
	canvas.DrawRect(cx-60, cy-12, 120, 24, render.ColorGray)
}

func (g *Game) drawPauseOverlay(canvas *glow.Canvas) {
	// Dithered checkerboard pattern for semi-transparency
	for y := 0; y < config.WindowHeight; y += 2 {
		for x := 0; x < config.WindowWidth; x += 2 {
			canvas.SetPixel(x, y, render.ColorBlack)
		}
	}
	cx := config.WindowWidth / 2
	cy := config.WindowHeight / 2
	canvas.DrawRect(cx-50, cy-12, 100, 24, render.ColorYellow)
}

func (g *Game) drawGameOver(canvas *glow.Canvas) {
	// Dithered overlay
	for y := 0; y < config.WindowHeight; y += 2 {
		for x := (y / 2) % 2; x < config.WindowWidth; x += 2 {
			canvas.SetPixel(x, y, render.ColorBlack)
		}
	}
	cx := config.WindowWidth / 2
	cy := config.WindowHeight / 2
	canvas.DrawRect(cx-70, cy-16, 140, 32, render.ColorRed)
}

func (g *Game) drawLevelComplete(canvas *glow.Canvas) {
	cx := config.WindowWidth / 2
	cy := config.WindowHeight / 2
	canvas.DrawRect(cx-80, cy-16, 160, 32, render.ColorPlayerBody)
}
