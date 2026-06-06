package dectalkdapi_test

import (
	"os"
	"testing"

	dectalkdapi "github.com/icedream/go-dectalkdapi"
)

func startTTS(t *testing.T) *dectalkdapi.TTS {
	t.Helper()
	tts, err := dectalkdapi.Startup(dectalkdapi.DoNotUseAudioDevice | dectalkdapi.ReportOpenError)
	if err != nil {
		t.Fatalf("Startup() failed: %v", err)
	}
	return tts
}

func TestStartupAndShutdown(t *testing.T) {
	tts := startTTS(t)
	if err := tts.Shutdown(); err != nil {
		t.Errorf("Shutdown() failed: %v", err)
	}
}

func TestGetRate(t *testing.T) {
	tts := startTTS(t)
	defer tts.Shutdown()

	rate, err := tts.GetRate()
	if err != nil {
		t.Fatalf("GetRate() failed: %v", err)
	}
	if rate == 0 {
		t.Error("rate should not be 0")
	}
	t.Logf("Rate: %d", rate)
}

func TestSetRate(t *testing.T) {
	tts := startTTS(t)
	defer tts.Shutdown()

	newRate := uint32(200)
	if err := tts.SetRate(newRate); err != nil {
		t.Fatalf("SetRate(%d) failed: %v", newRate, err)
	}

	rate, err := tts.GetRate()
	if err != nil {
		t.Fatalf("GetRate() after SetRate failed: %v", err)
	}
	if rate != newRate {
		t.Errorf("expected rate %d, got %d", newRate, rate)
	}
}

func TestGetSpeaker(t *testing.T) {
	tts := startTTS(t)
	defer tts.Shutdown()

	speaker, err := tts.GetSpeaker()
	if err != nil {
		t.Fatalf("GetSpeaker() failed: %v", err)
	}
	t.Logf("Speaker: %d", speaker)
}

func TestSetSpeaker(t *testing.T) {
	tts := startTTS(t)
	defer tts.Shutdown()

	chosenSpeaker := dectalkdapi.Speaker(0)
	if err := tts.SetSpeaker(chosenSpeaker); err != nil {
		t.Fatalf("SetSpeaker() failed: %v", err)
	}

	speaker, err := tts.GetSpeaker()
	if err != nil {
		t.Fatalf("GetSpeaker() after SetSpeaker failed: %v", err)
	}
	if speaker != chosenSpeaker {
		t.Errorf("expected speaker %d, got %d", chosenSpeaker, speaker)
	}
}

func TestSpeakToWaveFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "dectalk-test-*.wav")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	tts := startTTS(t)
	defer tts.Shutdown()

	if err := tts.OpenWaveOutFile(tmpPath, dectalkdapi.WaveFormat1M16); err != nil {
		t.Fatalf("OpenWaveOutFile() failed: %v", err)
	}

	if err := tts.Speak("Hello, world!", dectalkdapi.Normal); err != nil {
		t.Fatalf("Speak() failed: %v", err)
	}

	if err := tts.Sync(); err != nil {
		t.Logf("Sync() after Speak: %v", err)
	}

	if err := tts.CloseWaveOutFile(); err != nil {
		t.Fatalf("CloseWaveOutFile() failed: %v", err)
	}

	info, err := os.Stat(tmpPath)
	if err != nil {
		t.Fatalf("stat temp file failed: %v", err)
	}
	if info.Size() == 0 {
		t.Error("wave file should not be empty after Speak")
	}
	t.Logf("Wave file size: %d bytes", info.Size())
}

func TestPauseResume(t *testing.T) {
	tts := startTTS(t)
	defer tts.Shutdown()

	if err := tts.Pause(); err != nil {
		t.Logf("Pause() expected maybe: %v", err)
	}

	if err := tts.Resume(); err != nil {
		t.Logf("Resume() expected maybe: %v", err)
	}
}

func TestReset(t *testing.T) {
	tts := startTTS(t)
	defer tts.Shutdown()

	if err := tts.Reset(true); err != nil {
		t.Errorf("Reset() failed: %v", err)
	}
}

func TestSync(t *testing.T) {
	tts := startTTS(t)
	defer tts.Shutdown()

	tmpFile, err := os.CreateTemp("", "dectalk-test-*.wav")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	if err := tts.OpenWaveOutFile(tmpPath, dectalkdapi.WaveFormat1M16); err != nil {
		t.Fatalf("OpenWaveOutFile() failed: %v", err)
	}

	tts.Speak("Hello, world!", dectalkdapi.Normal)
	if err := tts.Sync(); err != nil {
		t.Fatalf("Sync() failed: %v", err)
	}

	if err := tts.CloseWaveOutFile(); err != nil {
		t.Fatalf("CloseWaveOutFile() failed: %v", err)
	}
}

func TestVersion(t *testing.T) {
	verStr, dtMaj, dtMin, daMaj, daMin := dectalkdapi.Version()
	if verStr == "" {
		t.Error("version string should not be empty")
	}
	if daMin == 0 && daMaj == 0 {
		t.Error("API version should not be 0.0")
	}
	t.Logf("Version: %q (DECtalk %d.%d, DAPI %d.%d)", verStr, dtMaj, dtMin, daMaj, daMin)
}

func TestStartLang(t *testing.T) {
	lang, err := dectalkdapi.StartLang("US")
	if err != nil {
		t.Fatalf("StartLang() failed: %v", err)
	}
	defer lang.Close()

	name := lang.Name()
	if name == "" {
		t.Error("Name() should not be empty")
	}
	t.Logf("Language: %s", name)
}

func TestSelectLang(t *testing.T) {
	lang, err := dectalkdapi.StartLang("US")
	if err != nil {
		t.Fatalf("StartLang() failed: %v", err)
	}
	defer lang.Close()

	ok := dectalkdapi.SelectLang(lang)
	t.Logf("SelectLang() returned %v", ok)
}

func TestLogFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "dectalk-log-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	tts := startTTS(t)
	defer tts.Shutdown()

	if err := tts.OpenLogFile(tmpPath, dectalkdapi.Text|dectalkdapi.Phonemes); err != nil {
		t.Fatalf("OpenLogFile() failed: %v", err)
	}

	if err := tts.CloseLogFile(); err != nil {
		t.Fatalf("CloseLogFile() failed: %v", err)
	}

	info, err := os.Stat(tmpPath)
	if err != nil {
		t.Fatalf("stat log file failed: %v", err)
	}
	t.Logf("Log file size: %d bytes", info.Size())
}

func TestSpeakWithForce(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "dectalk-test-*.wav")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	tts := startTTS(t)
	defer tts.Shutdown()

	if err := tts.OpenWaveOutFile(tmpPath, dectalkdapi.WaveFormat1M16); err != nil {
		t.Fatalf("OpenWaveOutFile() failed: %v", err)
	}

	if err := tts.Speak("Hello, world!", dectalkdapi.Force); err != nil {
		t.Fatalf("Speak() with Force failed: %v", err)
	}

	if err := tts.Sync(); err != nil {
		t.Logf("Sync() after Speak with Force: %v", err)
	}

	if err := tts.CloseWaveOutFile(); err != nil {
		t.Fatalf("CloseWaveOutFile() failed: %v", err)
	}

	info, err := os.Stat(tmpPath)
	if err != nil {
		t.Fatalf("stat temp file failed: %v", err)
	}
	if info.Size() == 0 {
		t.Error("wave file should not be empty after Speak with Force")
	}
}
