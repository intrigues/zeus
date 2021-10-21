#!/bin/bash
GOBIN=`go env GOPATH`/bin

# Installing gin package
if [[ -t $GOBIN/gin ]]
then
    echo "gin present on the system"
else
    go get github.com/codegangsta/gin
fi
$GOBIN/gin --build cmd/web/ --bin app --appPort 3000 --port 8080
# go build -o app cmd/web/*.go && ./app