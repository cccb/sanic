[![maintained](https://img.shields.io/maintenance/yes/2023)]()

# sanic

chaos music control inspired by [relaxx player][relaxx]


## features

- mpd web gui
  - search music
  - organize playlists
  - control current playback queue
- no authentication required to control music playback
- add music from other sources like youtube (`youtube-dl`)
- add playlists from internet radios (`*.m3u`, `*.pls`)

## development

Build nix flake:

```shell
nix build
```

## architecture

[![Architecture](https://git.berlin.ccc.de/cccb/sanic/raw/branch/main/architecture.drawio.svg)](https://app.diagrams.net/?mode=git.berlin.ccc.de#Hcccb%2Fsanic%2Fmain%2Farchitecture.drawio.svg)

---

Made with ❤️ and 🐍.

[relaxx]: http://relaxx.dirk-hoeschen.de/
