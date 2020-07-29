EXECUTABLE=memsfcr
APPNAME=MemsFCR
DEVID="Developer ID Application: Andrew Jackson (MD9E767XF5)"

DISTPATH=dist
RESOURCESPATH=resources

WINDOWSDISTPATH=dist/windows
WINDOWS=$(WINDOWSDISTPATH)/$(EXECUTABLE).exe

LINUXDISTPATH=dir/linux
LINUX=$(LINUXDISTPATH)/$(EXECUTABLE)

DARWINDISTPATH=dist/darwin
DARWIN=$(DARWINDISTPATH)/$(EXECUTABLE)

ARMDISTPATH=dist/arm
ARM=$(ARMDISTPATH)/$(EXECUTABLE)-arm

#VERSION=$(shell git describe --tags)
VERSION="V1.5.3"
BUILD=$(shell date +%FT%T%z)

.PHONY: all clean

all: build

build: darwin   ## Build binaries
	@echo version: $(VERSION)
	@echo appid: $(DEVID)

darwin: $(DARWIN) buildapp signapp packageapp ## Build for Darwin (macOS 10.15+)
arm: $(ARM) ## Build for Darwin 32bit (macOS <10.15)
linux: $(LINUX) ## Build for Linux
windows: $(WINDOWS) ## Build for Windows

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -i -v -o $(WINDOWS) -ldflags="-H windowsgui -s -w -X main.version=$(VERSION)"

$(LINUX):
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -i -v -o $(LINUX) -ldflags="-s -w -X main.version=$(VERSION)"

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -i -v -o $(DARWIN) -ldflags="-s -w -X main.Version=$(VERSION) -X main.Build=$(BUILD)"

$(ARM):
	env GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ go 

buildapp:
	# create the MacOS app
	cp -f "$(DARWINDISTPATH)/$(EXECUTABLE)" "$(RESOURCESPATH)/$(EXECUTABLE)"	
	./macapp -assets "$(RESOURCESPATH)" -bin $(EXECUTABLE) -icon "$(RESOURCESPATH)/icons/icon.png" -identifier "com.github.andrewdjackson.memsfcr" -name "$(APPNAME)" -o "$(DARWINDISTPATH)"

signapp:	
	# sign with the Developer ID
	codesign --force  --deep --verify --verbose -s $(DEVID) -v --timestamp --options runtime "$(DARWINDISTPATH)/$(APPNAME).app/Contents/MacOS/$(EXECUTABLE)" "$(DARWINDISTPATH)/$(APPNAME).app"


packageapp:
	appdmg $(DISTPATH)/dmgspec.json $(DARWINDISTPATH)/MemsFCR.dmg 

	# check signature with:
	#   codesign -display --deep -vvv $(APPNAME).app
	#
	# the app will need notarizing with the following command:
	#   xcrun altool --notarize-app -f $(APPNAME).dmg --primary-bundle-id "{bundle-id}" -u {username} -p {password}
	#
	# if successful 'staple' the app for offline installation
	#   xcrun stapler staple "$(APPNAME).app"
	#   xcrun stapler staple "$(APPNAME).dmg"
	

clean: ## Remove previous build
	rm -f $(WINDOWS) $(LINUX) $(DARWIN)
	rm -f $(DARWINDISTPATH)/Applications
	rm -fr $(DARWINDISTPATH)/*
	rm -fr $(DISTPATH)/$(APPNAME).dmg

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
