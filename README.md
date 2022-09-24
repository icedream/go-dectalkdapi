# DECtalk DAPI library for Golang

You remember the goofy text-to-speech that Moonbase Alpha had, yeah? And you
also like Golang, yeah? Well, this library is definitely right up your alley
then!

This library opens the gate for using DECtalk's DAPI to generate speech from any
given text. This means you can now fully automate this text-to-speech engine
straight from Go code.

You can fully use DECtalk's TTS engine including inline commands such as
modifying speed, pitch, inserting DTMF sounds, phonemes etc. through this
library!

## Requirements

In order to use this library you need a fully working copy of the DECtalk DAPI
SDK.

This library is able to target the DECtalk 4.61 or 5.0 DAPI on Windows and Linux
through cross-compiling. If you do not have a copy of DECtalk's files for DAPI,
a modernized copy of the DECtalk 5.0 beta source tree is available at
https://github.com/dectalk/dectalk which you can compile yourself.

If you want to cross-compile for Windows 32-bit (as some copies of the SDK only
ship with binaries for that) you can rely on [MinGW-w64](https://www.mingw-w64.org/):

```bash
export CGO_ENABLED=1 GOARCH=386 CC=i686-w64-mingw32-gcc
```

## Features

### Implemented features

- Support for Windows and Linux
- Full support for inline commands as-is
- Audio output to sound device
- Audio output to WAV file
- Manipulation of speech rate through API call
- Log output for text, phonemes, syllables
- Fast-pace single-letter speech output
- Speaker switching through API call
- Multi-language support
- Wrapping of native error codes to Go error objects
- Simple version querying

### Currently missing features

- Audio output to memory buffer
- Callback functionality
- Manipulation of speaker through API call
- Additional Go-side checks for bad code conditions such as those known to lead
  to deadlocks
- Engine status querying
- Features querying
- Engine capabilities querying
- Version querying to a struct

## Building

You need to provide Go with the correct paths for this library to find the
needed binary files. On Linux, you can do this before building with Go:

```bash
CGO_LDFLAGS="-L${PATH_TO_DECTALK_INSTALL}/Us" # DECtalk 5.0 Windows SDK
CGO_LDFLAGS="-L${PATH_TO_DECTALK_INSTALL}/lib" # DECtalk 4.61/5.0 Linux SDK
CGO_CFLAGS="-I${PATH_TO_DECTALK_INSTALL}/include"
export CGO_ENABLED=1 CGO_LDFLAGS CGO_CFLAGS
```

You need to replace or set `${PATH_TO_DECTALK_INSTALL}` to the path of your own copy of DECtalk SDK.

With the above information, you can build the example code provided in this
repository. It will simply render a popular speech-synthesized interpretation of
The Imperial March's first few notes and a congratulatory message to a file
called `test.wav`. One way of running it could be through [WINE](https://winehq.org):

```bash
# compile with Go (make sure you set up the environment as described above)
go build -v ./cmd/speak

# run with 32-bit wine
WINEARCH=win32 WINEPREFIX="$(pwd)/wineprefix" wine ./speak.exe

# listen to the result
aplay test.wav
```

## License

First of all, I am not providing any ready-built binaries or even test results
for this project, neither of my own code nor of DECtalk, mainly due to legal
constraints: DECtalk is abandonware but still not 100% open for any open-source
project to simply include or even link against.

What I am doing is to provide my own code, available for anyone to inspect the
work that I have done around DECtalk's API specification. That sole work I can
safely publish and license.

That said, this project's code is licensed under the [MIT license](LICENSE.txt).
