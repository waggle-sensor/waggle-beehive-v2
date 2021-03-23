#!/bin/bash

# WARNING the openssl on macos isn't exactly compatible with our flags / configs.
# it will probably break or lead to weird problems.

cd $(dirname $0)

echo "deploying rabbitmq"
kubectl apply -f kubernetes/beehive-rabbitmq.yaml

echo "deploying message logger"
./update-rabbitmq-auth.sh beehive-message-logger-auth beehive-message-logger '.*' '.*' '.*'
kubectl apply -f kubernetes/beehive-message-logger.yaml

echo "deploying upload server"
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
kubectl create secret generic beehive-influxdb-loader-influxdb-token \
    --from-literal=token="$token"
./update-rabbitmq-auth.sh beehive-influxdb-loader-auth beehive-influxdb-loader '.*' '.*' '.*'
kubectl apply -f kubernetes/beehive-influxdb-loader.yaml

echo "generating token for data api"
token=$(generate_influxdb_token --read-buckets)
kubectl create secret generic beehive-data-api-influxdb-token \
    --from-literal=token="$token"
kubectl apply -f kubernetes/beehive-data-api.yaml

echo "creating credentials for but will not deploy message generator"
./update-rabbitmq-auth.sh beehive-message-generator-auth beehive-message-generator '^$' '^waggle.msg$' '^$' 'impersonator'
# kubectl apply -f kubernetes/beehive-message-generator.yaml

echo "setting up ingress"
kubectl apply -f kubernetes/beehive-ingress.yaml
