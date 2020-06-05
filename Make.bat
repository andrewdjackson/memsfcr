set EXECUTABLE=memsfcr
set APPNAME=MemsFCR
set DISTPATH=dist
set RESOURCESPATH=resources
set WINDOWSDISTPATH=dist/windows
set WINDOWS=%WINDOWSDISTPATH%/%EXECUTABLE%.exe
set VERSION=1.0.1
set VCNEXE="C:\Program Files\CodeNotary\vcn"

go build -i -v -o %WINDOWS% -ldflags="-H windowsgui -s -w -X main.version=%VERSION%"
%VCNEXE% notarize %WINDOWS%
%VCNEXE% login
%VCNEXE% authenticate %WINDOWS%
