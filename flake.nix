{
  description = "sanic - chaos music control";
  inputs = {
    nixpkgs.url = github:NixOS/nixpkgs/nixos-unstable-small;
    flake-utils.url = github:numtide/flake-utils;
    nixos-generators = {
      url = "github:nix-community/nixos-generators";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    gomod2nix = {
      url = github:tweag/gomod2nix;
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };
  outputs = { self, nixpkgs, flake-utils, nixos-generators, gomod2nix }: flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs {
        inherit system;
        overlays = [ gomod2nix.overlays.default ];
      };
      sanic = pkgs.buildGoApplication {
        pname = "sanic";
        version = "0.0.1";
        src = ./.;
        modules = ./gomod2nix.toml;
      };
      go-test = pkgs.stdenvNoCC.mkDerivation {
        name = "go-test";
        dontBuild = true;
        src = ./.;
        doCheck = true;
        nativeBuildInputs = with pkgs; [
          go
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
          golangci-lint
          go
          writableTmpDirAsHomeHook
        ];
        checkPhase = ''
          golangci-lint run
        '';
        installPhase = ''
          mkdir "$out"
        '';
      };
    in
    {
      checks = { inherit go-test go-lint; };
      formatter = pkgs.nixfmt-tree;
      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          go-tools # staticcheck
          gomod2nix.packages.${system}.default
        ];
        packages = with pkgs; [
          mpd
          mpc-cli
          mkcert
        ];
      };
      packages = {
        default = sanic;
        container = pkgs.dockerTools.buildImage {
          name = "sanic";
          config = {
            Cmd = [ "${sanic}/bin/sanic" ];
          };
        };
        proxmox-lxc = nixos-generators.nixosGenerate {
          system = "x86_64-linux";
          format = "proxmox-lxc";
          #specialArgs = {
          #  pkgs = pkgs;
          #};
          modules = [
            ./option.nix
            ./proxmox-lxc.nix
          ];
        };
      };
      nixosModules.default = import ./option.nix;
    }
  );
}
