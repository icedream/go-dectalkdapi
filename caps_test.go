package dectalkdapi_test

import (
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

func TestGetCaps(t *testing.T) {
	tts := startTTS(t)
	defer tts.Shutdown()

	caps, err := tts.GetCaps()
	if err != nil {
		t.Fatalf("GetCaps() failed: %v", err)
	}
	if caps.NumberOfLanguages == 0 {
		t.Error("NumberOfLanguages should be >= 1")
	}
	if caps.SampleRate == 0 {
		t.Error("SampleRate should be > 0")
	}
	if caps.MinimumSpeakingRate == 0 || caps.MaximumSpeakingRate == 0 {
		t.Error("Speaking rate range should be valid")
	}
	if caps.NumberOfPredefinedSpeakers == 0 {
		t.Error("NumberOfPredefinedSpeakers should be > 0")
	}
	if caps.Version == 0 {
		t.Error("Version should be > 0")
	}
	t.Logf("Capabilities: %+v", caps)
}

func TestGetFeatures(t *testing.T) {
	features, err := dectalkdapi.GetFeatures()
	if err != nil {
		t.Fatalf("GetFeatures() failed: %v", err)
	}
	if features == 0 {
		t.Error("Features should not be 0")
	}
	t.Logf("Features bitmask: %08x", features)
}

func TestGetStatus(t *testing.T) {
	tts := startTTS(t)
	defer tts.Shutdown()

	identifiers := []dectalkdapi.StatusIdentifier{dectalkdapi.InputCharacterCount}
	status, err := tts.GetStatus(identifiers, 1)
	if err != nil {
		t.Fatalf("GetStatus() failed: %v", err)
	}
	if len(status) != 1 {
		t.Fatalf("Expected 1 status value, got %d", len(status))
	}
	t.Logf("Status: %+v", status)
}

func TestVersionEx(t *testing.T) {
	ver, err := dectalkdapi.VersionEx()
	if err != nil {
		t.Fatalf("VersionEx() failed: %v", err)
	}
	if ver.StructSize == 0 {
		t.Error("StructSize should be > 0")
	}
	if ver.StructVersion == 0 {
		t.Error("StructVersion should be > 0")
	}
	if ver.DLLVersion == 0 || ver.DTalkVersion == 0 {
		t.Error("Version numbers should be > 0")
	}
	if ver.VerString == "" {
		t.Error("VerString should not be empty")
	}
	if ver.Language == "" {
		t.Error("Language should not be empty")
	}
	if ver.Features == 0 {
		t.Error("Features should not be 0")
	}
	t.Logf("Version: %+v", ver)
}

func TestEnumLangs(t *testing.T) {
	langs, err := dectalkdapi.EnumLangs()
	if err != nil {
		t.Fatalf("EnumLangs() failed: %v", err)
	}
	if langs.Languages == 0 {
		t.Error("Languages should be > 0")
	}
	if langs.Entries == nil {
		t.Error("Entries should not be nil")
	}
	for i, entry := range langs.Entries {
		if string(entry.LangCode[:]) == "" {
			t.Errorf("Entry %d: LangCode should not be empty", i)
		}
		if string(entry.LangName[:]) == "" {
			t.Errorf("Entry %d: LangName should not be empty", i)
		}
		t.Logf("Language %d: %s (%s)", i, entry.GetLangName(), entry.GetLangCode())
	}
}
