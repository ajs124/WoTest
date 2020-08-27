FROM nixos/nix AS wotest-base

RUN adduser mosquitto -D -H
# arbitrary release-20.03 commit from 2020-08-19
RUN nix-channel --add https://github.com/NixOS/nixpkgs/archive/f0924dbf552e28ee0462b180116135c187eb41b4.tar.gz nixpkgs && nix-channel --update
# download and put all dependencies in a docker layer
COPY shell.nix .
RUN nix-shell --run true

FROM wotest-base AS wotest-wot-py
COPY implementations implementations
RUN nix-shell --run ./implementations/prepare_wot-py.sh

FROM wotest-base AS wotest-node-wot
COPY implementations implementations
RUN nix-shell --run ./implementations/prepare_node-wot.sh

FROM wotest-base AS wotest-sane-city
COPY .git/modules .git/modules
COPY implementations implementations
RUN nix-shell --run ./implementations/prepare_sane-city.sh

FROM wotest-base
# copy previous stages
COPY --from=wotest-wot-py implementations/wot-py implementations/
COPY --from=wotest-node-wot implementations/node-wot implementations/
COPY --from=wotest-node-wot implementations/node-wot-install implementations/
COPY --from=wotest-sand-city implementations/sane-city implementations/
# copy rest, build/run this
COPY go.mod src tests .
RUN nix-shell --run go run
