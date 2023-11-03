@echo off
go mod tidy

set exeFile=%1
cd %exeFile%

set GOARCH=amd64
set GOOS=linux
go build

set GOARCH=amd64
set GOOS=windows
go build

set srcPath=%CD%
set linuxPath=%CD%\Package\Linux
set windowsPath=%CD%\Package\Windows

IF exist %linuxPath% rd /S /Q %linuxPath%
mkdir %linuxPath%

IF exist %windowsPath% rd /S /Q %windowsPath%
mkdir %windowsPath%

copy /y %srcPath%\%exeFile%.exe %windowsPath%\%exeFile%.exe
copy /y %srcPath%\config.json %windowsPath%\config.json

copy /y %srcPath%\%exeFile% %linuxPath%\%exeFile%
copy /y %srcPath%\config.json  %linuxPath%\config.json

del %srcPath%\%exeFile%

@echo %date%%time%
cd ..