.PHONY: build_arm clean

# Variables
PATHROOT = $(shell pwd)
BINARY_NAME = app.exe
MAC_BINARY_NAME = app

download_vpx: clean
	@echo "Downloading vpx"
	git clone https://chromium.googlesource.com/webm/libvpx .vpxsrc

build_vpx_arm: clean
	@echo "Building vpx..."
	docker build . -f ./.docker/builder_vpx_arm.Dockerfile -t vpxbuilder
	docker run -d --name vpxbuilder vpxbuilder
	mkdir -p vpx/lib
	docker cp vpxbuilder:"/builder/libvpx.a" $(PATHROOT)/vpx/lib
	cp .vpxsrc/*.h $(PATHROOT)/vpx/lib
#	docker rm -f vpxbuilder
#	docker rmi -f vpxbuilder

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
	go build -o $(PATHROOT)/bin/$(MAC_BINARY_NAME) $(PATHROOT)/cmd/app/...

clean:
	@echo "Cleaning..."
	rm -rf $(PATHROOT)/bin/$(BINARY_NAME)
	docker rm -f winbalbuilder
	docker rm -f vpxbuilder
	docker rmi -f winbalbuilder
	docker rmi -f vpxbuilder

clean_all: clean
	@echo "Cleaning all..."
	rm -rf $(PATHROOT)/.vpxsrc
	rm -rf $(PATHROOT)/vpx/lib