FROM amd64/ubuntu:22.04

WORKDIR /builder

COPY ./.x264src .

RUN dpkg --add-architecture amd64 
RUN apt-get update
RUN apt-get install -y gcc-mingw-w64
RUN apt-get install -y make

# build x264 library
RUN ./configure --host=x86_64-w64-mingw32 --disable-cli --enable-shared --disable-asm \
                --cross-prefix=x86_64-w64-mingw32- && \
    make && make install-lib-static

RUN mv ./libx264*.dll ./libx264.dll

# wait for the container to copy the binary
CMD ["tail", "-f", "/dev/null"]