package game

import (
	"github.com/AchrafSoltani/TankStrike/render"
	"github.com/AchrafSoltani/TankStrike/world"
	"github.com/AchrafSoltani/glow"
)

// Game is the top-level game orchestrator.
type Game struct {
	State    GameState
	Grid     *world.Grid
	Renderer *render.Renderer
	Level    int
	Time     float64

	// Input state
	Keys map[glow.Key]bool
}

// NewGame creates a new game instance.
func NewGame() *Game {
	g := &Game{
		State:    StatePlaying,
		Grid:     world.NewGrid(),
		Renderer: render.NewRenderer(),
		Level:    0,
		Keys:     make(map[glow.Key]bool),
	}
	g.LoadLevel(0)
	return g
}

// LoadLevel loads a level by index.
func (g *Game) LoadLevel(index int) {
	if index >= 0 && index < len(world.Levels) {
		g.Level = index
		world.LoadLevel(g.Grid, world.Levels[index])
	}
}

// KeyDown handles key press events.
func (g *Game) KeyDown(key glow.Key) {
	g.Keys[key] = true
}

// KeyUp handles key release events.
func (g *Game) KeyUp(key glow.Key) {
	g.Keys[key] = false
}

// Update advances game state by dt seconds.
func (g *Game) Update(dt float64) {
	g.Time += dt
	g.Renderer.Time = g.Time
}

// Draw renders the current game state.
func (g *Game) Draw(canvas *glow.Canvas) {
	g.Renderer.DrawPlayAreaBorder(canvas)
	g.Renderer.DrawGrid(canvas, g.Grid)
}
