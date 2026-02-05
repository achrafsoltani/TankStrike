package entity

// Direction represents a facing direction.
type Direction int

const (
	DirUp Direction = iota
	DirDown
	DirLeft
	DirRight
)

// DX returns the X component of the direction vector.
func (d Direction) DX() float64 {
	switch d {
	case DirLeft:
		return -1
	case DirRight:
		return 1
	default:
		return 0
	}
}

// DY returns the Y component of the direction vector.
func (d Direction) DY() float64 {
	switch d {
	case DirUp:
		return -1
	case DirDown:
		return 1
	default:
		return 0
	}
}

// Opposite returns the opposite direction.
func (d Direction) Opposite() Direction {
	switch d {
	case DirUp:
		return DirDown
	case DirDown:
		return DirUp
	case DirLeft:
		return DirRight
	case DirRight:
		return DirLeft
	default:
		return d
	}
}
