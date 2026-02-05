package world

import "github.com/AchrafSoltani/TankStrike/config"

// LoadLevel parses a level string and populates the grid.
// Level data uses characters:
//
//	'.' = empty, 'B' = brick, 'S' = steel, 'W' = water,
//	'I' = ice, 'F' = forest, 'E' = eagle
//
// The string should have 26 lines of 26 characters each.
func LoadLevel(g *Grid, data string) {
	g.Clear()
	x, y := 0, 0
	for _, ch := range data {
		if ch == '\n' {
			y++
			x = 0
			continue
		}
		if x < config.GridWidth && y < config.GridHeight {
			switch ch {
			case '.':
				g.Tiles[y][x] = TileEmpty
			case 'B':
				g.Tiles[y][x] = TileBrick
			case 'S':
				g.Tiles[y][x] = TileSteel
			case 'W':
				g.Tiles[y][x] = TileWater
			case 'I':
				g.Tiles[y][x] = TileIce
			case 'F':
				g.Tiles[y][x] = TileForest
			case 'E':
				g.Tiles[y][x] = TileEagle
			}
		}
		x++
	}
}
