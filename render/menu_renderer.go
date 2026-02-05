package render

import (
	"fmt"

	"github.com/AchrafSoltani/TankStrike/config"
	"github.com/AchrafSoltani/glow"
)

// MenuOption represents a selectable menu item.
type MenuOption struct {
	Label string
}

// DrawTitleScreen renders the main menu title screen.
func DrawTitleScreen(canvas *glow.Canvas, options []MenuOption, selected int, time float64) {
	cx := config.WindowWidth / 2

	// Background
	canvas.Clear(glow.Black)

	// Decorative border
	canvas.DrawRectOutline(20, 20, config.WindowWidth-40, config.WindowHeight-40, ColorDarkGray)
	canvas.DrawRectOutline(22, 22, config.WindowWidth-44, config.WindowHeight-44, ColorDarkGray)

	// Title: "TANK STRIKE" in large text
	DrawTextCentered(canvas, "TANK", cx, 100, ColorYellow, 6)
	DrawTextCentered(canvas, "STRIKE", cx, 160, ColorOrange, 5)

	// Tank art — small tank formation
	drawMenuTankArt(canvas, cx-80, 240)

	// Menu options
	optY := 380
	for i, opt := range options {
		color := ColorGray
		if i == selected {
			color = ColorYellow
			// Selection indicator (arrow)
			DrawText(canvas, ">", cx-120, optY, ColorWhite, 2)
		}
		DrawTextCentered(canvas, opt.Label, cx, optY, color, 2)
		optY += 30
	}

	// Flashing prompt
	if int(time*2)%2 == 0 {
		DrawTextCentered(canvas, "UP/DOWN TO SELECT, ENTER TO CONFIRM", cx, 520, ColorDarkGray, 1)
	}

	// Credits
	DrawTextCentered(canvas, "A BATTLE CITY REMAKE", cx, config.WindowHeight-60, ColorDarkGray, 1)
	DrawTextCentered(canvas, "MADE WITH GLOW ENGINE", cx, config.WindowHeight-44, ColorDarkGray, 1)
}

func drawMenuTankArt(canvas *glow.Canvas, x, y int) {
	// Player tank (large)
	canvas.DrawRect(x+60, y, 48, 48, ColorPlayerBody)
	canvas.DrawRect(x+78, y-16, 12, 20, ColorPlayerTread)
	canvas.FillCircle(x+84, y+24, 10, ColorPlayerDark)

	// Enemy tanks approaching
	colors := []glow.Color{ColorEnemyBasicBody, ColorEnemyFastBody, ColorEnemyPowerBody}
	for i, c := range colors {
		ex := x + 20 + i*60
		canvas.DrawRect(ex, y+70, 32, 32, c)
		canvas.DrawRect(ex+12, y+60, 8, 14, ColorDarkGray)
	}
}

// DrawPauseScreen renders the pause overlay.
func DrawPauseScreen(canvas *glow.Canvas, time float64) {
	// Dithered checkerboard
	for y := 0; y < config.WindowHeight; y += 2 {
		for x := 0; x < config.WindowWidth; x += 2 {
			canvas.SetPixel(x, y, ColorBlack)
		}
	}

	cx := config.WindowWidth / 2
	cy := config.WindowHeight / 2

	// Pause box
	canvas.DrawRect(cx-120, cy-40, 240, 80, ColorBlack)
	canvas.DrawRectOutline(cx-120, cy-40, 240, 80, ColorYellow)

	DrawTextCentered(canvas, "PAUSED", cx, cy-20, ColorYellow, 3)

	if int(time*2)%2 == 0 {
		DrawTextCentered(canvas, "ESC TO RESUME", cx, cy+20, ColorGray, 1)
	}
}

// DrawGameOverScreen renders the game over screen.
func DrawGameOverScreen(canvas *glow.Canvas, score int, canContinue bool, time float64) {
	// Dithered overlay
	for y := 0; y < config.WindowHeight; y += 2 {
		for x := (y / 2) % 2; x < config.WindowWidth; x += 2 {
			canvas.SetPixel(x, y, ColorBlack)
		}
	}

	cx := config.WindowWidth / 2
	cy := config.WindowHeight / 2

	// Game over box
	canvas.DrawRect(cx-140, cy-60, 280, 120, ColorBlack)
	canvas.DrawRectOutline(cx-140, cy-60, 280, 120, ColorRed)

	DrawTextCentered(canvas, "GAME OVER", cx, cy-40, ColorRed, 3)

	scoreText := fmt.Sprintf("FINAL SCORE: %d", score)
	DrawTextCentered(canvas, scoreText, cx, cy, ColorWhite, 2)

	if canContinue && int(time*2)%2 == 0 {
		DrawTextCentered(canvas, "PRESS ENTER", cx, cy+40, ColorGray, 1)
	}
}

// DrawLevelIntro renders the level introduction screen.
func DrawLevelIntro(canvas *glow.Canvas, level int) {
	canvas.Clear(glow.RGB(40, 40, 40))

	cx := config.WindowWidth / 2
	cy := config.WindowHeight / 2

	text := fmt.Sprintf("STAGE %d", level+1)
	DrawTextCentered(canvas, text, cx, cy-20, ColorWhite, 4)

	// Decorative lines
	lineW := TextWidth(text, 4)
	canvas.DrawRect(cx-lineW/2, cy+20, lineW, 2, ColorYellow)
}

// DrawLevelComplete renders the level complete tally.
func DrawLevelComplete(canvas *glow.Canvas, level int, score int,
	killsBasic, killsFast, killsPower, killsArmour int,
	canContinue bool, time float64) {
	cx := config.WindowWidth / 2

	// Dark overlay
	canvas.DrawRect(cx-200, 60, 400, 540, ColorBlack)
	canvas.DrawRectOutline(cx-200, 60, 400, 540, ColorYellow)

	y := 80
	stageText := fmt.Sprintf("STAGE %d CLEAR!", level+1)
	DrawTextCentered(canvas, stageText, cx, y, ColorPlayerBody, 3)
	y += 40

	scoreText := fmt.Sprintf("SCORE: %06d", score)
	DrawTextCentered(canvas, scoreText, cx, y, ColorYellow, 2)
	y += 50

	// Kill tally header
	DrawTextCentered(canvas, "- KILL TALLY -", cx, y, ColorWhite, 2)
	y += 30

	// Headers
	tallyX := cx - 150
	DrawText(canvas, "TYPE", tallyX, y, ColorGray, 1)
	DrawText(canvas, "KILLS", tallyX+80, y, ColorGray, 1)
	DrawText(canvas, "PTS", tallyX+140, y, ColorGray, 1)
	DrawText(canvas, "TOTAL", tallyX+220, y, ColorGray, 1)
	y += 4
	canvas.DrawRect(tallyX, y+10, 300, 1, ColorDarkGray)
	y += 18

	drawTallyLine(canvas, tallyX, y, "BASIC", killsBasic, 100, ColorEnemyBasicBody)
	y += 24
	drawTallyLine(canvas, tallyX, y, "FAST", killsFast, 200, ColorEnemyFastBody)
	y += 24
	drawTallyLine(canvas, tallyX, y, "POWER", killsPower, 300, ColorEnemyPowerBody)
	y += 24
	drawTallyLine(canvas, tallyX, y, "ARMOUR", killsArmour, 400, ColorEnemyArmourBody)
	y += 8
	canvas.DrawRect(tallyX, y, 300, 1, ColorDarkGray)
	y += 12

	total := killsBasic + killsFast + killsPower + killsArmour
	totalPts := killsBasic*100 + killsFast*200 + killsPower*300 + killsArmour*400
	DrawText(canvas, fmt.Sprintf("TOTAL: %d KILLS  %d PTS", total, totalPts), tallyX, y, ColorWhite, 1)
	y += 30

	if canContinue && int(time*2)%2 == 0 {
		DrawTextCentered(canvas, "PRESS ENTER TO CONTINUE", cx, y, ColorYellow, 1)
	}
}

func drawTallyLine(canvas *glow.Canvas, x, y int, name string, kills, ptsEach int, color glow.Color) {
	canvas.DrawRect(x, y+2, 10, 8, color)
	DrawText(canvas, name, x+16, y, ColorWhite, 1)
	DrawText(canvas, fmt.Sprintf("%2d", kills), x+88, y, ColorYellow, 1)
	DrawText(canvas, fmt.Sprintf("%3d", ptsEach), x+140, y, ColorYellow, 1)
	DrawText(canvas, fmt.Sprintf("%5d", kills*ptsEach), x+220, y, ColorYellow, 1)
}
