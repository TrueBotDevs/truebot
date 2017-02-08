#!/bin/bash

echo "Changing to $GOPATH/src/truebot-2.0 and compiling."

cd $GOPATH/src/truebot-2.0

go build
