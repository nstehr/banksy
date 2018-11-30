#!/bin/bash

echo "Downloading goveralls"
go get github.com/mattn/goveralls

# executes the tests in each package, and collects their coverage
echo "mode: set" > acc.out
for Dir in $(go list ./...); 
do
    if [[ ${Dir} != *"/vendor/"* ]]
    then
        returnval=`go test -coverprofile=profile.out $Dir`
        echo ${returnval}
        if [[ ${returnval} != *FAIL* ]]
        then
            if [ -f profile.out ]
            then
                cat profile.out | grep -v "mode: set" >> acc.out 
            fi
        else
            exit 1
        fi
    else
        exit 1
    fi  

done

# upload to coveralls
$GOPATH/bin/goveralls -coverprofile=acc.out -service=travis-ci
  
for GOOS in darwin linux windows; do
  for GOARCH in 386 amd64; do
    echo "Building $GOOS-$GOARCH"
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    go build -o bin/banksy-$GOOS-$GOARCH
  done
done