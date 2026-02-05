package audio

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/oto/v2"
)

// Engine manages audio playback using oto.
type Engine struct {
	ctx *oto.Context

	// Pre-generated sound buffers
	shootBuf     []byte
	explodeBuf   []byte
	powerUpBuf   []byte
	gameOverBuf  []byte
	levelBuf     []byte
	menuSelBuf   []byte
}

// NewEngine initialises the audio subsystem.
func NewEngine() *Engine {
	ctx, ready, err := oto.NewContext(sampleRate, 1, 2) // sampleRate, channelCount, bitDepthBytes
	if err != nil {
		log.Printf("audio: failed to init oto: %v", err)
		return &Engine{}
	}
	<-ready

	e := &Engine{
		ctx:         ctx,
		shootBuf:    GenerateShoot(),
		explodeBuf:  GenerateExplode(),
		powerUpBuf:  GeneratePowerUp(),
		gameOverBuf: GenerateGameOver(),
		levelBuf:    GenerateLevelStart(),
		menuSelBuf:  GenerateMenuSelect(),
	}
	return e
}

func (e *Engine) play(buf []byte) {
	if e.ctx == nil || len(buf) == 0 {
		return
	}
	p := e.ctx.NewPlayer(bytes.NewReader(buf))
	p.Play()
	// Player will be garbage collected after playback completes
}

// PlayShoot plays the shooting sound.
func (e *Engine) PlayShoot() { e.play(e.shootBuf) }

// PlayExplode plays the explosion sound.
func (e *Engine) PlayExplode() { e.play(e.explodeBuf) }

// PlayPowerUp plays the power-up collection sound.
func (e *Engine) PlayPowerUp() { e.play(e.powerUpBuf) }

// PlayGameOver plays the game over sound.
func (e *Engine) PlayGameOver() { e.play(e.gameOverBuf) }

// PlayLevelStart plays the level start fanfare.
func (e *Engine) PlayLevelStart() { e.play(e.levelBuf) }

// PlayMenuSelect plays the menu selection blip.
func (e *Engine) PlayMenuSelect() { e.play(e.menuSelBuf) }
