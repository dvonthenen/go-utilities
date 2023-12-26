#!/bin/bash

rm -rf ./src
rm -rf ./dst

mkdir ./src
cp -r ./src-ORG/* ./src

mkdir ./dst
cp -r ./dst-ORG/* ./dst

rm -f ./diff-directory
go build .

if [[ -z "$(command -v ./diff-directory)" ]]; then
    echo "diff-directory did not compile"
else
    # ./diff-directory -src=./src -dst=./dst -skipsrc -logging=7
    # ./diff-directory -src=./src -dst=./dst -dryrun -skipsrc
    # ./diff-directory -src=./src -dst=./dst -dryrun -logging=7
    # ./diff-directory -src=./src -dst=./dst -dryrun -logging=2
    # ./diff-directory -src=./src -dst=./dst -skipsrc
    ./diff-directory -src=./src -dst=./dst
fi
