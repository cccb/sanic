{ lib, pkgs, ... }:

{
  networking = {
    hostName = "sanic";
    useNetworkd = true;
    nameservers = [
      "172.23.42.1"
    ];
    defaultGateway = {
      address = "172.23.42.1";
      interface = "eth0";
    };
    interfaces.eth0 = {
      ipv4.addresses = [{
        address = "172.23.43.102";
        prefixLength = 23;
      }];
    };
  };
  services.resolved = {
    enable = true;
    llmnr = "true";
    dnssec = "allow-downgrade";
    dnsovertls = "opportunistic";
  };
  time.timeZone = "Europe/Berlin";
  i18n.defaultLocale = "en_US.UTF-8";

  users.users.xengi = {
    isNormalUser = true;
    extraGroups = [ "wheel" ];
    shell = pkgs.fish;
    openssh.authorizedKeys.keys = [
      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICW1+Ml8R9x1LCJaZ8bIZ1qIV4HCuZ6x7DziFW+0Nn5T xengi@kanae_2022-12-09"
      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICmb+mJfo84IagUaRoDEqY9ROjjQUOQ7tMclpN6NDPrX xengi@kota_2022-01-16"
      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICyklb7dvEHH0VBEMmTUQFKHN6ekBQqkDKj09+EilUIQ xengi@lucy_2018-09-08"
      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGhyfD+8jMl6FDSADb11sfAsJk0KNoVzjjiDRZjUOtmf xengi@nana_2019-08-16"
      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICjv9W8WXq9QGkgmANNPQR24/I1Pm1ghxNIHftEI+jlZ xengi@mayu_2021-06-11"
      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMPtGqhV7io3mhIoZho4Yf7eCo0sUZvjT2NziM2PkXSo xengi@nyu_2017-10-11"
    ];
    packages = with pkgs; [
      kitty # for terminfo
    ];
  };

  nix = {
    optimise = {
      automatic = true;
      dates = [ "00:00" ];
    };
    settings = {
      auto-optimise-store = true;
      experimental-features = [ "nix-command" "flakes" ];
    };
    gc = {
      automatic = true;
      options = "--delete-older-than 10d";
    };
  };

  environment.systemPackages = with pkgs; [
    git # required for flakes
    vim
    nvd
  ];

  services = {
    openssh = {
      enable = true;
      settings.PasswordAuthentication = false;
    };
  };

  programs = {
    fish = {
      enable = true;
      interactiveShellInit = ''
        function upgrade --description "Upgrade NixOS system"
          cd /etc/nixos
          nix flake update
          cd -
          nixos-rebuild switch --upgrade
          nvd diff (ls -d1v /nix/var/nix/profiles/system-*-link|tail -n 2)
        end
      '';
    };
    vim.defaultEditor = true;
    mtr.enable = true;
  };

  security = {
    sudo.execWheelOnly = true;
  };

  networking.firewall = {
    enable = true;
    allowedTCPPorts = [
      80 # HTTP/1
      443 # HTTP/2
    ];
    allowedUDPPorts = [
      443 # HTTP/3
    ];
  };

  system.stateVersion = "24.05";
}

