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

// Test TTSBuffer.Data() method
func TestTTSBufferData(t *testing.T) {
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

	// Test Data() on nil buffer
	var nilBuffer *dectalkdapi.TTSBuffer
	_, err = nilBuffer.Data()
	if err == nil {
		t.Error("Data() on nil buffer should return error")
	}

	// Test Data() on empty buffer (before speech)
	// Note: buffer may not be properly initialized yet
	data, err := buffer.Data()
	// Empty buffer may return error or nil data
	if err != nil && err.Error() != "buffer has no data" && err.Error() != "buffer is nil" {
		t.Errorf("Data() returned unexpected error: %v", err)
	}
	if data != nil {
		t.Logf("Data() returned non-nil data for empty buffer: %d bytes", len(data))
	}

	// Close memory buffer
	err = tts.CloseInMemory()
	if err != nil {
		t.Fatalf("CloseInMemory() failed: %v", err)
	}
}

// Test TTSBuffer.BufferLength() method
func TestTTSBufferBufferLength(t *testing.T) {
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

	// Test BufferLength() on nil buffer
	var nilBuffer *dectalkdapi.TTSBuffer
	length := nilBuffer.BufferLength()
	if length != 0 {
		t.Errorf("BufferLength() on nil buffer should return 0, got %d", length)
	}

	// Test BufferLength() on empty buffer
	length = buffer.BufferLength()
	if length != 0 {
		t.Errorf("BufferLength() should return 0 for empty buffer, got %d", length)
	}

	// Close memory buffer
	err = tts.CloseInMemory()
	if err != nil {
		t.Fatalf("CloseInMemory() failed: %v", err)
	}
}

// Test TTSBuffer.PhonemeCount() method
func TestTTSBufferPhonemeCount(t *testing.T) {
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

	// Test PhonemeCount() on nil buffer
	var nilBuffer *dectalkdapi.TTSBuffer
	count := nilBuffer.PhonemeCount()
	if count != 0 {
		t.Errorf("PhonemeCount() on nil buffer should return 0, got %d", count)
	}

	// Test PhonemeCount() on empty buffer
	count = buffer.PhonemeCount()
	if count != 0 {
		t.Errorf("PhonemeCount() should return 0 for empty buffer, got %d", count)
	}

	// Close memory buffer
	err = tts.CloseInMemory()
	if err != nil {
		t.Fatalf("CloseInMemory() failed: %v", err)
	}
}

// Test TTSBuffer.IndexMarkCount() method
func TestTTSBufferIndexMarkCount(t *testing.T) {
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

	// Test IndexMarkCount() on nil buffer
	var nilBuffer *dectalkdapi.TTSBuffer
	count := nilBuffer.IndexMarkCount()
	if count != 0 {
		t.Errorf("IndexMarkCount() on nil buffer should return 0, got %d", count)
	}

	// Test IndexMarkCount() on empty buffer
	count = buffer.IndexMarkCount()
	if count != 0 {
		t.Errorf("IndexMarkCount() should return 0 for empty buffer, got %d", count)
	}

	// Close memory buffer
	err = tts.CloseInMemory()
	if err != nil {
		t.Fatalf("CloseInMemory() failed: %v", err)
	}
}

// Test TTSBuffer.MaximumBufferLength() method
func TestTTSBufferMaximumBufferLength(t *testing.T) {
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

	// Test MaximumBufferLength() on nil buffer
	var nilBuffer *dectalkdapi.TTSBuffer
	length := nilBuffer.MaximumBufferLength()
	if length != 0 {
		t.Errorf("MaximumBufferLength() on nil buffer should return 0, got %d", length)
	}

	// Test MaximumBufferLength() on buffer
	length = buffer.MaximumBufferLength()
	// Value may be 0 if buffer not yet initialized
	t.Logf("MaximumBufferLength(): %d", length)

	// Close memory buffer
	err = tts.CloseInMemory()
	if err != nil {
		t.Fatalf("CloseInMemory() failed: %v", err)
	}
}

// Test TTSBuffer.MaximumPhonemeChanges() method
func TestTTSBufferMaximumPhonemeChanges(t *testing.T) {
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

	// Test MaximumPhonemeChanges() on nil buffer
	var nilBuffer *dectalkdapi.TTSBuffer
	count := nilBuffer.MaximumPhonemeChanges()
	if count != 0 {
		t.Errorf("MaximumPhonemeChanges() on nil buffer should return 0, got %d", count)
	}

	// Test MaximumPhonemeChanges() on buffer
	count = buffer.MaximumPhonemeChanges()
	// Value may be 0 if buffer not yet initialized
	t.Logf("MaximumPhonemeChanges(): %d", count)

	// Close memory buffer
	err = tts.CloseInMemory()
	if err != nil {
		t.Fatalf("CloseInMemory() failed: %v", err)
	}
}

// Test TTSBuffer.MaximumIndexMarks() method
func TestTTSBufferMaximumIndexMarks(t *testing.T) {
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

	// Test MaximumIndexMarks() on nil buffer
	var nilBuffer *dectalkdapi.TTSBuffer
	count := nilBuffer.MaximumIndexMarks()
	if count != 0 {
		t.Errorf("MaximumIndexMarks() on nil buffer should return 0, got %d", count)
	}

	// Test MaximumIndexMarks() on buffer
	count = buffer.MaximumIndexMarks()
	// Value may be 0 if buffer not yet initialized
	t.Logf("MaximumIndexMarks(): %d", count)

	// Close memory buffer
	err = tts.CloseInMemory()
	if err != nil {
		t.Fatalf("CloseInMemory() failed: %v", err)
	}
}
