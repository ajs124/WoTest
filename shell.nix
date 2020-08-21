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
      tornado
      jsonschema
      six
      rx
      python-slugify
      hbmqtt
      zeroconf

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
  '';
}
