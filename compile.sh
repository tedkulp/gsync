#!/bin/sh

export CGO_ENABLED=0
for GOOS in darwin linux; do
  for GOARCH in 386 amd64; do
    echo "Building $GOOS-$GOARCH"
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    go build -o /go/bin/gsync-$GOOS-$GOARCH -a -ldflags "-extldflags '-static'"
  done
done
