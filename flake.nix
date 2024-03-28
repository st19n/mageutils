{
	description = "mageutils";

	inputs = {
		nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
	};

	outputs = { self, nixpkgs }:
	let
		system = "x86_64-linux";
		pkgs = nixpkgs.legacyPackages.${system};
	in
	{
		devShells.${system}.default =
			pkgs.mkShell
				{
					buildInputs = with pkgs; [
						go
						gopls
					];

					shellHook = ''
						export PATH="$(pwd)/.tmp/bin:$PATH"
					'';
			};
	};
}
