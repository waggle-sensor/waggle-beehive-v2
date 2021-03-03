#!/bin/bash

# WARNING the openssl on macos isn't exactly compatible with our flags / configs.
# it will probably break or lead to weird problems.

cd $(dirname $0)

echo "creating tls and ssh ca"
pki-tools/create-ca.sh

echo "deploying rabbitmq"
pki-tools/create-and-sign-tls-secret.sh beehive-rabbitmq beehive-rabbitmq-tls-secret
kubectl apply -f kubernetes/beehive-rabbitmq.yaml

echo "deploying message logger"
pki-tools/create-and-sign-tls-secret.sh beehive-message-logger beehive-message-logger-tls-secret
kubectl apply -f kubernetes/beehive-message-logger.yaml

echo "deploying upload server"
pki-tools/create-and-sign-ssh-host-key-secret.sh beehive-upload-server beehive-upload-server-ssh-secret
kubectl apply -f kubernetes/beehive-upload-server.yaml

echo "deploying influxdb"
kubectl apply -f kubernetes/beehive-influxdb.yaml

setup_influxdb() {
    kubectl exec svc/beehive-influxdb -- influx setup \
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
    kubectl exec svc/beehive-influxdb -- influx auth create \
        --user waggle \
        --org waggle \
        --hide-headers $* | awk '{print $2}'
}

echo "generating token for data loader"
token=$(generate_influxdb_token --write-buckets)
kubectl create secret generic beehive-influxdb-loader-secret \
    --from-literal=token="$token"
pki-tools/create-and-sign-tls-secret.sh beehive-influxdb-loader beehive-influxdb-loader-tls-secret
kubectl apply -f kubernetes/beehive-influxdb-loader.yaml

echo "generating token for data api"
token=$(generate_influxdb_token --read-buckets)
kubectl create secret generic beehive-influxdb-data-api-secret \
    --from-literal=token="$token"
kubectl apply -f kubernetes/beehive-influxdb-data-api.yaml

echo "creating credentials for but will not deploy message generator"
pki-tools/create-and-sign-tls-secret.sh beehive-message-generator beehive-message-generator-tls-secret
# kubectl apply -f kubernetes/beehive-message-generator.yaml
