{
  description = "chaos music control";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:tweag/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };
  outputs = { self, nixpkgs, flake-utils, gomod2nix }: flake-utils.lib.eachDefaultSystem (system:
    let
      #pkgs = nixpkgs.legacyPackages.${system};
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
    in {
      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gomod2nix.packages.${system}.default
        ];
        packages = with pkgs; [
          mpd
          mpc-cli
          sanic
        ];
      };
    }
  );
}
