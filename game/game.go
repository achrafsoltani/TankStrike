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
	HUD       *render.HUDRenderer
	Input     *system.Input
	Player    *entity.PlayerTank
	Eagle     *entity.Eagle
	Enemies   []*entity.EnemyTank
	Bullets   []*entity.Bullet
	PowerUps  []*entity.PowerUp
	Particles *render.ParticlePool
	Spawner   *system.Spawner
	Level     int
	Time      float64

	// Power-up timers
	ClockTimer  float64 // freeze enemies timer
	ShovelTimer float64 // fortified eagle timer

	// Level stats
	KillsBasic  int
	KillsFast   int
	KillsPower  int
	KillsArmour int

	// Menu state
	MenuSelection int
	MenuOptions   []render.MenuOption

	// Transition timers
	GameOverTimer   float64
	LevelComplTimer float64
	LevelIntroTimer float64
}

// NewGame creates a new game instance.
func NewGame() *Game {
	g := &Game{
		State:     StateMenu,
		Grid:      world.NewGrid(),
		Renderer:  render.NewRenderer(),
		HUD:       render.NewHUDRenderer(),
		Input:     system.NewInput(),
		Player:    entity.NewPlayerTank(),
		Particles: render.NewParticlePool(),
		Level:     0,
		MenuOptions: []render.MenuOption{
			{Label: "NEW GAME"},
			{Label: "CONTINUE"},
		},
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
		g.PowerUps = g.PowerUps[:0]
		g.ClockTimer = 0
		g.ShovelTimer = 0
		g.KillsBasic = 0
		g.KillsFast = 0
		g.KillsPower = 0
		g.KillsArmour = 0
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
		if g.Input.IsJustPressed(glow.KeyUp) || g.Input.IsJustPressed(glow.KeyW) {
			g.MenuSelection--
			if g.MenuSelection < 0 {
				g.MenuSelection = len(g.MenuOptions) - 1
			}
		}
		if g.Input.IsJustPressed(glow.KeyDown) || g.Input.IsJustPressed(glow.KeyS) {
			g.MenuSelection++
			if g.MenuSelection >= len(g.MenuOptions) {
				g.MenuSelection = 0
			}
		}
		if g.Input.IsJustPressed(glow.KeyEnter) || g.Input.IsJustPressed(glow.KeySpace) {
			switch g.MenuSelection {
			case 0: // New Game
				g.StartGame()
			case 1: // Continue (loads save or starts new)
				g.StartGame()
			}
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
					g.State = StateMenu
				}
			}
		}
	}
}

func (g *Game) updatePlaying(dt float64) {
	g.Player.HandleInput(g.Input.Keys)
	g.Player.UpdatePlayer(dt)

	otherTanks := g.enemyBBoxes()
	system.MovePlayerTank(g.Player, g.Grid, dt, otherTanks)

	if g.Player.WantsToShoot(g.Input.Keys) && g.Player.CanShoot() {
		if system.CountPlayerBullets(g.Bullets) < config.MaxPlayerBullets {
			bx, by := g.Player.Shoot()
			bullet := entity.NewBullet(bx, by, g.Player.Dir, g.Player.BulletSpeed, g.Player.PowerLevel, true)
			g.Bullets = append(g.Bullets, bullet)
		}
	}

	if enemy := g.Spawner.Update(dt, g.countAliveEnemies()); enemy != nil {
		g.Enemies = append(g.Enemies, enemy)
	}

	eagleCX, eagleCY := g.Eagle.CenterX(), g.Eagle.CenterY()
	frozen := g.ClockTimer > 0
	for _, e := range g.Enemies {
		if !e.Alive {
			continue
		}
		if frozen {
			e.UpdateEnemy(dt) // still animate flash, but don't move/shoot
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

	g.Eagle.Update(dt)

	for _, b := range g.Bullets {
		b.Update(dt)
	}

	for _, b := range g.Bullets {
		system.BulletGridCollision(b, g.Grid, g.Particles)
	}

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
					g.trackKill(e.Type)
					if e.HasPowerUp {
						g.PowerUps = append(g.PowerUps, entity.NewPowerUp())
					}
				}
				break
			}
		}
	}

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

	// Power-up collection
	if g.Player.Alive {
		for _, p := range g.PowerUps {
			if !p.Active {
				continue
			}
			// Simple AABB overlap between player and power-up
			if g.Player.X < p.X+24 && g.Player.X+float64(config.TankSize) > p.X &&
				g.Player.Y < p.Y+24 && g.Player.Y+float64(config.TankSize) > p.Y {
				p.Active = false
				g.applyPowerUp(p.Type)
			}
		}
	}

	// Update power-ups
	for _, p := range g.PowerUps {
		p.Update(dt)
	}

	// Clock timer (freeze enemies)
	if g.ClockTimer > 0 {
		g.ClockTimer -= dt
	}

	// Shovel timer (fortification)
	if g.ShovelTimer > 0 {
		g.ShovelTimer -= dt
		if g.ShovelTimer <= 0 {
			g.unfortifyEagle()
		}
	}

	// Clean up power-ups
	g.cleanPowerUps()

	if g.Eagle != nil {
		for y := 0; y < config.GridHeight; y++ {
			for x := 0; x < config.GridWidth; x++ {
				if g.Grid.Get(x, y) == world.TileEagleDead {
					g.Eagle.Alive = false
				}
			}
		}
	}

	if !g.Eagle.Alive || (g.Player.Lives <= 0 && !g.Player.Alive) {
		g.State = StateGameOver
		g.GameOverTimer = 2.0
	}

	if g.Spawner.Done() && g.countAliveEnemies() == 0 {
		g.State = StateLevelComplete
		g.LevelComplTimer = 1.5
	}

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

func (g *Game) trackKill(typ entity.EnemyType) {
	switch typ {
	case entity.EnemyBasic:
		g.KillsBasic++
	case entity.EnemyFast:
		g.KillsFast++
	case entity.EnemyPower:
		g.KillsPower++
	case entity.EnemyArmour:
		g.KillsArmour++
	}
}

func (g *Game) cleanPowerUps() {
	n := 0
	for _, p := range g.PowerUps {
		if p.Active {
			g.PowerUps[n] = p
			n++
		}
	}
	g.PowerUps = g.PowerUps[:n]
}

func (g *Game) applyPowerUp(typ entity.PowerUpType) {
	switch typ {
	case entity.PowerUpStar:
		g.Player.ApplyStar()
	case entity.PowerUpTank:
		g.Player.Lives++
	case entity.PowerUpHelmet:
		g.Player.ShieldTimer = config.PowerUpDuration
	case entity.PowerUpShovel:
		g.fortifyEagle()
		g.ShovelTimer = config.PowerUpDuration
	case entity.PowerUpBomb:
		for _, e := range g.Enemies {
			if e.Alive {
				e.Alive = false
				g.Player.Score += e.ScoreValue
				g.Particles.SpawnExplosion(e.CenterX(), e.CenterY(), 25)
			}
		}
	case entity.PowerUpClock:
		g.ClockTimer = config.PowerUpDuration
	}
	g.Player.Score += 500
}

func (g *Game) fortifyEagle() {
	if g.Eagle == nil {
		return
	}
	// Replace brick around eagle with steel
	ex := int(g.Eagle.X) / config.SubBlock
	ey := int(g.Eagle.Y) / config.SubBlock
	for dy := -1; dy <= 2; dy++ {
		for dx := -1; dx <= 2; dx++ {
			gx, gy := ex+dx, ey+dy
			tile := g.Grid.Get(gx, gy)
			if tile == world.TileBrick || tile == world.TileEmpty {
				// Only fortify the border cells
				if dx == -1 || dx == 2 || dy == -1 || dy == 2 {
					g.Grid.Set(gx, gy, world.TileSteel)
				}
			}
		}
	}
	g.Eagle.Fortified = true
}

func (g *Game) unfortifyEagle() {
	if g.Eagle == nil {
		return
	}
	ex := int(g.Eagle.X) / config.SubBlock
	ey := int(g.Eagle.Y) / config.SubBlock
	for dy := -1; dy <= 2; dy++ {
		for dx := -1; dx <= 2; dx++ {
			gx, gy := ex+dx, ey+dy
			tile := g.Grid.Get(gx, gy)
			if tile == world.TileSteel {
				if dx == -1 || dx == 2 || dy == -1 || dy == 2 {
					g.Grid.Set(gx, gy, world.TileBrick)
				}
			}
		}
	}
	g.Eagle.Fortified = false
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
		g.drawHUD(canvas)
		if g.State == StatePaused {
			g.drawPauseOverlay(canvas)
		}
	case StateGameOver:
		g.drawPlayField(canvas)
		g.drawHUD(canvas)
		g.drawGameOver(canvas)
	case StateLevelComplete:
		g.drawPlayField(canvas)
		g.drawHUD(canvas)
		g.drawLevelComplete(canvas)
	}
}

func (g *Game) drawPlayField(canvas *glow.Canvas) {
	g.Renderer.DrawPlayAreaBorder(canvas)
	g.Renderer.DrawGrid(canvas, g.Grid)

	for _, e := range g.Enemies {
		colors := enemyColors(e.Type)
		if e.IsFlashing() {
			colors = render.TankColors{Body: render.ColorWhite, Tread: render.ColorYellow, Dark: render.ColorGray}
		}
		render.DrawTank(canvas, &e.Tank, colors, config.Padding, config.Padding)
	}

	if g.Player.Alive {
		render.DrawTank(canvas, &g.Player.Tank, render.PlayerColors, config.Padding, config.Padding)
		if g.Player.IsInvulnerable() {
			render.DrawShield(canvas, &g.Player.Tank, config.Padding, config.Padding, g.Time)
		}
	}

	for _, b := range g.Bullets {
		render.DrawBullet(canvas, b, config.Padding, config.Padding)
	}

	// Power-ups
	for _, p := range g.PowerUps {
		render.DrawPowerUp(canvas, p, config.Padding, config.Padding)
	}

	g.Particles.Draw(canvas, config.Padding, config.Padding)
	g.Renderer.DrawForest(canvas, g.Grid)
}

func (g *Game) drawHUD(canvas *glow.Canvas) {
	remaining := g.Spawner.Remaining() + g.countAliveEnemies()
	g.HUD.DrawHUD(canvas, remaining, g.Player.Lives, g.Level, g.Player.Score)
}

func (g *Game) drawMenu(canvas *glow.Canvas) {
	render.DrawTitleScreen(canvas, g.MenuOptions, g.MenuSelection, g.Time)
}

func (g *Game) drawLevelIntro(canvas *glow.Canvas) {
	render.DrawLevelIntro(canvas, g.Level)
}

func (g *Game) drawPauseOverlay(canvas *glow.Canvas) {
	render.DrawPauseScreen(canvas, g.Time)
}

func (g *Game) drawGameOver(canvas *glow.Canvas) {
	render.DrawGameOverScreen(canvas, g.Player.Score, g.GameOverTimer <= 0, g.Time)
}

func (g *Game) drawLevelComplete(canvas *glow.Canvas) {
	render.DrawLevelComplete(canvas, g.Level, g.Player.Score,
		g.KillsBasic, g.KillsFast, g.KillsPower, g.KillsArmour,
		g.LevelComplTimer <= 0, g.Time)
}
