package config

// Layout holds scale and centering offsets for window resize support.
type Layout struct {
	Scale   float64
	OffsetX int
	OffsetY int
}

// NewLayout computes a layout that scales the original WindowWidth x WindowHeight
// content to fit within winW x winH, centred with letterboxing/pillarboxing.
func NewLayout(winW, winH int) Layout {
	sx := float64(winW) / float64(WindowWidth)
	sy := float64(winH) / float64(WindowHeight)
	scale := sx
	if sy < sx {
		scale = sy
	}
	if scale < 1.0 {
		scale = 1.0
	}

	scaledW := int(float64(WindowWidth) * scale)
	scaledH := int(float64(WindowHeight) * scale)
	ox := (winW - scaledW) / 2
	oy := (winH - scaledH) / 2

	return Layout{
		Scale:   scale,
		OffsetX: ox,
		OffsetY: oy,
	}
}
