#!/bin/bash

docker build -t waggle/beehive-rabbitmq beehive-rabbitmq && docker push waggle/beehive-rabbitmq
docker build -t waggle/beehive-data-logger beehive-data-logger && docker push waggle/beehive-data-logger
