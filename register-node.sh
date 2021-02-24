#!/bin/bash

echo "adding node to rabbitmq"
./add-nodes-to-rabbitmq.py "$1" > /dev/null

echo "adding node to upload server"
./add-nodes-to-upload-server.sh "$1" > /dev/null
