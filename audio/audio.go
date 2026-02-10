package audio

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/AchrafSoltani/glow"
)

// Engine manages audio playback using Glow's PulseAudio backend.
type Engine struct {
	ctx *glow.AudioContext

	// Pre-generated sound buffers
	shootBuf     []byte
	explodeBuf   []byte
	powerUpBuf   []byte
	gameOverBuf  []byte
	levelBuf     []byte
	menuSelBuf   []byte

	Muted  bool
	Volume float64
}

// NewEngine initialises the audio subsystem.
func NewEngine() *Engine {
	ctx, err := glow.NewAudioContext(sampleRate, 1, 2) // sampleRate, mono, 16-bit
	if err != nil {
		log.Printf("audio: failed to init audio: %v", err)
		return &Engine{}
	}

	e := &Engine{
		ctx:         ctx,
		shootBuf:    GenerateShoot(),
		explodeBuf:  GenerateExplode(),
		powerUpBuf:  GeneratePowerUp(),
		gameOverBuf: GenerateGameOver(),
		levelBuf:    GenerateLevelStart(),
		menuSelBuf:  GenerateMenuSelect(),
		Volume:      1.0,
	}
	return e
}

func (e *Engine) play(buf []byte) {
	if e.ctx == nil || len(buf) == 0 || e.Muted {
		return
	}

	scaled := buf
	if e.Volume < 1.0 {
		scaled = make([]byte, len(buf))
		for i := 0; i+1 < len(buf); i += 2 {
			sample := int16(binary.LittleEndian.Uint16(buf[i:]))
			sample = int16(float64(sample) * e.Volume)
			binary.LittleEndian.PutUint16(scaled[i:], uint16(sample))
		}
	}

	p := e.ctx.NewPlayer(bytes.NewReader(scaled))
	p.Play()
}

// ToggleMute toggles the muted state.
func (e *Engine) ToggleMute() {
	e.Muted = !e.Muted
}

// VolumeUp increases volume by 0.1, capped at 1.0.
func (e *Engine) VolumeUp() {
	e.Volume += 0.1
	if e.Volume > 1.0 {
		e.Volume = 1.0
	}
}

// VolumeDown decreases volume by 0.1, floored at 0.0.
func (e *Engine) VolumeDown() {
	e.Volume -= 0.1
	if e.Volume < 0.0 {
		e.Volume = 0.0
	}
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
