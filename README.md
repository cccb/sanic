[![maintained](https://img.shields.io/maintenance/yes/2024)]()

# 🦔 sanic

chaos music control inspired by [relaxx player][relaxx]

## ✨ Features

- mpd web gui
  - search music
  - organize playlists
  - control current playback queue
- no authentication required to control music playback
- add playlists from internet radios (`*.m3u`, `*.pls`)
- add music from other sources like youtube (`youtube-dl`)

## 🛠️ Development

sanic is developed using [Nix][nix], but you can also just use the usual Golang tooling.

Run local [MPD][mpd] instance for testing with `make mpd`.

### ❄️ w/ Nix

Enter development shell (also has [mpc][mpc] client installed for testing):

```shell
nix develop
```

Build nix flake:

```shell
nix build
```

### 🐧 w/o Nix

Use these Make targets for convenience:

- `run`: Run project
- `build`: Compile project
- `tidy`: Add missing and remove unused modules
- `verify`: Verify dependencies have expected content
- `test`: Run tests
- `cert`: Create https certificate for local testing

### 🐳 Container

You can run sanic in a container. Use these Make targets for convenience:

- `build-container`: Build container image
- `run-container`: Run container image

## 🗺️ Architecture

[![Architecture](https://git.berlin.ccc.de/cccb/sanic/raw/branch/main/architecture.drawio.svg)](https://app.diagrams.net/?mode=git.berlin.ccc.de#Hcccb%2Fsanic%2Fmain%2Farchitecture.drawio.svg)

---

Made with ❤️ and ![golang logo][golang].

[relaxx]: http://relaxx.dirk-hoeschen.de/
[nix]: https://nixos.org/manual/nix/stable/
[golang]: https://go.dev/images/favicon-gopher.svg
[mpd]: https://musicpd.org/
[mpc]: https://www.musicpd.org/clients/mpc/
