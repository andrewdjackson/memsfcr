set EXECUTABLE=memsfcr
set APPNAME=MemsFCR
set DISTPATH=dist
set RESOURCESPATH=resources
set WINDOWSDISTPATH=dist/windows
set WINDOWS=%WINDOWSDISTPATH%/%EXECUTABLE%.exe
set VERSION=2.6.2

go build -o %WINDOWS% -ldflags="-H windowsgui -s -w -X main.version=%VERSION%"
