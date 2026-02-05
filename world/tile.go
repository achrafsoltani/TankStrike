package world

// TileType represents a type of tile in the game world.
type TileType int

const (
	TileEmpty  TileType = iota
	TileBrick           // Destructible wall
	TileSteel           // Indestructible wall (unless 3-star)
	TileWater           // Impassable to tanks, bullets pass over
	TileIce             // Tanks slide on ice
	TileForest          // Visual overlay, tanks pass under
	TileEagle           // Player base (must protect)
	TileEagleDead       // Destroyed eagle
)

// IsPassable returns whether tanks can move through this tile.
func (t TileType) IsPassable() bool {
	switch t {
	case TileEmpty, TileIce, TileForest:
		return true
	default:
		return false
	}
}

// BlocksBullets returns whether bullets are stopped by this tile.
func (t TileType) BlocksBullets() bool {
	switch t {
	case TileBrick, TileSteel, TileEagle, TileEagleDead:
		return true
	default:
		return false
	}
}
