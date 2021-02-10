#!/bin/bash

for nodeID in $*; do
    username="node-$nodeID"

    kubectl exec --stdin deployment/upload-server -- sh -s <<EOF
adduser -D -g "" "$username"
passwd -u "$username"
true
EOF
done
