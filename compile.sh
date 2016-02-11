#!/bin/bash


docker run --rm -it -v "$GOPATH":/go -w /go/src/github.com/kildevaeld/nodester golang:1.4.2-cross sh -c '
for GOOS in darwin linux; do
  for GOARCH in 386 amd64; do
    echo "Building $GOOS-$GOARCH"
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    CGO_ENABLED=0 go build -a -installsuffix cgo -o bin/nodester-$GOOS-$GOARCH
  done
done
'

