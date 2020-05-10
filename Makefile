EXECUTABLE=memsfcr
APPNAME=MemsFCR

WINDOWSDISTPATH=dist/windows
WINDOWS=$(WINDOWSDISTPATH)/$(EXECUTABLE).exe
LINUXDISTPATH=dir/linux
LINUX=$(LINUXDISTPATH)/$(EXECUTABLE)
DARWINDISTPATH=dist/darwin
DARWIN=$(DARWINDISTPATH)/$(EXECUTABLE)
ARMDISTPATH=dist/arm
ARM=$(ARMDISTPATH)/$(EXECUTABLE)-arm
VERSION=$(shell git describe --tags --always --long --dirty)

.PHONY: all clean

all: build

build: darwin   ## Build binaries
	@echo version: $(VERSION)

darwin: $(DARWIN) buildapp ## Build for Darwin (macOS 10.15+)
arm: $(ARM) ## Build for Darwin 32bit (macOS <10.15)
linux: $(LINUX) ## Build for Linux
windows: $(WINDOWS) ## Build for Windows

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -i -v -o $(WINDOWS) -ldflags="-H windowsgui -s -w -X main.version=$(VERSION)"

$(LINUX):
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -i -v -o $(LINUX) -ldflags="-s -w -X main.version=$(VERSION)"

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -i -v -o $(DARWIN) -ldflags="-s -w -X main.version=$(VERSION)"

$(ARM):
	env GOOS=android GOARCH=arm GOARM=5 CGO_ENABLED=1 go build -i -v -o $(DARWIN) -ldflags="-extldflags=-Wl,-z,norelro"

buildapp:
	mkdir "$(DARWINDISTPATH)/$(APPNAME).app"
	mkdir "$(DARWINDISTPATH)/$(APPNAME).app/Contents"
	mkdir "$(DARWINDISTPATH)/$(APPNAME).app/Contents/MacOS"
	mkdir "$(DARWINDISTPATH)/$(APPNAME).app/Contents/Resources"
	mkdir "$(DARWINDISTPATH)/$(APPNAME).app/Contents/MacOS/logs"

	cp resources/icons/icon.icns "$(DARWINDISTPATH)/$(APPNAME).app/Contents/Resources"
	cp resources/darwin/Info.plist "$(DARWINDISTPATH)/$(APPNAME).app/Contents"
	mv $(DARWIN) "$(DARWINDISTPATH)/$(APPNAME).app/Contents/MacOS/$(EXECUTABLE)"
	cp memsfcr.cfg "$(DARWINDISTPATH)/$(APPNAME).app/Contents/MacOS"
	cp -r ./public "$(DARWINDISTPATH)/$(APPNAME).app/Contents/MacOS"
	ln -s /Applications "$(DARWINDISTPATH)/Applications"

	sips -i resources/icons/icon.png
	DeRez -only icns resources/icons/icon.png > resources/icons/icns.rsrc
	hdiutil create /tmp/tmp.dmg -ov -volname "MemsFCR" -fs HFS+ -srcfolder "$(DARWINDISTPATH)" 
	hdiutil convert /tmp/tmp.dmg -format UDZO -o "$(DARWINDISTPATH)/$(APPNAME).dmg"
	Rez -append resources/icons/icns.rsrc -o "$(DARWINDISTPATH)/$(APPNAME).dmg"
	SetFile -a C "$(DARWINDISTPATH)/$(APPNAME).dmg"
	mv "$(DARWINDISTPATH)/$(APPNAME).dmg" dist/$(APPNAME).dmg

clean: ## Remove previous build
	rm -f $(WINDOWS) $(LINUX) $(DARWIN)
	rm -fr $(DARWINDISTPATH)/$(APPNAME).app
	rm -fr $(DARWINDISTPATH)/$(APPNAME).dmg
	rm -f $(DARWINDISTPATH)/Applications

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
