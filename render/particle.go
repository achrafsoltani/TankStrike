package render

import (
	"math"
	"math/rand"

	"github.com/AchrafSoltani/TankStrike/config"
	"github.com/AchrafSoltani/glow"
)

// Particle is a single visual effect particle.
type Particle struct {
	X, Y     float64
	VX, VY   float64
	Life     float64
	MaxLife  float64
	Size     float64
	Color    glow.Color
	IsCircle bool
	Active   bool
}

// ParticlePool manages a fixed-size pool of particles.
type ParticlePool struct {
	Particles [config.MaxParticles]Particle
}

// NewParticlePool creates a new particle pool.
func NewParticlePool() *ParticlePool {
	return &ParticlePool{}
}

// Emit activates a particle with the given properties.
func (pp *ParticlePool) Emit(x, y, vx, vy, life, size float64, color glow.Color, isCircle bool) {
	for i := range pp.Particles {
		if !pp.Particles[i].Active {
			pp.Particles[i] = Particle{
				X: x, Y: y, VX: vx, VY: vy,
				Life: life, MaxLife: life,
				Size: size, Color: color,
				IsCircle: isCircle, Active: true,
			}
			return
		}
	}
}

// SpawnExplosion creates a burst of explosion particles.
func (pp *ParticlePool) SpawnExplosion(x, y float64, count int) {
	colors := []glow.Color{ColorExplosion1, ColorExplosion2, ColorExplosion3, ColorExplosion4}
	for i := 0; i < count; i++ {
		angle := rand.Float64() * math.Pi * 2
		speed := 40 + rand.Float64()*120
		vx := math.Cos(angle) * speed
		vy := math.Sin(angle) * speed
		life := 0.3 + rand.Float64()*0.5
		size := 2 + rand.Float64()*4
		color := colors[rand.Intn(len(colors))]
		isCircle := rand.Float64() > 0.5
		pp.Emit(x, y, vx, vy, life, size, color, isCircle)
	}
}

// SpawnSpark creates small sparks (for bullet hitting steel).
func (pp *ParticlePool) SpawnSpark(x, y float64) {
	for i := 0; i < 8; i++ {
		angle := rand.Float64() * math.Pi * 2
		speed := 30 + rand.Float64()*80
		vx := math.Cos(angle) * speed
		vy := math.Sin(angle) * speed
		life := 0.1 + rand.Float64()*0.2
		pp.Emit(x, y, vx, vy, life, 2, ColorSpark, false)
	}
}

// SpawnDebris creates brick debris particles.
func (pp *ParticlePool) SpawnDebris(x, y float64) {
	colors := []glow.Color{ColorDebris1, ColorDebris2, ColorBrick}
	for i := 0; i < 12; i++ {
		angle := rand.Float64() * math.Pi * 2
		speed := 20 + rand.Float64()*60
		vx := math.Cos(angle) * speed
		vy := math.Sin(angle)*speed - 20
		life := 0.3 + rand.Float64()*0.4
		size := 2 + rand.Float64()*3
		color := colors[rand.Intn(len(colors))]
		pp.Emit(x, y, vx, vy, life, size, color, false)
	}
}

// Update advances all active particles.
func (pp *ParticlePool) Update(dt float64) {
	for i := range pp.Particles {
		p := &pp.Particles[i]
		if !p.Active {
			continue
		}
		p.Life -= dt
		if p.Life <= 0 {
			p.Active = false
			continue
		}
		p.X += p.VX * dt
		p.Y += p.VY * dt
		p.VY += 100 * dt // gravity
		// Shrink over time
		ratio := p.Life / p.MaxLife
		p.Size *= (0.98 + 0.02*ratio)
	}
}

// Draw renders all active particles.
func (pp *ParticlePool) Draw(canvas *ScaledCanvas, offsetX, offsetY int) {
	for i := range pp.Particles {
		p := &pp.Particles[i]
		if !p.Active {
			continue
		}
		px := int(p.X) + offsetX
		py := int(p.Y) + offsetY
		s := int(p.Size)
		if s < 1 {
			s = 1
		}
		if p.IsCircle {
			canvas.FillCircle(px, py, s, p.Color)
		} else {
			canvas.DrawRect(px-s/2, py-s/2, s, s, p.Color)
		}
	}
}
