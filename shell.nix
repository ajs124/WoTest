with import <nixpkgs> { }; # TODO: pin nixpkgs to some revision/release

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
    export SOURCE_DATE_EPOCH=315532800
    alias pip="PIP_PREFIX='$PWD/.pip_packages' \pip"
    export PYTHONPATH="$PWD/.pip_packages/lib/${pkgs.python3.libPrefix}/site-packages:$PYTHONPATH"
    cd src
  '';
}
