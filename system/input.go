package system

import "github.com/AchrafSoltani/glow"

// Input tracks keyboard state.
type Input struct {
	Keys     map[glow.Key]bool
	JustDown map[glow.Key]bool // true only on the frame the key was first pressed
	prev     map[glow.Key]bool
}

// NewInput creates a new input tracker.
func NewInput() *Input {
	return &Input{
		Keys:     make(map[glow.Key]bool),
		JustDown: make(map[glow.Key]bool),
		prev:     make(map[glow.Key]bool),
	}
}

// KeyDown registers a key press.
func (inp *Input) KeyDown(key glow.Key) {
	inp.Keys[key] = true
}

// KeyUp registers a key release.
func (inp *Input) KeyUp(key glow.Key) {
	inp.Keys[key] = false
}

// Update should be called once per frame to compute JustDown.
func (inp *Input) Update() {
	for k := range inp.JustDown {
		delete(inp.JustDown, k)
	}
	for k, v := range inp.Keys {
		if v && !inp.prev[k] {
			inp.JustDown[k] = true
		}
	}
	for k := range inp.prev {
		delete(inp.prev, k)
	}
	for k, v := range inp.Keys {
		inp.prev[k] = v
	}
}

// IsJustPressed returns true only on the first frame a key is held.
func (inp *Input) IsJustPressed(key glow.Key) bool {
	return inp.JustDown[key]
}
