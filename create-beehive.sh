#!/bin/bash

files_exist() {
    for f in $*; do
        if ! test -e "$f"; then
            return 1
        fi
    done
}

# using a self-signed certificate for all public tls endpoints for now.
if ! files_exist cacert.pem cert.pem key.pem; then
    openssl req -newkey rsa:2048 -nodes -keyout key.pem -x509 -days 365 -out cert.pem -subj "/CN=beehive"
    cp cert.pem cacert.pem
fi

# TOTO manage these better
admin_password=admin
service_password=$(openssl rand -base64 16)

# generate rabbitmq definitions file. this only creates / updates the config of things
# in this definitions file - other preexisting resources are not affected.
cat <<EOF > definitions.json
{
    "users": [
        {
            "name": "admin",
            "password": "$admin_password",
            "tags": "administrator",
            "limits": {}
        },
        {
            "name": "service",
            "password": "$service_password",
            "tags": "",
            "limits": {}
        }
    ],
    "vhosts": [
        {
            "name": "/"
        }
    ],
    "permissions": [
        {
            "user": "admin",
            "vhost": "/",
            "configure": ".*",
            "write": ".*",
            "read": ".*"
        },
        {
            "user": "service",
            "vhost": "/",
            "configure": "^$",
            "write": ".*",
            "read": ".*"
        }
    ],
    "topic_permissions": [],
    "parameters": [],
    "policies": [],
    "queues": [],
    "exchanges": [
        {
            "name": "waggle.msg",
            "vhost": "/",
            "type": "topic",
            "durable": true,
            "auto_delete": false,
            "internal": false,
            "arguments": {}
        }
    ],
    "bindings": []
}
EOF

cat <<EOF > rabbitmq.conf
load_definitions = /etc/rabbitmq/definitions.json

ssl_options.cacertfile           = /etc/rabbitmq/cacert.pem
ssl_options.certfile             = /etc/rabbitmq/cert.pem
ssl_options.keyfile              = /etc/rabbitmq/key.pem
ssl_options.fail_if_no_peer_cert = false
# ssl_options.verify               = verify_peer

management.ssl.cacertfile = /etc/rabbitmq/cacert.pem
management.ssl.certfile   = /etc/rabbitmq/cert.pem
management.ssl.keyfile    = /etc/rabbitmq/key.pem
EOF

# define config and secrets for rabbitmq
kubectl create secret generic rabbitmq-config-secret \
    --from-file=rabbitmq.conf=rabbitmq.conf \
    --from-file=definitions.json=definitions.json \
    --from-file=cacert.pem=cacert.pem \
    --from-file=cert.pem=cert.pem \
    --from-file=key.pem=key.pem

# define rabbitmq credentials for beehive services
kubectl create secret generic rabbitmq-service-secret \
    --from-literal=RABBITMQ_USERNAME="service" \
    --from-literal=RABBITMQ_PASSWORD="$service_password"

kubectl create -f rabbitmq.yaml
