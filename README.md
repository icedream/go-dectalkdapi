# DECtalk DAPI library for Golang

You remember the goofy text-to-speech that Moonbase Alpha had, yeah? And you
also like Golang, yeah? Well, this library is definitely right up your alley
then!

This library opens the gate for using DECtalk's DAPI to generate speech from any
given text. This means you can now fully automate this text-to-speech engine
straight from Go code.

## Requirements

In order to use this library you need a fully working copy of the DECtalk DAPI
SDK. If you search around on the internet for long enough, you may find one, but
I'm not going to provide such a copy here for copyright/licensing reasons.

Also, as of right now this library only really supports building against the
32-bit Windows binary files, so the only way to build anything with this library
currently is with `CGO_ENABLED=1 GOOS=windows GOARCH=386`. Sorry.

If you want to cross-compile you can rely on mingw64
(`CC=i686-w64-mingw32-gcc`).

## Building

You need to provide Go with the correct paths for this library to find the
needed binary files. On Linux, you can do this before building with Go:

```bash
CGO_LDFLAGS="-L${PATH_TO_DECTALK_INSTALL}/Us/"
CGO_CFLAGS="-I${PATH_TO_DECTALK_INSTALL}/include"
export CGO_ENABLED=1 GOOS=windows GOARCH=386 \
    CGO_LDFLAGS CGO_CFLAGS
```

You need to replace or set `${PATH_TO_DECTALK_INSTALL}` to the path of your own copy of DECtalk SDK.

If you want to use a different default language, you can replace `Us/` respectively as well.

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
