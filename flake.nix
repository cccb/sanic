{
  description = "sanic - chaos music control";
  inputs = {
    nixpkgs.url = github:NixOS/nixpkgs/nixpkgs-unstable;
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
      packages.default = sanic;
      nixosModules.default = { config, lib, pkgs, options, ... }:
      let
        cfg = config.services.sanic;
        configFile = pkgs.writeText "config.ini" ''
          [ui]
          host=${cfg.ui.host}
          port=${cfg.ui.port}
          tls=${cfg.ui.tls}
          certificate=${cfg.ui.certificate}
          key=${cfg.ui.key}

          [mpd]
          host=${cfg.backend.host}
          port=${cfg.backend.port}
        '';
        execCommand = "${cfg.package}/bin/sanic -c '${configFile}'";
      in
      {
        options.services.sanic = {
          enable = lib.mkEnableOption "Enables the sanic systemd service.";
          package = lib.mkOption {
            description = "Package to use.";
            type = lib.types.package;
            default = sanic;
          };
          ui = lib.mkOption {
            description = "Setting for HTTP(S) UI.";
            example = lib.literalExpression ''
              {
                host = "[::1]";
                port = 443;
                tls = true;
                certificate = "${config.security.acme.certs."sanic.example.com".directory}/fullchain.pem";
                key = "${config.security.acme.certs."sanic.example.com".directory}/key.pem";
              }
            '';
            default = {
              host = "[::1]";
              port = 80;
              tls = false;
            };
            type = lib.types.submodule {
              options = {
                host = lib.mkOption {
                  type = lib.types.str;
                  default = "[::1]";
                  description = "Host to bind to.";
                };
                port = lib.mkOption {
                  type = lib.types.port;
                  default = 80;
                  description = "Port to listen on.";
                };
                tls = lib.mkOption {
                  type = lib.types.bool;
                  default = false;
                  description = "Enables HTTPS.";
                };
                certificate = lib.mkOption {
                  type = lib.types.nullOr lib.types.path;
                  default = null;
                  description = "Path to TLS certificate for HTTPS.";
                };
                key = lib.mkOption {
                  type = lib.types.nullOr lib.types.path;
                  default = null;
                  description = "Path to TLS key for HTTPS.";
                };
              };
            };
          };
          backend = lib.mkOption {
            description = "Configure MPD backend.";
            example = lib.literalExpression ''
              {
                host = "localhost";
                port = 6600;
              }
            '';
            default = {
              host = "localhost";
              port = 6600;
            };
            type = lib.types.submodule {
              options = {
                host = lib.mkOption {
                  type = lib.types.str;
                  default = "localhost";
                  description = "Hostname or IP of MPD instance.";
                };
                port = lib.mkOption {
                  type = lib.types.port;
                  default = 6600;
                  description = "Port of MPD instance.";
                };
              };
            };
          };
        };

        config = lib.mkIf cfg.enable {
          systemd.services."sanic" = {
            description = "sanic - chaos music control";
            wants = [ "network-online.target" ];
            after = [ "network-online.target" ];
            serviceConfig = {
              Restart = "always";
              RestartSec = 30;
              ExecStart = execCommand;
              User = "sanic";
              Group = "sanic";
              AmbientCapabilities = lib.mkIf (cfg.ui.port < 1000) [ "CAP_NET_BIND_SERVICE" ];
              CapabilityBoundingSet = lib.mkIf (cfg.ui.port < 1000) [ "CAP_NET_BIND_SERVICE" ];
              NoNewPrivileges = true;
            };
            wantedBy = [ "multi-user.target" ];
          };
        };

        #meta = {
        #  maintainers = with lib.maintainers; [ xengi ];
        #  doc = ./default.xml;
        #};
      };
    }
  );
}
