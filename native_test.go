package dectalkdapi_test

import (
	"testing"

	"github.com/icedream/go-dectalkdapi"
)

// Test memory buffer workflow
func TestMemoryBufferWorkflow(t *testing.T) {
	tts, err := dectalkdapi.Startup(dectalkdapi.DoNotUseAudioDevice | dectalkdapi.ReportOpenError)
	if err != nil {
		t.Fatalf("Startup() failed: %v", err)
	}
	defer func() {
		if err := tts.Shutdown(); err != nil {
			panic(err)
		}
	}()

	// Open memory buffer
	err = tts.OpenInMemory(dectalkdapi.WaveFormat1M16)
	if err != nil {
		t.Fatalf("OpenInMemory() failed: %v", err)
	}
	defer func() {
		if err := tts.Reset(true); err != nil {
			panic(err)
		}
	}()

	// Return buffer
	buffer, err := tts.ReturnBuffer()
	if err != nil {
		t.Fatalf("ReturnBuffer() failed: %v", err)
	}
	if buffer == nil {
		t.Fatal("ReturnBuffer() returned nil")
	}

	// Verify buffer can be used with AddBuffer
	// Note: AddBuffer requires properly initialized buffer, so we expect an error
	err = tts.AddBuffer(buffer)
	if err == nil {
		t.Fatal("AddBuffer() expected an error with uninitialized buffer")
	}

	// Close memory buffer
	err = tts.CloseInMemory()
	if err != nil {
		t.Fatalf("CloseInMemory() failed: %v", err)
	}
}
