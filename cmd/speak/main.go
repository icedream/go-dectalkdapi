package main

import (
	"log"
	"os"

	"github.com/icedream/go-dectalkdapi"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func main() {
	tts, err := dectalkdapi.Startup(
		dectalkdapi.DoNotUseAudioDevice |
			dectalkdapi.ReportOpenError)
	must(err)
	defer tts.Shutdown()

	must(tts.OpenWaveOutFile("test.wav",
		dectalkdapi.WaveFormat1M16))
	defer tts.CloseWaveOutFile()

	defer tts.Sync()

	must(tts.Speak(
		`[:PHONE ON]`,
		dectalkdapi.Normal))

	must(tts.Speak(
		`[dah<600,20>][dah<600,20>][dah<600,20>]`+
			`[dah<500,16>][dah<130,23>][dah<600,20>]`+
			`[dah<500,16>][dah<130,23>][dah<600,20>]`,
		dectalkdapi.Force))

	must(tts.Speak(
		`Congratulations, you can now synthesize any text `+
			`the same way Moonbase Alpha did! aeiou`,
		dectalkdapi.Force))
}
