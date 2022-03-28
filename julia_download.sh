#!/usr/bin/env bash
ARCH=$(uname -m)

if [[ "${ARCH}" == "aarch64" ]]; then
	wget -q https://julialang-s3.julialang.org/bin/linux/aarch64/1.7/julia-1.7.2-linux-aarch64.tar.gz
	tar -zxf julia-1.7.2-linux-aarch64.tar.gz
	rm -rf julia-1.7.2-linux-aarch64.tar.gz
elif [[ "${ARCH}" == "x86_64" ]]; then
	wget -q https://julialang-s3.julialang.org/bin/linux/x64/1.7/julia-1.7.2-linux-x86_64.tar.gz
	tar -zxf julia-1.7.2-linux-x86_64.tar.gz
	rm -rf julia-1.7.2-linux-x86_64.tar.gz
else
	echo "unsupported platform ${ARCH}"
	exit 1
fi

mv julia-1.7.2 /usr/local/julia
ln -s /usr/local/julia/bin/julia /usr/local/bin

