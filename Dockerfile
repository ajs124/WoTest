FROM nixos/nix

RUN adduser mosquitto -D -H
# arbitrary release-20.03 commit from 2020-08-19
RUN nix-channel --add https://github.com/NixOS/nixpkgs/archive/f0924dbf552e28ee0462b180116135c187eb41b4.tar.gz nixpkgs && nix-channel --update
# download and put all dependencies in a docker layer
COPY shell.nix .
RUN nix-shell --run true
# copy and prepare implementations
COPY implementations implementations
RUN nix-shell --run ./implementations/prepare.sh
# copy rest, build/run this
COPY go.mod src tests .
RUN nix-shell --run go run
