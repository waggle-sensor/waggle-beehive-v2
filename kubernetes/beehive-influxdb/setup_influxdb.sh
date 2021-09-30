
#!/bin/bash

NAMESPACE=$1
PWD=$2

kubectl exec svc/beehive-influxdb -n ${NAMESPACE} -- influx setup \
        --org waggle \
        --bucket waggle \
        --username waggle \
        --password ${PWD} \
        --force