EXECUTABLE=memsfcr
APPNAME=MemsFCR
DISTPATH=dist
RESOURCESPATH=resources
DARWINDISTPATH=dist/darwin
EXEPATH=$(DARWINDISTPATH)/$(EXECUTABLE)
MIN_DEPLOYMENT_TARGET=-mmacosx-version-min=11.6.7

DEVID="Developer ID Application: Andrew Jackson (MD9E767XF5)"
LOCAL_DISTID="Developer ID Application: Andrew Jackson (MD9E767XF5)"
LOCAL_INSTID="Developer ID Installer: Andrew Jackson (MD9E767XF5)"
STORE_DISTID="3rd Party Mac Developer Application: Andrew Jackson (MD9E767XF5)"
STORE_INSTID="3rd Party Mac Developer Installer: Andrew Jackson (MD9E767XF5)"

VERSION=$(shell cat version)
BUILD=$(shell date +%FT%T%z)

.PHONY: clean

package: package_local_app notarize_local_package
build: build_darwin create_darwin_app package_local_app
local: build_darwin create_darwin_app sign_app_local package_local_app notarize_local_package
store: build_darwin create_darwin_app sign_app_store upload_app_to_store

build_darwin:
	# Version: $(VERSION)
	# Build: $(BUILD)
	env GOOS=darwin GOARCH=amd64 CGO_CFLAGS="$(MIN_DEPLOYMENT_TARGET)" CGO_LDFLAGS="$(MIN_DEPLOYMENT_TARGET)" go build -v -o $(EXEPATH) -ldflags="-s -w -X main.Version=$(VERSION) -X main.Build=$(BUILD)"

create_darwin_app:
	# copy the binary to the distribution folder
	cp -f "$(DARWINDISTPATH)/$(EXECUTABLE)" "$(RESOURCESPATH)/$(EXECUTABLE)"

	# create the MacOS app
	./macapp -assets "$(RESOURCESPATH)" -bin $(EXECUTABLE) -icon "$(RESOURCESPATH)/icons/icon.png" -identifier "com.github.andrewdjackson.memsfcr" -name "$(APPNAME)" -o "$(DARWINDISTPATH)"
	# copy the info and entitlement plists into the application structure
	cp -f "$(DISTPATH)/Info.plist" "$(DARWINDISTPATH)/$(APPNAME).app/Contents/Info.plist"
	cp -f "$(DISTPATH)/entitlements.plist" "$(DARWINDISTPATH)/$(APPNAME).app/Contents/entitlements.plist"
	cp -f "$(DISTPATH)/rosco.env" "$(DARWINDISTPATH)/$(APPNAME).app/Contents/MacOS/rosco.env"

sign_app_local:
	# sign with the app
	codesign --force  --deep --verify --verbose=4 -s $(DEVID) --timestamp --options runtime "$(DARWINDISTPATH)/$(APPNAME).app/Contents/MacOS/$(EXECUTABLE)" "$(DARWINDISTPATH)/$(APPNAME).app"
	# build and sign installer PKG
	#productbuild --component $(DARWINDISTPATH)/$(APPNAME).app /Applications --sign $(LOCAL_INSTID) --product $(DARWINDISTPATH)/$(APPNAME).app/Contents/Info.plist $(DARWINDISTPATH)/$(APPNAME).pkg
	# sign the PKG with the entitlements
	#codesign --force  --deep --verify --verbose=4 -s $(LOCAL_DISTID) --timestamp --entitlements "$(DARWINDISTPATH)/$(APPNAME).app/Contents/entitlements.plist" --options runtime "$(DARWINDISTPATH)/$(APPNAME).app/Contents/MacOS/$(EXECUTABLE)" "$(DARWINDISTPATH)/$(APPNAME).pkg"

sign_app_store:
	# sign with the app
	codesign --force  --deep --verify --verbose=4 -s $(STORE_DISTID) --timestamp --entitlements "$(DARWINDISTPATH)/$(APPNAME).app/Contents/entitlements.plist" --options runtime "$(DARWINDISTPATH)/$(APPNAME).app/Contents/MacOS/$(EXECUTABLE)" "$(DARWINDISTPATH)/$(APPNAME).app"
	# build and sign the installer PKG
	productbuild --component $(DARWINDISTPATH)/$(APPNAME).app /Applications --sign $(STORE_INSTID) --product $(DARWINDISTPATH)/$(APPNAME).app/Contents/Info.plist $(DARWINDISTPATH)/$(APPNAME).pkg
	# sign the PKG with the entitlements
	codesign --force  --deep --verify --verbose=4 -s $(STORE_DISTID) --timestamp --entitlements "$(DARWINDISTPATH)/$(APPNAME).app/Contents/entitlements.plist" --options runtime "$(DARWINDISTPATH)/$(APPNAME).app/Contents/MacOS/$(EXECUTABLE)" "$(DARWINDISTPATH)/$(APPNAME).pkg"


package_local_app:
	rm -f $(DARWINDISTPATH)/$(APPNAME).dmg
	# create a DMG for local distributions
	appdmg $(DISTPATH)/dmgspec.json $(DARWINDISTPATH)/$(APPNAME).dmg

notarize_local_package:
	# notarize the DMG
	xcrun notarytool submit $(DARWINDISTPATH)/$(APPNAME).dmg --wait --keychain-profile "APPLEDEV"
	#xcrun altool --notarize-app -f $(DARWINDISTPATH)/$(APPNAME).dmg --primary-bundle-id "com.github.andrewdjackson.memsfcr" -u $(APPLEDEVUSR) -p $(APPLEDEVPWD)

	# notarize the PKG
	xcrun notarytool submit $(DARWINDISTPATH)/$(APPNAME).pkg --wait --keychain-profile "APPLEDEV"
	#xcrun altool --notarize-app -f $(DARWINDISTPATH)/$(APPNAME).pkg --primary-bundle-id "com.github.andrewdjackson.memsfcr" -u $(APPLEDEVUSR) -p $(APPLEDEVPWD)

	#
	# if successful staple the app for offline installation
	xcrun stapler staple $(DARWINDISTPATH)/$(APPNAME).app & xcrun stapler staple $(DARWINDISTPATH)/$(APPNAME).pkg & xcrun stapler staple $(DARWINDISTPATH)/$(APPNAME).dmg

upload_app_to_store:
	xcrun altool --upload-app -f $(DARWINDISTPATH)/$(APPNAME).pkg --primary-bundle-id "com.github.andrewdjackson.memsfcr" -u $(APPLEDEVUSR) -p $(APPLEDEVPWD)


clean: ## Remove previous build
	rm -f $(EXEPATH)
	rm -f $(DARWINDISTPATH)/Applications
	rm -fr $(DARWINDISTPATH)/*
	rm -fr $(DISTPATH)/$(APPNAME).dmg
	rm -fr $(DISTPATH)/$(APPNAME).pkg


help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
