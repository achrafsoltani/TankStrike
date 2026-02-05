package audio

import "math"

const sampleRate = 44100

// GenerateShoot creates a short high-pitched burst.
func GenerateShoot() []byte {
	duration := 0.1 // seconds
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*2) // 16-bit mono

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		progress := float64(i) / float64(samples)

		// Descending frequency sweep
		freq := 800.0 - 400.0*progress
		val := math.Sin(2 * math.Pi * freq * t)

		// Envelope: quick attack, fast decay
		env := 1.0 - progress
		env *= env

		sample := int16(val * env * 8000)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	return buf
}

// GenerateExplode creates a noise-based explosion sound.
func GenerateExplode() []byte {
	duration := 0.4
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*2)

	// Simple LFSR noise
	lfsr := uint16(0xACE1)

	for i := 0; i < samples; i++ {
		progress := float64(i) / float64(samples)

		// Noise via LFSR
		bit := ((lfsr >> 0) ^ (lfsr >> 2) ^ (lfsr >> 3) ^ (lfsr >> 5)) & 1
		lfsr = (lfsr >> 1) | (bit << 15)
		noise := float64(int16(lfsr)) / 32768.0

		// Low-frequency rumble
		t := float64(i) / float64(sampleRate)
		rumble := math.Sin(2*math.Pi*60*t) * 0.5

		// Mix
		val := (noise*0.6 + rumble*0.4)

		// Envelope
		env := 1.0 - progress
		env = env * env * env

		sample := int16(val * env * 10000)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	return buf
}

// GeneratePowerUp creates a pleasant ascending arpeggio.
func GeneratePowerUp() []byte {
	duration := 0.5
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*2)

	notes := []float64{523.25, 659.25, 783.99, 1046.50} // C5, E5, G5, C6
	noteLen := samples / len(notes)

	for i := 0; i < samples; i++ {
		noteIdx := i / noteLen
		if noteIdx >= len(notes) {
			noteIdx = len(notes) - 1
		}
		freq := notes[noteIdx]
		t := float64(i) / float64(sampleRate)
		progress := float64(i) / float64(samples)

		val := math.Sin(2*math.Pi*freq*t) * 0.7
		val += math.Sin(2*math.Pi*freq*2*t) * 0.2 // harmonic

		env := 1.0 - progress*0.5
		sample := int16(val * env * 6000)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	return buf
}

// GenerateGameOver creates a descending sad tone.
func GenerateGameOver() []byte {
	duration := 1.0
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*2)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		progress := float64(i) / float64(samples)

		freq := 400.0 - 200.0*progress
		val := math.Sin(2*math.Pi*freq*t)*0.6 + math.Sin(2*math.Pi*freq*0.5*t)*0.3

		env := 1.0 - progress
		sample := int16(val * env * 8000)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	return buf
}

// GenerateLevelStart creates a brief fanfare.
func GenerateLevelStart() []byte {
	duration := 0.6
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*2)

	notes := []float64{392.0, 523.25, 659.25} // G4, C5, E5
	noteLen := samples / len(notes)

	for i := 0; i < samples; i++ {
		noteIdx := i / noteLen
		if noteIdx >= len(notes) {
			noteIdx = len(notes) - 1
		}
		freq := notes[noteIdx]
		t := float64(i) / float64(sampleRate)
		localT := float64(i%noteLen) / float64(noteLen)

		val := math.Sin(2 * math.Pi * freq * t)

		// Per-note envelope
		env := 1.0
		if localT > 0.8 {
			env = (1.0 - localT) * 5.0
		}

		sample := int16(val * env * 6000)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	return buf
}

// GenerateMenuSelect creates a short click/blip.
func GenerateMenuSelect() []byte {
	duration := 0.05
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*2)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		progress := float64(i) / float64(samples)

		val := math.Sin(2 * math.Pi * 1000 * t)
		env := 1.0 - progress
		sample := int16(val * env * 5000)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	return buf
}
