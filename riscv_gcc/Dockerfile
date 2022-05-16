FROM ubuntu:22.04

RUN apt-get update && apt-get -y upgrade && apt-get -y install git autoconf automake autotools-dev curl python3 libmpc-dev libmpfr-dev libgmp-dev gawk build-essential bison flex texinfo gperf libtool patchutils bc zlib1g-dev libexpat-dev
RUN git clone https://github.com/riscv-collab/riscv-gnu-toolchain.git
WORKDIR riscv-gnu-toolchain

ENV LC_ALL=C

# Initialization of submodules takes quite long to be worth of caching in a
# separate layer.
RUN git submodule update --init riscv-gcc riscv-binutils glibc
RUN ./configure --prefix=/opt/riscv --with-arch=rv64ima
RUN make linux
ENV PATH="/opt/riscv/bin:${PATH}" CC="riscv64-unknown-linux-gnu-gcc -misa-spec=2.2 -mabi=lp64 -march=rv64im"

# Install dependencies required to bootstrap grep.
RUN apt-get -y install rsync wget gettext autopoint pkg-config

RUN useradd -N -d /build -K UID_MAX=1000000 -u 1000000 build
USER build
WORKDIR /build
