EXECUTABLE=memsfcr
APPNAME=MemsFCR
WINDOWS=dist/$(EXECUTABLE)_windows_amd64.exe
LINUX=dist/$(EXECUTABLE)_linux_amd64
DARWIN=dist/$(EXECUTABLE)_darwin_amd64
VERSION=$(shell git describe --tags --always --long --dirty)

.PHONY: all clean

all: build

build: darwin windows  ## Build binaries
	@echo version: $(VERSION)

darwin: $(DARWIN) buildapp ## Build for Darwin (macOS)
linux: $(LINUX) ## Build for Linux
windows: $(WINDOWS) ## Build for Windows

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -i -v -o $(WINDOWS) -ldflags="-H windowsgui -s -w -X main.version=$(VERSION)"

$(LINUX):
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -i -v -o $(LINUX) -ldflags="-s -w -X main.version=$(VERSION)"

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -i -v -o $(DARWIN) -ldflags="-s -w -X main.version=$(VERSION)"

buildapp:
	mkdir "dist/$(APPNAME).app"
	mkdir "dist/$(APPNAME).app/Contents"
	mkdir "dist/$(APPNAME).app/Contents/MacOS"
	mkdir "dist/$(APPNAME).app/Contents/Resources"
	mkdir "dist/$(APPNAME).app/Contents/MacOS/logs"
	cp resources/icons/icon.icns "dist/$(APPNAME).app/Contents/Resources"
	cp resources/Info.plist "dist/$(APPNAME).app/Contents"
	cp $(DARWIN) "dist/$(APPNAME).app/Contents/MacOS/$(EXECUTABLE)"
	cp memsfcr.cfg "dist/$(APPNAME).app/Contents/MacOS"
	cp -r ./public "dist/$(APPNAME).app/Contents/MacOS"

clean: ## Remove previous build
	rm -f $(WINDOWS) $(LINUX) $(DARWIN)
	rm -fr dist/$(APPNAME).app

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
