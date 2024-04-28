FROM amd64/ubuntu:22.04

WORKDIR /builder

COPY ./.vpxsrc .

RUN dpkg --add-architecture amd64 
RUN apt-get update
RUN apt-get install -y gcc-mingw-w64
RUN apt-get install -y g++-mingw-w64
RUN apt-get install -y make
RUN apt-get install -y yasm

# build vpx library
RUN CC=x86_64-w64-mingw32-gcc \
    CXX=x86_64-w64-mingw32-g++ \
    AR=x86_64-w64-mingw32-ar \
    CROSS=x86_64-w64-mingw32- \
    ./configure \
        --target=x86_64-win64-gcc \
        --disable-multithread \
        --enable-vp8 \
        --enable-vp9 \
        --disable-examples && \
    make
RUN make install

# wait for the container to copy the binary
CMD ["tail", "-f", "/dev/null"]