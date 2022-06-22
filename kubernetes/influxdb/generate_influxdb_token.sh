#!/bin/bash

set -e

# example: generate_influxdb_token.sh default --write-buckets
NAMESPACE=$1
PERMISSION=$2

if [ "$#" -ne 2 ]; then
     echo "usage: generate_influxdb_token.sh <namespace> <permission_arg>"
     echo "example: generate_influxdb_token.sh default --write-buckets"
     exit 1
fi

# token is emitted in second column
kubectl exec svc/beehive-influxdb -n ${NAMESPACE} -- influx auth create \
    --user waggle \
    --org waggle \
    --hide-headers ${PERMISSION} | awk '{print $2}'


#echo "generating token for data loader"
#echo "token="$(generate_influxdb_token ${PERMISSION})



