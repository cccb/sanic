{
  description = "sanic - chaos music control";
  inputs = {
    nixpkgs.url = github:NixOS/nixpkgs/nixos-24.05;
    flake-utils.url = github:numtide/flake-utils;
    gomod2nix = {
      url = github:tweag/gomod2nix;
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };
  outputs = { self, nixpkgs, flake-utils, gomod2nix }: flake-utils.lib.eachDefaultSystem (system:
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
    in
    {
      packages.default = sanic;
      formatter = pkgs.nixpkgs-fmt;
      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          go-tools # staticcheck
          gomod2nix.packages.${system}.default
        ];
        packages = with pkgs; [
          mpd
          mpc-cli
        ];
      };
      nixosConfigurations."sanic" = nixpkgs.lib.nixosSystem {
        inherit system;
        modules = [
          {
            environment.systemPackages = with pkgs; [
              sanic
            ];
          }
          "${nixpkgs}/nixos/modules/virtualisation/proxmox-lxc.nix"
          ./configuration.nix
        ];
      };
    }
  );
}
