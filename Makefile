.PHONY: build_arm clean

# Variables
PATHROOT = $(shell pwd)
BINARY_NAME = app.exe
MAC_BINARY_NAME = app

download_x264: clean
	@echo "Downloading x264..."
	wget https://code.videolan.org/videolan/x264/-/archive/31e19f92f00c7003fa115047ce50978bc98c3a0d/x264-31e19f92f00c7003fa115047ce50978bc98c3a0d.tar.gz
	tar -xvf x264-31e19f92f00c7003fa115047ce50978bc98c3a0d.tar.gz
	rm x264-31e19f92f00c7003fa115047ce50978bc98c3a0d.tar.gz
	mv -f x264-31e19f92f00c7003fa115047ce50978bc98c3a0d .x264src

build_x264_arm: clean
	@echo "Building x264..."
	docker build . -f ./.docker/builder_x264_arm.Dockerfile -t x264builder
	docker run -d --name x264builder x264builder
	mkdir -p x264/lib
	docker cp x264builder:"/builder/libx264.a" $(PATHROOT)/x264/lib
	docker cp x264builder:/builder/x264.pc $(PATHROOT)/x264/lib
	cp .x264src/x264.h $(PATHROOT)/x264/lib
	cp .x264src/x264cli.h $(PATHROOT)/x264/lib
	cp x264/x264_config.h $(PATHROOT)/x264/lib
#	docker rm -f x264builder
#	docker rmi -f x264builder

build_arm: 
	@echo "Building..."
	docker build . -f ./.docker/builder_arm.Dockerfile -t winbalbuilder
	docker run -d --name winbalbuilder winbalbuilder
	mkdir -p bin
	docker cp winbalbuilder:/app/$(BINARY_NAME) $(PATHROOT)/bin/
	docker rm -f winbalbuilder
	docker rmi -f winbalbuilder

build_mac:
	@echo "Building..."
	go build -o $(PATHROOT)/bin/$(MAC_BINARY_NAME) $(PATHROOT)/cmd/macosapp/main.go

clean:
	@echo "Cleaning..."
	rm -rf $(PATHROOT)/bin/$(BINARY_NAME)
	docker rm -f winbalbuilder
	docker rm -f x264builder
	docker rmi -f winbalbuilder
	docker rmi -f x264builder

clean_all: clean
	@echo "Cleaning all..."
	rm -rf $(PATHROOT)/.x264src
	rm -rf $(PATHROOT)/x264/lib