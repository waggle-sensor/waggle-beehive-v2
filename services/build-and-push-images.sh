#!/bin/bash -e

for dockerfilePath in $(find . -name Dockerfile); do
    dir=$(dirname "$dockerfilePath")
    name=$(basename "$dir")
    docker build -t "waggle/$name" "$dir"
    docker push "waggle/$name"
done
