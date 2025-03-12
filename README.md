
<p align="center">
    <img src="./.github/assets/wordmark.png" alt="OpenTurntable" style="width:50%; height:auto;">
    <br/>
    A (hopefully good) OSS music player.
</p>

## Features
### Current
- Supports .mp3, .flac, .wav, and .ogg playback
- Gathers metadata from files (title, artist, album art, etc)
- Volume control

### Planned
- Library system to store a collection of music (see issue [#1](https://github.com/TheRandomMelon/OpenTurntable/issues/1))
- Playlist system (see issue [#2](https://github.com/TheRandomMelon/OpenTurntable/issues/2))
- Last.fm scrobbling support (see issue [#4](https://github.com/TheRandomMelon/OpenTurntable/issues/4))
- EQ (equalizer) (see issue [#5](https://github.com/TheRandomMelon/OpenTurntable/issues/5))

## Current Stack
- [Wails](https://wails.io)
- [Nuxt 3](https://nuxt.com)

## Using the Source Code

### Live Development
To run in live development mode, run `wails dev` in the project directory.

### Building
To build a redistributable, production mode package, use `wails build -tags "production"`.
