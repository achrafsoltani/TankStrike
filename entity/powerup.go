package entity

import "math/rand"

// PowerUpType represents the type of power-up.
type PowerUpType int

const (
	PowerUpStar   PowerUpType = iota // Upgrade tank fire power
	PowerUpTank                      // Extra life
	PowerUpHelmet                    // Temporary invulnerability
	PowerUpShovel                    // Fortify eagle with steel
	PowerUpBomb                      // Destroy all active enemies
	PowerUpClock                     // Freeze all enemies
)

// PowerUp represents a collectible power-up on the field.
type PowerUp struct {
	X, Y        float64
	Type        PowerUpType
	Active      bool
	FlashTimer  float64
}

// NewPowerUp creates a new power-up at a random position within the play area.
func NewPowerUp() *PowerUp {
	types := []PowerUpType{PowerUpStar, PowerUpTank, PowerUpHelmet, PowerUpShovel, PowerUpBomb, PowerUpClock}
	typ := types[rand.Intn(len(types))]

	// Random position, snapped to sub-block grid, avoiding edges
	x := float64(2+rand.Intn(22)) * 24
	y := float64(2+rand.Intn(22)) * 24

	return &PowerUp{
		X:      x,
		Y:      y,
		Type:   typ,
		Active: true,
	}
}

// Update handles the power-up flash timer.
func (p *PowerUp) Update(dt float64) {
	p.FlashTimer += dt
}

// IsVisible returns whether the power-up is currently visible (flashing).
func (p *PowerUp) IsVisible() bool {
	return int(p.FlashTimer*4)%2 == 0
}

// TypeName returns a display name for the power-up type.
func (p *PowerUp) TypeName() string {
	switch p.Type {
	case PowerUpStar:
		return "STAR"
	case PowerUpTank:
		return "TANK"
	case PowerUpHelmet:
		return "HELMET"
	case PowerUpShovel:
		return "SHOVEL"
	case PowerUpBomb:
		return "BOMB"
	case PowerUpClock:
		return "CLOCK"
	default:
		return ""
	}
}
