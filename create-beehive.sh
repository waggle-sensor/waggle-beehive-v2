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

echo "deploying influxdb"
kubectl apply -f kubernetes/influxdb.yaml

setup_influxdb() {
    kubectl exec svc/influxdb -- influx setup \
        --org waggle \
        --bucket waggle \
        --username waggle \
        --password wagglewaggle \
        --force
}

echo "setting up influxdb"
while true; do
    if msg=$(setup_influxdb); then
        break
    fi
    if [[ "$msg" == *"already been setup"* ]]; then
        echo "influxdb already setup. skipping."
        break
    fi
    echo "influxdb not ready... will retry."
    sleep 3
done

generate_influxdb_token() {
    # token is emitted in second column
    kubectl exec svc/influxdb -- influx auth create \
        --user waggle \
        --org waggle \
        --hide-headers $* | awk '{print $2}'
}

echo "generating token for data loader"
token=$(generate_influxdb_token --write-buckets)
kubectl create secret generic influxdb-loader-secret \
    --from-literal=token="$token"
pki-tools/create-and-sign-tls-secret.sh influxdb-loader influxdb-loader-tls-secret
kubectl apply -f kubernetes/influxdb-loader.yaml

echo "generating token for data api"
token=$(generate_influxdb_token --read-buckets)
kubectl create secret generic influxdb-data-api-secret \
    --from-literal=token="$token"
kubectl apply -f kubernetes/influxdb-data-api.yaml

echo "creating credentials for but will not deploy message generator"
pki-tools/create-and-sign-tls-secret.sh message-generator message-generator-tls-secret
# kubectl apply -f kubernetes/message-generator.yaml
