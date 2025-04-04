{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.utils.url = "github:numtide/flake-utils";
  inputs.gomod2nix = {
    url = "github:tweag/gomod2nix";
    inputs.nixpkgs.follows = "nixpkgs";
  };
  outputs =
    {
      self,
      nixpkgs,
      utils,
      gomod2nix,
    }:
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ gomod2nix.overlays.default ];
        };
        name = "todo";
      in
      {
        packages.default = pkgs.buildGoApplication {
          pname = name;
          version = "1.0.2";
          src = ./.;
          modules = ./gomod2nix.toml;
        };
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            gomod2nix.packages.${system}.default
            go
            gotools
            gopls
          ];
        };
        nixosModules.default =
          { lib, config, ... }:
          let
            cfg = config.services.${name};
          in
          {
            options.services.${name} = {
              enable = lib.mkEnableOption "this ${name} service";

              package = lib.mkOption {
                type = lib.types.package;
                default = self.packages.${pkgs.system}.default;
                description = "package to use for the service";
              };

              port = lib.mkOption {
                type = lib.types.port;
                default = 8080;
                description = "port to listen on";
              };

              key = lib.mkOption {
                type = lib.types.str;
                default = "todo";
                description = "the key for storing todos";
              };
            };

            config = lib.mkIf cfg.enable {
              environment.systemPackages = [
                (pkgs.writeShellScriptBin "todocli" ''
                  export TODO_KEY="${cfg.key}"; ${cfg.package}/bin/todocli $@
                '')
              ];
              systemd.services.${name} = {
                wantedBy = [ "multi-user.target" ];
                serviceConfig = {
                  ExecStart = "${cfg.package}/bin/todosrv -key ${cfg.key} -port :${toString cfg.port}";
                  Restart = "on-failure";
                  RestartSec = "5s";
                };
              };
              services.redis.servers."".enable = true;
            };
          };
      }
    );
}
