{
  description = "sanic - chaos music control";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable-small";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:tweag/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };
  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      gomod2nix,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ gomod2nix.overlays.default ];
        };
        golang = pkgs.go_1_25;
        sanic = pkgs.buildGoApplication {
          pname = "sanic";
          version = "0.0.1";
          go = golang;
          src = ./.;
          modules = ./gomod2nix.toml;
        };
        go-test = pkgs.stdenvNoCC.mkDerivation {
          name = "go-test";
          dontBuild = true;
          src = ./.;
          doCheck = true;
          nativeBuildInputs = with pkgs; [
            golang
            writableTmpDirAsHomeHook
          ];
          checkPhase = ''
            go test -v ./...
          '';
          installPhase = ''
            mkdir "$out"
          '';
        };
        go-lint = pkgs.stdenvNoCC.mkDerivation {
          name = "go-lint";
          dontBuild = true;
          src = ./.;
          doCheck = true;
          nativeBuildInputs = with pkgs; [
            #golangci-lint
            golang
            go-tools
            writableTmpDirAsHomeHook
          ];
          checkPhase = ''
            #golangci-lint run
            staticcheck -f stylish ./...
          '';
          installPhase = ''
            mkdir "$out"
          '';
        };
      in
      {
        apps = {
          test = {
            type = "app";
            program = "${pkgs.writeShellScript "test" ''
              ${golang}/bin/go test -v ./...
            ''}";
          };
          lint = {
            type = "app";
            program = "${pkgs.writeShellScript "lint" ''
              ${pkgs.go-tools}/bin/staticcheck -f stylish ./...
            ''}";
          };
        };
        checks = { inherit go-test go-lint; };
        formatter = pkgs.nixfmt-tree;
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            golang
            go-tools # staticcheck
            gomod2nix.packages.${system}.default
          ];
          packages = with pkgs; [
            mpd
            mpc
            mkcert
          ];
        };
        packages = {
          default = sanic;
          container_old = pkgs.dockerTools.buildImage {
            name = "sanic";
            config = {
              Cmd = [ "${sanic}/bin/sanic" ];
            };
          };
          container = pkgs.dockerTools.buildLayeredImage {
            name = "sanic";
            tag = sanic.version;
            contents = [ sanic ];
            config.Cmd = [ "${sanic}/bin/sanic" ];
          };
        };
        nixosModules.default = import ./option.nix;
      }
    );
}
