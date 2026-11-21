//go:build !linux || android || !cgo

package sound

// CheckPlaybackDevice skips the ALSA probe on other platforms.
func CheckPlaybackDevice() error { return nil }
