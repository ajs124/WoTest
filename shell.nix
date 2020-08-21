with import <nixpkgs> { };

stdenv.mkDerivation {
  name = "go";
  buildInputs = [
    # delve
    # our code
    go
    # node-wot
    nodejs
    nodePackages.npm
    # sane-city
    jdk12_headless
    maven
    # wot-py
    (python3.withPackages (p: with p; [
      setuptools
      pip

      # for tests
      pytest
      faker
      aiozeroconf
      mock
      pyopenssl
    ]))
    # unstable.jetbrains.goland
    mosquitto
    curl
  ];
  shellHook = ''
    export GOPATH=$PWD/gopath
    export NODE_PATH=$PWD/implementations/node-wot-install/node_modules:$NODE_PATH
    export PYTHONPATH=$PWD/implementations/wot-py/install:$PYTHONPATH
  '';
}
