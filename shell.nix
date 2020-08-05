with import <nixpkgs> { };

stdenv.mkDerivation {
  name = "go";
  buildInputs = [
    delve
    go
  ];
  shellHook = ''
    export GOPATH=$PWD/gopath
  '';
}
