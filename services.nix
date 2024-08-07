{ self, ...}: {config, lib, pkgs, ...}:

let
  cfg = config.services.sanic;
  format = pkgs.formats.ini { };
in
{
  options.services.sanic = {
    enable = mkEnableOption (lib.mdDoc "sanic");
    settings = mkOption {
      type = format.type;
      default = { };
      description = lib.mkDoc ''
      '';
    };
  };

  config = mkIf cfg.enable {
    systemd.services.sanic = {
      description = "chaos music control";
      wantedBy = [ "multi-user.target" "default.target" ];
      serviceConfig = {
        DynamicUser = true;
        ExecStart = "${self.packages.${pkgs.system}.default}/bin/sanic";
        Restart = "on-failure";
        AmbientCapabilities = [ "CAP_NET_BIND_SERVICE" ];
      };
    };
  };
}
