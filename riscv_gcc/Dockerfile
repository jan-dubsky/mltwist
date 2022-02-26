FROM archlinux:latest

RUN pacman -Syyu --noconfirm autoconf automake curl python3 libmpc mpfr gmp gawk base-devel bison flex texinfo gperf libtool patchutils bc zlib expat
RUN pacman -S --noconfirm git
RUN git clone https://github.com/riscv-collab/riscv-gnu-toolchain.git
WORKDIR riscv-gnu-toolchain

RUN ./configure --prefix=/opt/riscv --with-arch=rv64ima
RUN make linux
ENV PATH="/opt/riscv/bin:${PATH}"

RUN mkdir /build
WORKDIR /build
