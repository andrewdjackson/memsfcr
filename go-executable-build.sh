#!/usr/bin/env bash

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi
package_split=(${package//\// })
package_name=${package_split[-1]}

platforms=("windows/amd64" "darwin/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi


#CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -i -v -o $(WINDOWS) -ldflags="-H windowsgui -s -w -X main.version=$(VERSION)"

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi

    # build MacOS app
    if [ $GOOS = "darwin" ]; then
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
    fi
done