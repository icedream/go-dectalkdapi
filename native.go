//go:build (windows && 386) || linux
// +build windows,386 linux

package dectalkdapi

/*
#cgo windows LDFLAGS: -ldectalk -lwinmm
#cgo linux LDFLAGS: -ltts

#if defined WIN32

#include <windows.h>
#include <TTSAPI.H>

const HWND ZeroHWND = NULL;

#else

#include <stdlib.h>
#include <dtk/ttsapi.h>

// define constants that cgo can't use
#define WAVERR_BADFORMAT 32

// define structs that don't exist on linux
typedef uint HWND;
const HWND ZeroHWND = 0;

#endif

// functions which have different signatures between linux and windows
#if defined __linux__ || defined VXWORKS || defined _SPARC_SOLARIS_ || defined __osf__
MMRESULT FixedTextToSpeechStartup(
	HWND hWnd,
	LPTTS_HANDLE_T * pphTTS,
	UINT uiDeviceNumber,
	DWORD dwDeviceOptions
) {
    // Linux code has a "New Audio Integration", hence the parameters are all over the place…
	return TextToSpeechStartup(
		pphTTS,
		uiDeviceNumber,
		dwDeviceOptions,
		NULL,
		hWnd // code mentioned "Backward compatibilty for TextToSpeechStartupEx" here
	);
}
#endif
#ifdef WIN32
MMRESULT FixedTextToSpeechStartup(HWND a, LPTTS_HANDLE_T *b, UINT c, DWORD d) {
	return TextToSpeechStartup(a, b, c, d);
}
#endif
*/
import "C"

import (
	"errors"
	"unsafe"
)

func parseIntAsBool(value C.BOOL) bool {
	return value == 0
}

func boolToInt(value bool) C.BOOL {
	if value {
		return 0
	}
	return 1
}

const waveMapper = C.WAVE_MAPPER

type Speaker C.SPEAKER_T

const (
	Paul   = C.PAUL
	Betty  = C.BETTY
	Harry  = C.HARRY
	Frank  = C.FRANK
	Dennis = C.DENNIS
	Kit    = C.KIT
	Ursula = C.URSULA
	Rita   = C.RITA
	Wendy  = C.WENDY
)

type Log C.DWORD

const (
	// Log text.
	Text Log = C.LOG_TEXT

	// Log phonemes.
	Phonemes = C.LOG_PHONEMES

	// Log syllables.
	Syllables = C.LOG_SYLLABLES
)

type MMResult C.DWORD

const (
	// No error occurred.
	NoError MMResult = C.MMSYSERR_NOERROR

	// No audio driver available.
	//
	// Will not occur with #DoNotUseAudioDevice.
	NoDriver = C.MMSYSERR_NODRIVER

	// DECtalk dictionary not found.
	Error = C.MMSYSERR_ERROR

	// No more DECtalk License units available.
	Allocated = C.MMSYSERR_ALLOCATED

	// No DECtalk License unit exists at all.
	NotEnabled = C.MMSYSERR_NOTENABLED

	// Memory allocation error.
	NoMem = C.MMSYSERR_NOMEM

	// Invalid handle.
	InvalidHandle = C.MMSYSERR_INVALHANDLE

	// Invalid parameter.
	InvalidParam = C.MMSYSERR_INVALPARAM

	// Device ID out of range.
	BadDeviceID = C.MMSYSERR_BADDEVICEID

	// Wave output device does not support request format.
	//
	// This error is only returned on Windows NT.
	BadFormat = C.WAVERR_BADFORMAT

	// Specified alias not found in WIN.INI.
	// InvalidAlias = C.MSYSERR_INVALIDALIAS

	// Invalid flag passed.
	// InvalidFlag = C.MSYSERR_INVALFLAG
)

func mmResultToError(codeC C.uint) error {
	code := MMResult(codeC)
	if code == NoError {
		return nil
	}

	return &MMError{
		Code: code,
	}
}

// TODO - make messages for Error more specific to each function

var mmErrorMessages = map[MMResult]string{
	NoDriver:      "no audio driver available",
	Error:         "general error occurred",
	Allocated:     "no more DECtalk license units available",
	NoMem:         "no memory available",
	InvalidHandle: "invalid handle",
	InvalidParam:  "invalid parameter",
}

var ErrCanNotLoadLanguage = errors.New("can not load language")

type TTSLangErrorCode C.DWORD

const (
	NotSupported TTSLangErrorCode = C.TTS_NOT_SUPPORTED
	NotAvailable                  = C.TTS_NOT_AVAILABLE
)

var ttsErrorMessages = map[TTSLangErrorCode]string{
	NotSupported: "not supported",
	NotAvailable: "not available",
}

type TTSLangError struct {
	Code TTSLangErrorCode
}

func (e *TTSLangError) Error() string {
	return ttsErrorMessages[e.Code]
}

func checkTTSLang(handle C.uint) error {
	if handle&C.TTS_LANG_ERROR != 0 {
		return &TTSLangError{
			Code: TTSLangErrorCode(handle),
		}
	}
	return nil
}

// WaveFormat represents an audio sample format.
type WaveFormat C.DWORD

const (
	// Mono, 8-bit 11.025 kHz sample rate
	WaveFormat1M08 = C.WAVE_FORMAT_1M08

	// Mono, 16-bit 11.025 kHz sample rate
	WaveFormat1M16 = C.WAVE_FORMAT_1M16

	// Mono, 8-bit μ-law, 8 kHz sample rate
	WaveFormat08M08 = C.WAVE_FORMAT_08M08
)

type TTSFlags C.DWORD

const (
	Normal TTSFlags = C.TTS_NORMAL
	Force           = C.TTS_FORCE
)

type DeviceOption C.DWORD

const (
	DoNotUseAudioDevice DeviceOption = C.DO_NOT_USE_AUDIO_DEVICE
	OwnAudioDevice      DeviceOption = C.OWN_AUDIO_DEVICE
	ReportOpenError     DeviceOption = C.REPORT_OPEN_ERROR
	UseSAPI5AudioDevice DeviceOption = C.USE_SAPI5_AUDIO_DEVICE
)

// TTSLanguage represents a language loaed into the DECtalk Multi-Language (ML)
// engine.
type TTSLanguage struct {
	name   *C.char
	handle C.uint
}

// Name returns the 2-character language ID.
func (l *TTSLanguage) Name() string {
	return C.GoString(l.name)
}

// Close closes an instance for an installed language and attempts to unload it
// from the DECtalk Multi-Language (ML) engine.
//
// Returns TRUE when a language is successfully unloaded, or FALSE when the
// operation cannot be completed or more instances have the thread started.
func (l *TTSLanguage) Close() bool {
	return parseIntAsBool(C.TextToSpeechCloseLang(l.name))
}

type TTS struct {
	handle C.LPTTS_HANDLE_T
}

// TODO - Add indicator to TTS for whether the engine is still active and check
// it on calls that would otherwise segfault

func Startup(deviceOptions DeviceOption) (*TTS, error) {
	tts := new(TTS)
	if err := tts.startup(C.ZeroHWND, deviceOptions); err != nil {
		return nil, err
	}
	return tts, nil
}

func (t *TTS) startup(windowHandle C.HWND, deviceOptions DeviceOption) error {
	return mmResultToError(C.FixedTextToSpeechStartup(windowHandle, &t.handle, waveMapper, C.DWORD(deviceOptions)))
}

// UnloadUserDictionary unloads a user dictionary. You must unload any
// previously loaded dictionary before you can load a new one. That is, only one
// user dictionary can be loaded at a time.
//
// A user dictionary is created using the User Dictionary Build tool.
func (t *TTS) UnloadUserDictionary() error {
	return mmResultToError(C.TextToSpeechUnloadUserDictionary(t.handle))
}

// Version requests version information from DECtalk Software that allows a
// calling application to test for DECtalk Software API (DAPI) compatibility.
// The function returns a numerically encoded version number and additionally
// may return a pointer to text information.
func Version() (
	versionStr string,
	dectalkMajor, dectalkMinor, dapiMajor, dapiMinor byte,
) {
	var versionStrC C.LPSTR
	versionNum := uint64(C.TextToSpeechVersion(&versionStrC))
	versionStr = C.GoString(versionStrC)

	dapiMajor = byte(versionNum & 0xff)
	dapiMinor = byte((versionNum >> 8) & 0xff)
	dectalkMinor = byte((versionNum >> 16) & 0xff)
	dectalkMajor = byte((versionNum >> 24) & 0xff)
	return
}

// LoadUserDictionary loads a user-defined pronunciation dictionary into the
// text-to-speech system.
//
// This function loads a dictionary created by the windict or userdict applet
// (Linux or UNIX) or the windic applet (Windows). Any previously loaded user
// dictionary must be unloaded before loading a new user dictionary. Note that
// the text-to-speech system will automatically load a user dictionary, user.dic
// (or udict_langcode.dic for Linux), at startup if it exists in the home
// directory.
func (t *TTS) LoadUserDictionary(dictFile string) error {
	dictFileC := C.CString(dictFile)
	defer C.free(unsafe.Pointer(dictFileC))
	return mmResultToError(C.TextToSpeechLoadUserDictionary(t.handle, dictFileC))
}

// OpenWaveOutFile opens the specified wave file and causes the text-to-speech
// system to enter into wave-file mode. This mode indicates that the speech
// samples are to be written in wave format into the wave file each time #Speak
// is called. The text-to-speech system remains in the wave-file mode until
// #CloseWaveOutFile is called.
//
// This function automatically resumes audio output if the text-to-speech system
// is in a paused state by a previously issued #Pause call.
func (t *TTS) OpenWaveOutFile(outFile string, format WaveFormat) error {
	outFileC := C.CString(outFile)
	defer C.free(unsafe.Pointer(outFileC))
	return mmResultToError(C.TextToSpeechOpenWaveOutFile(t.handle, outFileC, C.DWORD(format)))
}

// CloseWaveOutFile closes a wave file opened by the #OpenWaveOutFile function
// and returns to the startup state. The speech samples are then ignored or sent
// to an audio device, depending on the setting of the deviceOptions parameter
// in the startup function.
//
// The application must have called #OpenWaveOutFile before calling
// CloseWaveOutFile.
func (t *TTS) CloseWaveOutFile() error {
	return mmResultToError(C.TextToSpeechCloseWaveOutFile(t.handle))
}

// OpenLogFile opens the specified log file and causes the text-to-speech system
// to enter into the log-file mode. This mode indicates that the speech samples
// are to be written as text, phonemes, or syllables into the log file each time
// #Speak is called. The phonemes and syllables are written using the arpabet
// alphabet. The text-to-speech system remains in the log-file mode until
// #CloseLogFile is called.
//
// If more than one of the dwFlags are passed, the logged output is mixed in an
// unpredictable fashion.
//
// If a log file is open already, this function returns an error. The Log
// voice-control command also has no effect when a log file is open already.
//
// OpenLogFile automatically resumes audio output if the text-to-speech system
// is in a paused state by a previously issued #Pause call.
func (t *TTS) OpenLogFile(outFile string, log Log) error {
	outFileC := C.CString(outFile)
	defer C.free(unsafe.Pointer(outFileC))
	return mmResultToError(C.TextToSpeechOpenLogFile(t.handle, outFileC, C.DWORD(log)))
}

// CloseLogFile closes a log file opened by #OpenLogFile and returns to the
// startup state. The speech samples are then ignored or sent to an audio
// device, depending on the setting of the deviceOptions parameter in the
// startup function.
//
// CloseLogFile closes any open log file, even if it was opened with the Log
// command.
//
// The application must have called #OpenLogFile before calling CloseLogFile.
func (t *TTS) CloseLogFile() error {
	return mmResultToError(C.TextToSpeechCloseLogFile(t.handle))
}

// Speak queues a null-terminated string to the text-to-speech system.
//
// While the text-to-speech system is in the startup state, speech samples are
// routed to the audio device or ignored, depending on whether the startup
// function flag DO_NOT_USE_AUDIO_DEVICE is clear or set in the dwDeviceOptions
// parameter of the startup function.
//
// If the text-to-speech system is in a special mode (wave-file, log-file, or
// speech-to-memory modes), the speech samples are handled as the mode dictates.
//
// The speaker, speaking rate, and volume also can be changed in the text string
// by inserting voice- control commands, as shown in the following example:
//
// [:name paul] I am Paul. [:nb] I am Betty. [:volume set 50] The volume has
// been set to 50% of the maximum level. [:rate 120] I am speaking at 120 words
// per minute.
func (t *TTS) Speak(text string, flags TTSFlags) error {
	textC := C.CString(text)
	defer C.free(unsafe.Pointer(textC))
	return mmResultToError(C.TextToSpeechSpeak(t.handle, textC, C.DWORD(flags)))
}

// StartLang checks whether the specified language is installed and, if so,
// loads the language into the DECtalk ML engine.
//
// StartLang must be called before a language can be selected and opened in a
// multi-language application.
func StartLang(lang string) (*TTSLanguage, error) {
	langC := C.CString(lang)
	defer C.free(unsafe.Pointer(langC))

	ttsLang := C.TextToSpeechStartLang(langC)
	if err := checkTTSLang(ttsLang); err != nil {
		return nil, err
	}

	return &TTSLanguage{
		name:   langC,
		handle: ttsLang,
	}, nil
}

// SelectLang selects a loaded language for a program thread.
func SelectLang(lang *TTSLanguage) (ok bool) {
	okC := C.TextToSpeechSelectLang(nil, lang.handle)
	return parseIntAsBool(okC)
}

// // Do not use - first parameter is reserved.
// func (t *TTS) SelectLang(lang *TTSLanguage) (ok bool) {
// 	okC := C.TextToSpeechSelectLang(t.handle, lang.handle)
// 	return parseIntAsBool(okC)
// }

// Sync blocks until all previously queued text is processed.
//
// This function automatically resumes audio output if the text-to-speech system
// is in a paused state by a previously issued #Pause call.
func (t *TTS) Sync() error {
	return mmResultToError(C.TextToSpeechSync(t.handle))
}

// Typing speaks a single letter as quickly as possible, aborting any previously
// queued speech. This is somewhat slower if #Speak has been called since the
// last Typing or #Reset call.
//
// This function is primarily useful with the Access32 versions of DECtalk
// Software. The function exists in non-Access32 versions, but is not fast.
//
// This function should be called only when the application is synthesizing
// directly to an audio device (not to memory or to a file).
func (t *TTS) Typing(letter rune) {
	C.TextToSpeechTyping(t.handle, C.uchar(letter))
}

// Shutdown shuts down the text-to-speech system and frees all its system
// resources.
//
// Shutdown is called to close an application. Any user-defined dictionaries
// that were previously loaded are unloaded. All previously queued text is
// discarded, and the text-to-speech system immediately stops speaking.
func (t *TTS) Shutdown() error {
	// TODO - check where a left-over wave file is open, otherwise the API will hang!

	return mmResultToError(C.TextToSpeechShutdown(t.handle))
}

// Pause pauses text-to-speech audio output.
//
// This function affects only the audio output and has no effect when writing
// log files or wave files, or when using the speech-to-memory capability of the
// text-to-speech system.
//
// If the text-to-speech system owns the audio device (that is, #OwnAudioDevice
// was specified in the startup function), then the text-to-speech system
// remains paused until #Resume, #Sync, #OpenInMemory, #OpenLogFile, or
// #OpenWaveOutFile is called.
//
// If the text-to-speech system does not own the audio device (#OwnAudioDevice
// was NOT specified in the startup function) and #Pause is called while the
// system is speaking, the text-to-speech system remains paused until the system
// has completed speaking.
//
// In this case, the wave output device is released when #Reset is called. It
// will also be released if #Sync, #OpenInMemory, #OpenLogFile, or
// #OpenWaveOutFile is called AND the system has completed speaking.
//
// Note that #Pause will NOT resume audio output if the text-to-speech system is
// paused by #Pause.
func (t *TTS) Pause() error {
	return mmResultToError(C.TextToSpeechPause(t.handle))
}

// Resume resumes text-to-speech output after it was paused by calling
// TextToSpeechPause.
//
// This function affects only audio output and has no effect when writing log
// files or wave files or when writing speech samples to memory.
func (t *TTS) Resume() error {
	return mmResultToError(C.TextToSpeechResume(t.handle))
}

// Reset flushes all previously queued text from the text-to-speech system and
// stops any audio output.
//
// If the #OpenInMemory function has enabled writing speech samples to memory,
// all queued memory buffers are returned to the calling application.
//
// If the fullReset flag is on and the text-to-speech system is in one of its
// special modes (log-file, wave-file, or speech-to-memory mode), all files are
// closed and the text-to-speech system is returned to the startup state.
//
// #Reset should be called before calling #CloseInMemory. Failing to do this in
// a situation where the synthesizer is busy may result in a deadlock.
func (t *TTS) Reset(fullReset bool) error {
	// TODO - implement deadlock checks (CloseInMemory called before Reset?)

	return mmResultToError(C.TextToSpeechReset(t.handle, boolToInt(fullReset)))
}

// GetRate returns the current setting of the speaking rate.
//
// Valid values range from 75 to 600 words per minute.
//
// The current setting of the speaking rate is returned even if the speaking
// rate change has not yet occurred. This may occur when the #SetRate function
// is used without the #Sync function. The speaking-rate change occurs on clause
// boundaries.
func (t *TTS) GetRate() (uint32, error) {
	var rateC C.DWORD
	if err := mmResultToError(C.TextToSpeechGetRate(t.handle, &rateC)); err != nil {
		return 0, err
	}
	return uint32(rateC), nil
}

// SetRate sets the text-to-speech speaking rate.
//
// The speaking rate change is not effective until the next phrase boundary. All
// the queued audio encountered before the phrase boundary is unaffected.
func (t *TTS) SetRate(rate uint32) error {
	return mmResultToError(C.TextToSpeechSetRate(t.handle, C.DWORD(rate)))
}

// GetSpeaker returns the value of the identifier for the last voice that has
// spoken.
//
// Note that even after calling #SetSpeaker(), #GetSpeaker() returns the value
// for the previous speaking voice until the new voice actually speaks.
func (t *TTS) GetSpeaker() (Speaker, error) {
	var speakerC C.SPEAKER_T
	err := mmResultToError(C.TextToSpeechGetSpeaker(t.handle, &speakerC))
	if err != nil {
		return 0, err
	}
	return Speaker(speakerC), nil
}

// SetSpeaker sets the voice of the speaker that the text-to-speech system is to
// use.
//
// The change in speaking voice is not effective until the next phrase boundary.
// All queued audio encountered before the phrase boundary is unaffected.
func (t *TTS) SetSpeaker(speaker Speaker) error {
	return mmResultToError(C.TextToSpeechSetSpeaker(t.handle, C.SPEAKER_T(speaker)))
}

// TODO - MMRESULT #AddBuffer(LPTTS_HANDLE_T phTTS, LPTTS_BUFFER_T pTTSbuffer) Adds a shared-memory buffer allocated by the calling application to the memory buffer list.
// TODO - MMRESULT #CloseInMemory(LPTTS_HANDLE_T phTTS) Returns the text-to-speech system to its startup state.
// TODO - DWORD #EnumLangs(LPLANG_ENUM *langs) retrieves information about what languages are available in the system.
// TODO - MMRESULT #GetCaps(LPTTS_CAPS_T lpTTScaps) Retrieves the capabilities of the text-to-speech system
// TODO - DWORD #GetFeatures(void) Retrieves information, in the form of a bitmask, about the features of DECtalk Software. (maskable to the list supplied in the header file TTSFEAT.H.)
// TODO - MMRESULT #GetRate(LPTTS_HANDLE_T phTTS, LPDWORD pdwRate) Returns the speaking rate of the text-to-speech system.
// TODO - MMRESULT #GetStatus(LPTTS_HANDLE_T phTTS, LPDWORD dwIdentifier[ ], LPDWORD dwStatus[ ], DWORD dwNumberOfStatusValues) Gets the status of the text-to-speech system
// TODO - MMRESULT #OpenInMemory(LPTTS_HANDLE_T phTTS, DWORD dwFormat) <requires TextToSpeechAddBuffer> Produces buffered speech samples in wave format whenever #Speak function is called. The calling application is notified when memory buffer is filled.
// TODO - MMRESULT #ReturnBuffer(LPTTS_HANDLE_T phTTS, LPTTS_BUFFER_T *ppTTSbuffer) Returns the current shared-memory buffer.
// TODO - MMRESULT #StartupEx(LPTTS_HANDLE_T *phTTS, UINT uiDeviceNumber, DWORD dwDeviceOptions, VOID (*DtCallbackRoutine)(), LONG dwCallbackParameter) TextToSpeechStartup but with custom callback
// TODO - ULONG TextToSpeechVersionEx(LPVERSION_INFO *ver)
// TODO - struct LPVERSION_INFO
