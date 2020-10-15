#!/bin/sh

OS=$(uname -s)
ARCH=$(uname -m)
if hash protoc 2>/dev/null; then
    echo "have protoc"
    exit 0;
else
    if [ "$OS" = "Darwin" ]; then 
      curl -o protoc.zip -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-osx-${ARCH}.zip
    else
      curl -o protoc.zip -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-${OS}-${ARCH}.zip
    fi;

    unzip -o protoc.zip -d /usr/local bin/protoc
    unzip -o protoc.zip -d /usr/local 'include/*'
    rm -f protoc.zip
fi
