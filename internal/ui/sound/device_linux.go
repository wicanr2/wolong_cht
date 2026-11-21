//go:build linux && !android && cgo

package sound

/*
#cgo pkg-config: alsa
#include <alsa/asoundlib.h>
#include <stdlib.h>
#include <string.h>

static int can_open_pcm(const char *name) {
	snd_pcm_t *pcm;
	if (snd_pcm_open(&pcm, name, SND_PCM_STREAM_PLAYBACK, SND_PCM_NONBLOCK) < 0) {
		return 0;
	}
	snd_pcm_close(pcm);
	return 1;
}

static int has_playback_device(void) {
	if (can_open_pcm("default") || can_open_pcm("plug:default")) {
		return 1;
	}

	void **hints;
	if (snd_device_name_hint(-1, "pcm", &hints) < 0) {
		return 0;
	}
	int available = 0;
	for (void **hint = hints; *hint && !available; ++hint) {
		char *direction = snd_device_name_get_hint(*hint, "IOID");
		char *name = snd_device_name_get_hint(*hint, "NAME");
		if ((!direction || strcmp(direction, "Input") != 0) &&
		    name && strcmp(name, "null") != 0 && strcmp(name, "default") != 0) {
			available = can_open_pcm(name);
		}
		free(direction);
		free(name);
	}
	snd_device_name_free_hint(hints);
	return available;
}
*/
import "C"

import "errors"

// CheckPlaybackDevice checks whether at least one ALSA playback PCM can open.
// Oto also tries devices listed by ALSA, including outputs without /dev/snd.
func CheckPlaybackDevice() error {
	if C.has_playback_device() == 0 {
		return errors.New("找不到可開啟的 ALSA 播放裝置")
	}
	return nil
}
