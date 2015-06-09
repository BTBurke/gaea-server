#/bin/bash

VERSION="1.4.2"
ARCH="amd64"
OS="linux"

wget https://storage.googleapis.com/golang/go$VERSION.$OS-$ARCH.tar.gz
tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
rm go$VERSION.$OS-$ARCH.tar.gz

mkdir -p ~/workspace/src/github.com/BTBurke
chown -R ubuntu:ubuntu ~/workspace
echo "PATH=/usr/local/go/bin:/home/ubuntu/workspace/bin:$PATH" >> ~/.bashrc
echo "GOPATH=~/workspace" >> ~/.bashrc
echo "GOROOT=/usr/local/go" >> ~/.bashrc
bash
echo "Open a new terminal or run bash for changes to take effect..."