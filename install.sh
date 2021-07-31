#!/bin/bash

set -e

if [ "$(uname -m)" != "x86_64" ]; then
	echo "Error: Unsupported architecture $(uname -m). Only x64 binaries are available." 1>&2
	exit 1
fi

if ! command -v tar >/dev/null; then
	echo "Error: tar is required to install." 1>&2
	exit 1
fi

goos=$(uname)

install() {
  file_name=tw_${goos}_amd64
  suffix=.tar.gz
  download_name="${file_name}${suffix}"
  uri=https://github.com/xjh22222228/translate-tw/releases/latest/download/${download_name}

  echo -e "Download ${uri} \n"

  # Remove current pkg
  rm -f "$download_name"

  curl "$uri" -OL --progress --retry 2 2>&1

  if [ $? -ne 0 ]; then
    rm -f "${download_name}"
    echo "Download failed"
    exit 1
  fi

  tar -xvf "${download_name}"

  if [ $? -ne 0 ]; then
    rm -f "${download_name}"
    echo "Installation failed"
    exit 1
  fi

  chmod +x tw_build/$file_name
  rm -f "${download_name}"

  mv -f tw_build/$file_name /usr/local/bin/tw
  rm -rf tw_build
}

if [ $goos = "Darwin" ]; then
  goos=darwin
else
  goos=linux
fi

install


echo -e "\n\033[1;32mTRANSLATE-TW was installed successfully\033[0m"

echo -e "\nRun \"tw -h\" \n"

