set EXECUTABLE=memsfcr
set APPNAME=MemsFCR
set DISTPATH=dist
set RESOURCESPATH=resources
set WINDOWSDISTPATH=dist/windows
set WINDOWS=%WINDOWSDISTPATH%/%EXECUTABLE%.exe
set VERSION=1.0.1
set VCNEXE="C:\Program Files\CodeNotary\vcn"

go build -i -v -o %WINDOWS% -ldflags="-H windowsgui -s -w -X main.version=%VERSION%"
rem %VCNEXE% notarize %WINDOWS%
rem %VCNEXE% login
rem %VCNEXE% authenticate %WINDOWS%
