{ config, pkgs, lib, ... }:

{
  networking = {
    hostName = "sanic";
    dhcpcd.enable = false;
    useDHCP = lib.mkDefault true;
    useNetworkd = true;
    #nftables.enable = true;
  };
  services.resolved = {
    enable = true;
    llmnr = "true";
    dnssec = "allow-downgrade";
    dnsovertls = "opportunistic";
    fallbackDns = [
      "2606:4700:4700::1111#one.one.one.one"
      "1.1.1.1#one.one.one.one"
    ];
  };

  time.timeZone = "Etc/UTC";

  i18n.defaultLocale = "en_US.UTF-8";
  console = {
    font = "Lat2-Terminus16";
    useXkbConfig = true;
  };

  nix = {
    optimise = {
      automatic = true;
      dates = [ "00:00" ];
    };
    gc = {
      automatic = true;
      dates = "weekly";
      options = "--delete-older-than 7d";
    };
    settings = {
      auto-optimise-store = true;
      trusted-users = [ "root" "@wheel" ];
      experimental-features = [ "nix-command" "flakes" ];
    };
  };

  environment.systemPackages = with pkgs; [ vim git ];

  services = {
    qemuGuest.enable = true;
    openssh = {
      enable = true;
      settings = {
        PasswordAuthentication = false;
        KbdInteractiveAuthentication = false;
      };
      banner = ''
                                __
          ____     __      ___ /\_\    ___
         /',__\  /'__`\  /' _ `\/\ \  /'___\
        /\__, `\/\ \L\.\_/\ \/\ \ \ \/\ \__/
        \/\____/\ \__/.\_\ \_\ \_\ \_\ \____\
         \/___/  \/__/\/_/\/_/\/_/\/_/\/____/
      '';
    };
  };

  programs = {
    ssh.startAgent = true;
    vim = {
      enable = true;
      defaultEditor = true;
    };
  };

  security.sudo.wheelNeedsPassword = false;

  system.stateVersion = "25.05";
}

