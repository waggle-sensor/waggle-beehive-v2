#!/bin/bash

for svc in beehive-data-logger beehive-message-generator; do
    docker build -t "waggle/$svc" "$svc" && docker push "waggle/$svc"
done
