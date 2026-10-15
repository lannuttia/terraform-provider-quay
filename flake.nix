{
  description = "Quay Terraform Provider development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfreePredicate = pkg:
            builtins.elem (nixpkgs.lib.getName pkg) [ "terraform" ];
        };
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            # Go toolchain — matches the version declared in go.mod
            go_1_27

            # Linting and formatting
            golangci-lint
            gotools       # goimports, godoc, etc.

            # Language server for editor integration
            gopls

            # Terraform CLI — used by tfplugindocs and acceptance tests
            terraform

            # Release tooling
            goreleaser

            # Version control
            git

            # Pre-commit hooks
            pre-commit
          ];

          shellHook = ''
            # Point GOROOT at the Nix-provided Go so that go vet, golangci-lint,
            # and other tools don't pick up a different system installation.
            export GOROOT="${pkgs.go_1_27}/share/go"
            echo "quay-terraform-provider dev shell"
            echo "  go            $(go version)"
            echo "  terraform     $(terraform version -json | ${pkgs.jq}/bin/jq -r '.terraform_version')"
            echo "  golangci-lint $(golangci-lint version --short)"
            pre-commit install --install-hooks 2>/dev/null
          '';
        };
      });
}
