#!/bin/bash
#
# allow specifying different destination directory
DIR="${DIR:-"$HOME/.local/bin"}"

# map different architecture variations to the available binaries
ARCH=$(uname | tr [:upper:] [:lower:])

# prepare the download URL
GITHUB_LATEST_VERSION=$(curl -L -s -H 'Accept: application/json' https://github.com/denetpro/node/releases/latest | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
GITHUB_FILE="denode-$ARCH.zip"
GITHUB_URL="https://github.com/denetpro/node/releases/download/${GITHUB_LATEST_VERSION}/${GITHUB_FILE}"

# install/update the local binary
curl -L -o denode.zip $GITHUB_URL
unzip -a denode.zip 
mv builds/denode ./denode && rm -rf builds denode.zip
mkdir -p $DIR
install -m 555 denode -t "$DIR"
