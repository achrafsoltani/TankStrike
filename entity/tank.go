package entity

// Tank is the base struct for all tanks (player and enemy).
type Tank struct {
	X, Y      float64   // pixel position (top-left of 2x2 area)
	Dir       Direction  // facing direction
	Speed     float64    // pixels per second
	HP        int        // hit points
	MaxHP     int        // max hit points
	Alive     bool
	Moving    bool       // whether the tank is currently moving

	// Animation
	TreadFrame int     // alternates for tread animation
	TreadTimer float64 // time accumulator for tread animation

	// Shooting
	ShootCooldown float64 // time remaining before next shot
	CooldownRate  float64 // seconds between shots
	BulletSpeed   float64
	PowerLevel    int // star upgrades (0-3), affects bullet power
}

// NewTank creates a base tank with default values.
func NewTank(x, y float64, speed float64, hp int) Tank {
	return Tank{
		X:            x,
		Y:            y,
		Dir:          DirUp,
		Speed:        speed,
		HP:           hp,
		MaxHP:        hp,
		Alive:        true,
		CooldownRate: 0.5,
		BulletSpeed:  300,
	}
}

// Update updates tank animation state.
func (t *Tank) Update(dt float64) {
	if t.ShootCooldown > 0 {
		t.ShootCooldown -= dt
	}

	if t.Moving {
		t.TreadTimer += dt
		if t.TreadTimer > 0.08 {
			t.TreadTimer = 0
			t.TreadFrame = (t.TreadFrame + 1) % 4
		}
	}
}

// CanShoot returns whether the tank can fire.
func (t *Tank) CanShoot() bool {
	return t.Alive && t.ShootCooldown <= 0
}

// Shoot puts the tank on cooldown. Returns the bullet spawn position.
func (t *Tank) Shoot() (float64, float64) {
	t.ShootCooldown = t.CooldownRate

	// Bullet spawns at the barrel tip
	cx := t.X + 24 // centre of 48px tank
	cy := t.Y + 24
	bx := cx + t.Dir.DX()*28 - 2 // offset to barrel tip, centred on 4px bullet
	by := cy + t.Dir.DY()*28 - 2
	return bx, by
}

// CenterX returns the centre X of the tank.
func (t *Tank) CenterX() float64 { return t.X + 24 }

// CenterY returns the centre Y of the tank.
func (t *Tank) CenterY() float64 { return t.Y + 24 }

// Hit reduces HP. Returns true if the tank is destroyed.
func (t *Tank) Hit(damage int) bool {
	t.HP -= damage
	if t.HP <= 0 {
		t.HP = 0
		t.Alive = false
		return true
	}
	return false
}
