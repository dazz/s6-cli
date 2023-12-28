
let
   nixpkgs = import <nixpkgs> {};
   inherit (nixpkgs) pkgs; # AKA pkgs = nixpkgs.pkgs

   s6-cli = pkgs.buildGoModule {
     name = "s6-cli";
     src = ./.;
     vendorHash = "sha256-XjmoKtlcynjSLdikWVGdsn2y3SY3ximWmoe2wEMhfXs=";
   };
in

pkgs.mkShell {
    packages = [
     s6-cli
    ];
}