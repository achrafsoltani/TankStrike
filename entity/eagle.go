package entity

import "github.com/AchrafSoltani/TankStrike/config"

// Eagle represents the player's base that must be protected.
type Eagle struct {
	X, Y      float64 // pixel position (top-left of 2x2 area)
	Alive     bool
	Fortified bool    // whether surrounded by steel (Shovel power-up)
	FortTimer float64 // remaining fortification time
}

// NewEagle creates an eagle at the standard position.
// It scans the grid to find the eagle tile position.
func NewEagle(gridX, gridY int) *Eagle {
	return &Eagle{
		X:     float64(gridX * config.SubBlock),
		Y:     float64(gridY * config.SubBlock),
		Alive: true,
	}
}

// CenterX returns the centre X pixel position.
func (e *Eagle) CenterX() float64 {
	return e.X + float64(config.SubBlock)
}

// CenterY returns the centre Y pixel position.
func (e *Eagle) CenterY() float64 {
	return e.Y + float64(config.SubBlock)
}

// Update handles eagle state updates.
func (e *Eagle) Update(dt float64) {
	if e.Fortified {
		e.FortTimer -= dt
		if e.FortTimer <= 0 {
			e.Fortified = false
		}
	}
}
