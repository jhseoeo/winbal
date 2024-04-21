FROM --platform=linux/amd64 golang:1.22-bullseye

WORKDIR /app

COPY go.mod ./
RUN go mod download

# install libraries for cross-compilation
RUN dpkg --add-architecture amd64 
RUN apt-get update
RUN apt-get install -y gcc-mingw-w64

COPY . .  

# enable CGO and flags to cross-compile for windows
ENV CGO_ENABLED=1 GOOS=windows GOARCH=amd64
ENV CC=x86_64-w64-mingw32-gcc
RUN go build -o app.exe ./cmd/windowapp

# wait for the container to copy the binary
CMD ["tail", "-f", "/dev/null"]