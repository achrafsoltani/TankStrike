package game

// GameState represents the current state of the game.
type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
	StateLevelComplete
	StateLevelIntro
)
