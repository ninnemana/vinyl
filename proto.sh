#! /bin/bash

PKG="$GOPATH/src/github.com/ninnemana/vinyl"
mkdir -p $PKG/openapi $PKG/sdks/java $PKG/sdks/csharp

prototool generate prototool.yaml --debug

npm install -g redoc-cli
redoc-cli bundle \
$PKG/openapi/vinyl.swagger.json \
-o="$PKG/openapi/index.html" \
--title "Vinyl Registry API"
