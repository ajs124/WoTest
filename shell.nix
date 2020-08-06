with import <nixpkgs> { }; # TODO: pin nixpkgs to some revision/release

stdenv.mkDerivation {
  name = "go";
  buildInputs = [
    # for out code
    go
    # for the implementations
    ## node-wot
    nodejs
    nodePackages.npm
    ## wot-py
    (python3.withPackages (p: with p; [
      pip
      setuptools
      netifaces
      numpy
      sphinx
      sphinx_rtd_theme
      tornado
      virtualenvwrapper
    ]))
    ## sane-city
    maven
    jdk12
  ];
  shellHook = ''
    export GOPATH=$PWD/gopath

    export SOURCE_DATE_EPOCH=315532800
    alias pip="PIP_PREFIX='$PWD/.pip_packages' \pip"
    export PYTHONPATH="$PWD/.pip_packages/lib/${pkgs.python3.libPrefix}/site-packages:$PYTHONPATH"
    cd src
  '';
}
