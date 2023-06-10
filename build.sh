#!/bin/bash

GOARCH=amd64

macBuildName=tw_darwin_${GOARCH}
winBuildName=tw_windows_${GOARCH}
linuxBuildName=tw_linux_${GOARCH}

# clear build/
rm -rf tw_build
mkdir tw_build

build() {
  # Build
  GOOS=darwin GOARCH=${GOARCH} go build -o $macBuildName translate.go
  GOOS=linux GOARCH=${GOARCH} go build -o $linuxBuildName translate.go
  GOOS=windows GOARCH=${GOARCH} go build -o $winBuildName.exe translate.go

  # Compress
#  upx $macBuildName
#  upx $linuxBuildName
#  upx $winBuildName.exe

  # Move
  mv -f $macBuildName tw_build/$macBuildName
  mv -f $linuxBuildName tw_build/$linuxBuildName
  mv -f $winBuildName.exe tw_build/$winBuildName.exe

  # gzip
  tar -cvf tw_build/${macBuildName}.tar tw_build/${macBuildName} && gzip tw_build/${macBuildName}.tar
  tar -cvf tw_build/${linuxBuildName}.tar tw_build/${linuxBuildName} && gzip tw_build/${linuxBuildName}.tar
  zip -j tw_build/${winBuildName}.zip tw_build/${winBuildName}.exe
}

build
