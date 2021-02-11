#!/bin/bash

cd $(dirname $0)

echo "creating tls and ssh ca"
pki-tools/create-ca.sh

echo "deploying rabbitmq"
if kubectl get secret rabbitmq-config-secret &> /dev/null; then
    kubectl delete secret rabbitmq-config-secret
fi

kubectl create secret generic rabbitmq-config-secret \
    --from-file=rabbitmq.conf=config/rabbitmq/rabbitmq.conf \
    --from-file=enabled_plugins=config/rabbitmq/enabled_plugins \
    --from-file=definitions.json=config/rabbitmq/definitions.json

pki-tools/create-and-sign-tls-secret.sh rabbitmq rabbitmq-tls-secret
kubectl apply -f kubernetes/rabbitmq.yaml

echo "deploying message logger"
pki-tools/create-and-sign-tls-secret.sh message-logger message-logger-tls-secret
kubectl apply -f kubernetes/message-logger.yaml

echo "deploying upload server"
pki-tools/create-and-sign-ssh-host-key-secret.sh beehive-upload-server upload-server-ssh-host-key-secret
kubectl apply -f kubernetes/upload-server.yaml

echo "creating credentials for but will not deploy message generator"
pki-tools/create-and-sign-tls-secret.sh message-generator message-generator-tls-secret
# kubectl apply -f kubernetes/message-generator.yaml
