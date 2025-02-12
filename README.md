# play-publisher
![Coverage](https://img.shields.io/badge/Coverage-96.3%25-brightgreen)

[Installation](#installation) | [Docs](#docs) | [Notes](#notes)

`play-publisher` is a command line tool to upload an `APK`, or an `AAB` to the `Play Store`.

```
A Play Store uploader for your convenience

Usage:
  play-publisher [flags]
  play-publisher [command]

Available Commands:
  help          Help about any command
  info          Display the package name and type of the app
  upload        Upload an APK or AAB to the Play Store
  version       Print the version

Flags:
  -h, --help   help for play-publisher

Use "play-publisher [command] --help" for more information about a command.
```

## Demo:[^1]

![Demo](images/demo.gif "Demo")

## Why?

After years of using `fastlane`, a dedicated `GitHub action` or a `Gradle task` in various combinations to manage uploads, I wanted something simpler like `altool`, that manages `iOS` uploads to the `App Store`.

While searching for a command line tool to handle uploads to the `Play Store`, I've found that every solution I've come across requires a runtime. Be it `Node` or `Java`.

The reason for the existence of this repository is that I couldn't find a standalone solution without additional managed dependencies.

## Installation

### Go

```
go install github.com/anselstetter/play-publisher/cmd/play-publisher@latest
```

### Manually

Download the latest binaries from [here](https://github.com/anselstetter/play-publisher/releases) and copy them to your desired location.

### Other

There are no releases to other package repositories yet.

## Docs:

Info about commands are [here](./docs/play-publisher.md).

## Notes:

### macOS:

The releases for macOS are not signed, so macOS will deny running the binary.

To get around this issue you can delete the quarantine extended attribute with:

```
xattr -d com.apple.quarantine <binary>
```

## Windows:

This binary is cross compiled for Windows `AMD64` and `ARM64`, but has never been tested.

Your mileage may vary on these platforms.

[^1]: This is a real upload recorded with [asciinema](https://asciinema.org). Only the package name has been edited.
