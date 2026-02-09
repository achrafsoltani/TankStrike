package render

import (
	"github.com/AchrafSoltani/TankStrike/config"
	"github.com/AchrafSoltani/TankStrike/entity"
)

// DrawBullet draws a bullet with its trail.
func DrawBullet(canvas *ScaledCanvas, b *entity.Bullet, offsetX, offsetY int) {
	if !b.Active {
		return
	}

	px := int(b.X) + offsetX
	py := int(b.Y) + offsetY

	// Trail
	for i := 0; i < b.TrailCount; i++ {
		tx := int(b.TrailX[i]) + offsetX
		ty := int(b.TrailY[i]) + offsetY
		size := 2
		canvas.DrawRect(tx+1, ty+1, size, size, ColorBulletTrail)
	}

	// Bullet body
	color := ColorBulletPlayer
	if !b.IsPlayer {
		color = ColorBulletEnemy
	}
	canvas.DrawRect(px, py, config.BulletSize, config.BulletSize, color)
}
