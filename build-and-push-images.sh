#!/bin/bash

for svc in beehive-message-logger beehive-message-generator; do
    docker build -t "waggle/$svc" "$svc" && docker push "waggle/$svc"
done
