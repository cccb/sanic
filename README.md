![Codecov](https://img.shields.io/codecov/c/github/cccb/sonic)
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/cccb/sonic/test)
[![license](https://img.shields.io/gitlab/license/xengi/dotfiles)](https://choosealicense.com/licenses/mit/)
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
- add playlists from internet radios (`*.m3u`)

## development

Build nix flake:

```shell
nix build
```

## architecture

[![Architecture](https://github.com/cccb/sanic/raw/main/architecture.drawio.svg)](https://app.diagrams.net/?mode=github#Hcccb%2Fsanic%2Fmain%2Farchitecture.drawio.svg)

---

Made with ❤️ and 🐍.

[relaxx]: http://relaxx.dirk-hoeschen.de/
