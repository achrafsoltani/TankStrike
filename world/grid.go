package world

import "github.com/AchrafSoltani/TankStrike/config"

// Grid represents the 26x26 sub-block game world.
type Grid struct {
	Tiles [config.GridHeight][config.GridWidth]TileType
}

// NewGrid creates an empty grid.
func NewGrid() *Grid {
	return &Grid{}
}

// Get returns the tile type at the given sub-block position.
func (g *Grid) Get(x, y int) TileType {
	if x < 0 || x >= config.GridWidth || y < 0 || y >= config.GridHeight {
		return TileSteel // Out of bounds is impassable
	}
	return g.Tiles[y][x]
}

// Set places a tile at the given sub-block position.
func (g *Grid) Set(x, y int, t TileType) {
	if x >= 0 && x < config.GridWidth && y >= 0 && y < config.GridHeight {
		g.Tiles[y][x] = t
	}
}

// Destroy removes a tile (sets to empty) if destructible. Returns true if destroyed.
func (g *Grid) Destroy(x, y int, powerLevel int) bool {
	if x < 0 || x >= config.GridWidth || y < 0 || y >= config.GridHeight {
		return false
	}
	t := g.Tiles[y][x]
	switch t {
	case TileBrick:
		g.Tiles[y][x] = TileEmpty
		return true
	case TileSteel:
		if powerLevel >= 3 {
			g.Tiles[y][x] = TileEmpty
			return true
		}
		return false
	case TileEagle:
		g.Tiles[y][x] = TileEagleDead
		return true
	default:
		return false
	}
}

// IsPassable checks if a 2x2 tank footprint can occupy the given sub-block position.
// (x, y) is the top-left sub-block of the 2x2 area.
func (g *Grid) IsPassable(x, y int) bool {
	for dy := 0; dy < 2; dy++ {
		for dx := 0; dx < 2; dx++ {
			if !g.Get(x+dx, y+dy).IsPassable() {
				return false
			}
		}
	}
	return true
}

// Clear resets all tiles to empty.
func (g *Grid) Clear() {
	for y := 0; y < config.GridHeight; y++ {
		for x := 0; x < config.GridWidth; x++ {
			g.Tiles[y][x] = TileEmpty
		}
	}
}

// GetTileAt returns the tile type at a pixel position within the play area.
func (g *Grid) GetTileAt(px, py int) TileType {
	sx := px / config.SubBlock
	sy := py / config.SubBlock
	return g.Get(sx, sy)
}
