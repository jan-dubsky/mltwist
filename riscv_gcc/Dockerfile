FROM archlinux:latest

RUN pacman -Syyu --noconfirm autoconf automake curl python3 libmpc mpfr gmp gawk base-devel bison flex texinfo gperf libtool patchutils bc zlib expat
RUN pacman -S --noconfirm git
RUN git clone https://github.com/riscv-collab/riscv-gnu-toolchain.git
WORKDIR riscv-gnu-toolchain

RUN ./configure --prefix=/opt/riscv --with-arch=rv64ima
RUN make linux
ENV PATH="/opt/riscv/bin:${PATH}"
ENV CC="riscv64-unknown-linux-gnu-gcc -misa-spec=2.2 -mabi=lp64 -march=rv64im"

RUN pacman -S --noconfirm rsync wget

RUN useradd -N -K UID_MAX=1000000 -u 1000000 build
RUN mkdir /build && chown build /build
USER build
WORKDIR /build
