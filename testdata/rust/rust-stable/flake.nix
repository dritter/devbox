{
    description = "A devbox shell";
    outputs = { self, nixpkgs }:
        # TODO generalize system using flake-util
        let pkgs = nixpkgs.legacyPackages.x86_64-darwin;

        in {
            devShell.x86_64-darwin = pkgs.mkShell {
              shellHook =
                ''
                  echo "Starting a devbox shell..."

                  # We're technically no longer in a Nix shell after this hook because we
                  # exec a devbox shell.
                  export IN_NIX_SHELL=0
                  export DEVBOX_SHELL_ENABLED=1

                  # Undo the effects of `nix-shell --pure` on SSL certs.
                  # See https://github.com/NixOS/nixpkgs/blob/dae204faa0243b4d0c0234a5f5f83a2549ecb5b7/pkgs/stdenv/generic/setup.sh#L677-L685
                  if [ "$NIX_SSL_CERT_FILE" == "/no-cert-file.crt" ]; then
                     unset NIX_SSL_CERT_FILE
                  fi
                  if [ "$SSL_CERT_FILE" == "/no-cert-file.crt" ]; then
                     unset SSL_CERT_FILE
                  fi

                  # Append the parent shell's PATH so that we retain access to
                  # non-Nix programs, while still preferring the Nix ones.
                  export "PATH=$PATH:$PARENT_PATH"

                  
                '';
              buildInputs = with pkgs; [
                    rustup
                  
                    libiconv
                  ];
            };
        };
}
