#!/bin/sh

set -e -x

# Setup the gopath based on current directory.
export GOPATH=$PWD

ls -la
# Now we must move our code from the current directory ./hello-go to $GOPATH/src/github.com/JeffDeCola/hello-go
mkdir -p src/github.com/oswemo
cp -R ./UrlCheck src/github.com/oswemo/.

# All set and everything is in the right place for go
echo "Gopath is: " $GOPATH
echo "pwd is: " $PWD
cd src/github.com/oswemo/UrlCheck
ls -lat

make deps build

