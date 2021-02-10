#!/bin/bash

# ensure working in beehive-v2 directory
cd $(dirname $0)

# create dev/test tls credentials for beehive services
pki-tools/create-ca.sh

# deploy rabbitmq
if kubectl get secret rabbitmq-config-secret &> /dev/null; then
    kubectl delete secret rabbitmq-config-secret
fi

kubectl create secret generic rabbitmq-config-secret \
    --from-file=rabbitmq.conf=config/rabbitmq/rabbitmq.conf \
    --from-file=enabled_plugins=config/rabbitmq/enabled_plugins \
    --from-file=definitions.json=config/rabbitmq/definitions.json

pki-tools/create-and-sign-tls-secret.sh rabbitmq rabbitmq-tls-secret
kubectl apply -f kubernetes/rabbitmq.yaml

# deploy message logger
pki-tools/create-and-sign-tls-secret.sh message-logger message-logger-tls-secret
kubectl apply -f kubernetes/message-logger.yaml

# deploy upload server
pki-tools/create-and-sign-ssh-host-key-secret.sh beehive-upload-server upload-server-ssh-host-key-secret
kubectl apply -f kubernetes/upload-server.yaml

# create credentials for but don't deploy message generator
pki-tools/create-and-sign-tls-secret.sh message-generator message-generator-tls-secret
# kubectl apply -f kubernetes/message-generator.yaml
