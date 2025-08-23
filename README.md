[![maintained](https://img.shields.io/maintenance/yes/2025?style=flat-square)]()
[![AUR Version](https://img.shields.io/aur/version/sanic?style=flat-square&logo=archlinux) ![AUR Last Modified](https://img.shields.io/aur/last-modified/sanic?style=flat-square&logo=archlinux)](https://aur.archlinux.org/packages/sanic)
![GitLab Release](https://img.shields.io/gitlab/v/release/xengi%2Fsanic?style=flat-square&logo=gitlab)


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

## 👩‍💻 Installation

### ❄️ NixOS (flakes)

Example flake setup (untested):

```nix
{
  description = "Example Flake to install sanic on your host";
  inputs = {
    nixpkgs.url = github:NixOS/nixpkgs/nixos-24.05;
    sanic = {
      url = gitlab:XenGi/sanic/main;
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };
  outputs = { self, nixpkgs, sanic }:
  let
    system = "x86_64-linux";
    pkgs = import nixpkgs { inherit system; };
  in
  {
    nixosConfigurations."myhostname".nixpkgs.lib.nixosSystem = {
      inherit system;
      modules = [
        { sanic.nixosModules.default }
        { services.sanic = {
            ui = {
              host = "[::1]";
              port = 8080;
              tls = false;
            };
            backend = {
              host = "localhost";
              port = 6600;
            };
          };
        }
      ];
    };
  };
}
```

### 🇦 Arch Linux

Install from the AUR:

```shell
yay -S sanic
```

### 🐳 Container

Run as daemon:

```shell
podman run -d -v ./config.ini:/config.ini -p 8080:8080 registry.gitlab.com/XenGi/sanic:latest
```

## 🛠️ Development

sanic is developed using [Nix][nix], but you can also just use the usual Golang tooling.

Run local [MPD][mpd] instance for testing with `make mpd`.

Update go dependencies like this:

```shell
go get -u  # or `make update`
go mod tidy  # or `make tidy`
gomod2nix  # sync go deps with nix
```

### ❄️ w/ Nix

Enter development shell (also has [mpc][mpc] client installed for testing):

```shell
nix develop
```

Build sanic:

```shell
nix build
```

### 🐧 w/o Nix

Use these Make targets for your convenience:

- `run`: Run project
- `build`: Compile project
- `tidy`: Add missing and remove unused modules
- `verify`: Verify dependencies have expected content
- `format`: Format go code
- `lint`: Run linter (staticcheck)
- `test`: Run tests
- `cert`: Create https certificate for local testing

### 🐳 Container

You can run sanic in a container. Use these Make targets for convenience:

- `build-container`: Build container image
- `run-container`: Run container image

## 🗺️ Architecture

[![Architecture](https://gitlab.com/XenGi/sanic/-/raw/main/architecture.drawio.svg)](https://app.diagrams.net/?mode=gitlab.com#AXenGi%2Fsanic%2Fmain%2Farchitecture.drawio.svg)

---

Made with ❤️ and ![golang logo][golang].

[relaxx]: http://relaxx.dirk-hoeschen.de/
[nix]: https://nixos.org/manual/nix/stable/
[golang]: https://go.dev/images/favicon-gopher.svg
[mpd]: https://musicpd.org/
[mpc]: https://www.musicpd.org/clients/mpc/

