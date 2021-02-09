#!/bin/bash

# ensure working in beehive-v2 directory
cd $(dirname $0)

# create dev/test tls credentials for beehive services
pki-tools/create-ca.sh
pki-tools/create-and-sign-tls-secret.sh rabbitmq rabbitmq-tls-secret
pki-tools/create-and-sign-tls-secret.sh message-logger message-logger-tls-secret
pki-tools/create-and-sign-tls-secret.sh message-generator message-generator-tls-secret

# define config and secrets for rabbitmq
if kubectl get secret rabbitmq-config-secret &> /dev/null; then
    kubectl delete secret rabbitmq-config-secret
fi

kubectl create secret generic rabbitmq-config-secret \
    --from-file=rabbitmq.conf=config/rabbitmq/rabbitmq.conf \
    --from-file=enabled_plugins=config/rabbitmq/enabled_plugins \
    --from-file=definitions.json=config/rabbitmq/definitions.json

# ensure that rabbitmq is recreated with these credentials
kubectl apply -f rabbitmq.yaml
kubectl apply -f message-logger.yaml

# NOTE kubectl exec -i svc/rabbitmq -- rabbitmqctl --timeout 300 import_definitions <<EOF ...
