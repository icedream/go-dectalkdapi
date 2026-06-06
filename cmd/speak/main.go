package main

import (
	"bytes"
	"errors"
	"fmt"
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

func renderInMemory(tts *dectalkdapi.TTS) (err error) {
	log.Println("=== renderInMemory starting ===")
	if err = tts.OpenInMemory(dectalkdapi.WaveFormat1M16); err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, tts.CloseInMemory())
	}()

	// Note: Memory buffer mode in DECtalk requires manual buffer allocation
	// and proper initialization before adding to queue. The Go wrapper
	// ReturnBuffer() returns buffers that need to be initialized.
	// For now, let's skip the buffer queueing and try a simpler approach.

	// Speak a test phrase
	phrase := `Hello, world!`

	log.Printf("Speaking: %q\n", phrase)
	if err = tts.Speak(phrase, dectalkdapi.Normal); err != nil {
		return fmt.Errorf("Speak failed: %w", err)
	}
	log.Println("Speak completed")

	// Get the audio buffer
	log.Println("Getting buffer...")
	buffer, err := tts.ReturnBuffer()
	if err != nil {
		return fmt.Errorf("ReturnBuffer failed: %w", err)
	}
	log.Println("Got buffer")

	// Check buffer length
	log.Printf("Buffer length: %d", buffer.BufferLength())

	// Get the audio data
	log.Println("Getting audio data...")
	audioData, err := buffer.Data()
	if err != nil {
		return fmt.Errorf("Data() failed: %w", err)
	}
	log.Printf("Got %d bytes of audio", len(audioData))

	// Add the buffer back to the queue for reuse
	log.Println("Adding buffer back to queue...")
	if err = tts.AddBuffer(buffer); err != nil {
		return fmt.Errorf("Failed to add buffer back: %w", err)
	}
	log.Println("Buffer added back to queue")

	// Speak a test phrase
	phrase := `Hello, world!`

	log.Printf("Speaking: %q\n", phrase)
	if err = tts.Speak(phrase, dectalkdapi.Normal); err != nil {
		return fmt.Errorf("Speak failed: %w", err)
	}
	log.Println("Speak completed")

	// Get the audio buffer
	log.Println("Getting buffer...")
	buffer, err := tts.ReturnBuffer()
	if err != nil {
		return fmt.Errorf("ReturnBuffer failed: %w", err)
	}
	log.Println("Got buffer")

	// Check buffer length
	log.Printf("Buffer length: %d", buffer.BufferLength())

	// Get the audio data
	log.Println("Getting audio data...")
	audioData, err := buffer.Data()
	if err != nil {
		return fmt.Errorf("Data() failed: %w", err)
	}
	log.Printf("Got %d bytes of audio", len(audioData))

	// Add the buffer back to the queue for reuse
	log.Println("Adding buffer back to queue...")
	if err = tts.AddBuffer(buffer); err != nil {
		return fmt.Errorf("Failed to add buffer back: %w", err)
	}
	log.Println("Buffer added back to queue")

	// Write raw PCM data to output
	var output bytes.Buffer
	output.Write(audioData)

	log.Printf("Received %d bytes of audio", len(audioData))

	// Output the combined audio data to file for debugging
	log.Println("\n=== Writing output to file ===")
	outputFile, err := os.Create("output.pcm")
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	log.Printf("Writing %d bytes to output.pcm\n", output.Len())
	if _, err := outputFile.Write(output.Bytes()); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}
	log.Println("Successfully written to output.pcm")

	log.Println("\n=== Raw PCM Audio Output ===")
	log.Printf("Total audio data: %d bytes (%.2f seconds @ 8kHz mono 8-bit PCM)",
		output.Len(), float64(output.Len())/8000)
	log.Println("\nAudio saved to: output.pcm")
	log.Println("You can convert it using:")
	log.Println("  ffmpeg -f s16le -ar 8000 -ac 1 -i output.pcm output.wav")

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

func main() {
	log.Println("Starting speak program...")
	tts, err := dectalkdapi.Startup(
		dectalkdapi.DoNotUseAudioDevice |
			dectalkdapi.ReportOpenError)
	if err != nil {
		log.Fatal(err)
	}
	defer tts.Shutdown()

	log.Println("Running renderInMemory...")
	renderInMemory(tts)
	os.Exit(0)
}