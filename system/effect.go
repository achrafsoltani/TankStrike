package system

// ScreenShake manages screen shake effects.
type ScreenShake struct {
	OffsetX  int
	OffsetY  int
	Duration float64
	Intensity float64
	Timer    float64
}

// Trigger starts a screen shake.
func (s *ScreenShake) Trigger(duration, intensity float64) {
	s.Duration = duration
	s.Intensity = intensity
	s.Timer = duration
}

// Update advances the screen shake timer.
func (s *ScreenShake) Update(dt float64) {
	if s.Timer <= 0 {
		s.OffsetX = 0
		s.OffsetY = 0
		return
	}

	s.Timer -= dt
	progress := s.Timer / s.Duration
	magnitude := s.Intensity * progress

	// Simple alternating shake
	frame := int(s.Timer * 30)
	switch frame % 4 {
	case 0:
		s.OffsetX = int(magnitude)
		s.OffsetY = 0
	case 1:
		s.OffsetX = 0
		s.OffsetY = int(magnitude)
	case 2:
		s.OffsetX = -int(magnitude)
		s.OffsetY = 0
	case 3:
		s.OffsetX = 0
		s.OffsetY = -int(magnitude)
	}
}

// IsActive returns whether the shake is currently active.
func (s *ScreenShake) IsActive() bool {
	return s.Timer > 0
}
