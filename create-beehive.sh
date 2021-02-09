#!/bin/bash

# using a self-signed certificate for all public tls endpoints for now.
openssl req -newkey rsa:2048 -nodes -keyout key.pem -x509 -days 365 -out cert.pem -subj "/CN=beehive"
cp cert.pem cacert.pem

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
            "configure": ".*",
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

listeners.ssl.default = 5671
listeners.tcp = none
ssl_options.cacertfile           = /etc/rabbitmq/cacert.pem
ssl_options.certfile             = /etc/rabbitmq/cert.pem
ssl_options.keyfile              = /etc/rabbitmq/key.pem
ssl_options.fail_if_no_peer_cert = false
# ssl_options.verify               = verify_peer

management.ssl.port       = 15671
management.ssl.cacertfile = /etc/rabbitmq/cacert.pem
management.ssl.certfile   = /etc/rabbitmq/cert.pem
management.ssl.keyfile    = /etc/rabbitmq/key.pem
EOF

# define config and secrets for rabbitmq
kubectl delete secret rabbitmq-config-secret
kubectl create secret generic rabbitmq-config-secret \
    --from-file=rabbitmq.conf=rabbitmq.conf \
    --from-file=definitions.json=definitions.json \
    --from-file=cacert.pem=cacert.pem \
    --from-file=cert.pem=cert.pem \
    --from-file=key.pem=key.pem

# define rabbitmq credentials for beehive services
kubectl delete secret beehive-service-secret
kubectl create secret generic beehive-service-secret \
    --from-file=cacert.pem=cacert.pem \
    --from-literal=username="service" \
    --from-literal=password="$service_password"

# clean up all configs and secrets now that they should be in kubernetes
rm -f rabbitmq.conf definitions.json cacert.pem cert.pem key.pem

# ensure that rabbitmq is recreated with these credentials
kubectl delete -f rabbitmq.yaml
kubectl create -f rabbitmq.yaml
