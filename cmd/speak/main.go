package main

import (
	"errors"
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
	defer must(tts.Shutdown())

	// TODO - renderInMemory(tts)
	_ = renderInMemory
	renderToWAV(tts)
}

func renderInMemory(tts *dectalkdapi.TTS) (err error) {
	if err = tts.OpenInMemory(dectalkdapi.WaveFormat1M16); err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, tts.CloseInMemory())
	}()

	buffer, err := tts.ReturnBuffer()
	must(err)

	_ = buffer // TODO

	err = speakExample(tts)
	return
}

func renderToWAV(tts *dectalkdapi.TTS) (err error) {
	if err = tts.OpenWaveOutFile("test.wav",
		dectalkdapi.WaveFormat1M16); err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, tts.CloseWaveOutFile())
	}()

	err = speakExample(tts)
	return
}

func speakExample(tts *dectalkdapi.TTS) (err error) {
	defer func() {
		err = errors.Join(err, tts.Sync())
	}()

	if err = tts.Speak(
		`[:PHONE ON]`,
		dectalkdapi.Normal); err != nil {
		return err
	}

	if err = tts.Speak(
		`[dah<600,20>][dah<600,20>][dah<600,20>]`+
			`[dah<500,16>][dah<130,23>][dah<600,20>]`+
			`[dah<500,16>][dah<130,23>][dah<600,20>]`,
		dectalkdapi.Force); err != nil {
		return err
	}

	if err = tts.Speak(
		`Congratulations, you can now synthesize any text `+
			`the same way Moonbase Alpha did! aeiou`,
		dectalkdapi.Force); err != nil {
		return err
	}

	return
}
