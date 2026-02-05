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
	Bullets   []*entity.Bullet
	Particles *render.ParticlePool
	Level     int
	Time      float64
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
	// Player input
	g.Player.HandleInput(g.Input.Keys)
	g.Player.UpdatePlayer(dt)

	// Move player
	var otherTanks []system.BBox
	system.MovePlayerTank(g.Player, g.Grid, dt, otherTanks)

	// Player shooting
	if g.Player.WantsToShoot(g.Input.Keys) && g.Player.CanShoot() {
		if system.CountPlayerBullets(g.Bullets) < config.MaxPlayerBullets {
			bx, by := g.Player.Shoot()
			bullet := entity.NewBullet(bx, by, g.Player.Dir, g.Player.BulletSpeed, g.Player.PowerLevel, true)
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

	// Clean up inactive bullets
	g.cleanBullets()

	// Update particles
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

// Draw renders the current game state.
func (g *Game) Draw(canvas *glow.Canvas) {
	g.Renderer.DrawPlayAreaBorder(canvas)
	g.Renderer.DrawGrid(canvas, g.Grid)

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
