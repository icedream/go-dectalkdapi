package dectalkdapi_test

import (
	"github.com/icedream/go-dectalkdapi"
	"testing"
)

func TestGetStatus(t *testing.T) {
	tts, err := dectalkdapi.Startup(dectalkdapi.DoNotUseAudioDevice | dectalkdapi.ReportOpenError)
	if err != nil {
		t.Fatalf("Startup() failed: %v", err)
	}
	defer tts.Shutdown()

	identifiers := []dectalkdapi.StatusIdentifier{dectalkdapi.StatusIdentifierSpeaking}
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
	if ver.VerString() == "" {
		t.Error("VerString should not be nil or empty")
	}
	if ver.Language() == "" {
		t.Error("Language should not be nil or empty")
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
		if entry.LangCode() == "" {
			t.Errorf("Entry %d: LangCode should not be empty", i)
		}
		if entry.LangName() == "" {
			t.Errorf("Entry %d: LangName should not be empty", i)
		}
		t.Logf("Language %d: %s (%s)", i, entry.LangName(), entry.LangCode())
	}
}

func TestGetStatusMultipleIdentifiers(t *testing.T) {
	tts, err := dectalkdapi.Startup(dectalkdapi.DoNotUseAudioDevice | dectalkdapi.ReportOpenError)
	if err != nil {
		t.Fatalf("Startup() failed: %v", err)
	}
	defer tts.Shutdown()

	identifiers := []dectalkdapi.StatusIdentifier{
		dectalkdapi.StatusIdentifierSpeaking,
		dectalkdapi.StatusIdentifierPaused,
		dectalkdapi.StatusIdentifierSilent,
		dectalkdapi.StatusIdentifierError,
	}
	status, err := tts.GetStatus(identifiers, uint32(len(identifiers)))
	if err != nil {
		t.Fatalf("GetStatus() failed: %v", err)
	}

	if len(status) != len(identifiers) {
		t.Fatalf("Expected %d status values, got %d", len(identifiers), len(status))
	}

	for i, s := range status {
		t.Logf("Status %d: Identifier=%v, Value=%d", i, s.Identifier, s.Value)
	}
}
