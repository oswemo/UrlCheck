#!/bin/sh

set -e -x

# Setup the gopath based on current directory.
export GOPATH=$PWD

ls -la
# Now we must move our code from the current directory ./hello-go to $GOPATH/src/github.com/JeffDeCola/hello-go
mkdir -p src/github.com/oswemo
cp -R ./urlCheck src/github.com/oswemo/.

# All set and everything is in the right place for go
echo "Gopath is: " $GOPATH
echo "pwd is: " $PWD
cd src/github.com/oswemo/urlCheck
ls -lat

# RUN unit_tests and it shows the percentage coverage
# print to stdout and file using tee
go test -cover ./... | tee test_coverage.txt
# add some whitespace to the begining of each line
sed -i -e 's/^/     /' test_coverage.txt

# Move to coverage-results directory.
mv test_coverage.txt $GOPATH/coverage-results/.