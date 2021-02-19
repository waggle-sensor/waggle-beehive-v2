#!/bin/bash -e

basedir="$1"

if [ -z "$basedir" ]; then
    basedir=.
fi

for dockerfilePath in $(find "$basedir" -name Dockerfile); do
    dir=$(dirname "$dockerfilePath")
    name=$(basename "$dir")
    docker build -t "waggle/$name" "$dir"
    docker push "waggle/$name"
done
