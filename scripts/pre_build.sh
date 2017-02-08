#!/bin/bash

echo "Getting External Dependencies..."
go get github.com/fatih/color
go get github.com/bwmarrin/discordgo
go get github.com/mattn/go-sqlite3

echo "Creating $GOPATH/src/truebot-2.0"

mkdir -p $GOPATH/src/truebot-2.0

echo "Moving files to $GOPATH/src/truebot-2.0"

mv $CODEBUILD_SRC_DIR/*.go $GOPATH/src/truebot-2.0/

echo "Ready to compile."
