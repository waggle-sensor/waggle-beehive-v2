#!/bin/bash

# ensure working in beehive-v2 directory
cd $(dirname $0)

# create dev/test tls credentials for beehive services
pki-tools/create-ca.sh
pki-tools/create-and-sign-tls-secret.sh rabbitmq rabbitmq-tls-secret
pki-tools/create-and-sign-tls-secret.sh data-logger data-logger-tls-secret

# define rabbitmq config
cat <<EOF > rabbitmq.conf
load_definitions = /etc/rabbitmq/definitions.json

default_vhost = /
default_user = admin
default_pass = admin

default_permissions.configure = .*
default_permissions.read = .*
default_permissions.write = .*

listeners.tcp = none
listeners.ssl.default = 5671
ssl_options.cacertfile           = /etc/tls/cacert.pem
ssl_options.certfile             = /etc/tls/cert.pem
ssl_options.keyfile              = /etc/tls/key.pem
ssl_options.fail_if_no_peer_cert = false
ssl_options.verify               = verify_peer

auth_mechanisms.1 = PLAIN
auth_mechanisms.2 = AMQPLAIN
auth_mechanisms.3 = EXTERNAL

ssl_cert_login_from   = common_name

management.ssl.port       = 15671
management.ssl.cacertfile = /etc/tls/cacert.pem
management.ssl.certfile   = /etc/tls/cert.pem
management.ssl.keyfile    = /etc/tls/key.pem
EOF

cat <<EOF > enabled_plugins
[rabbitmq_prometheus,rabbitmq_management,rabbitmq_management_agent,rabbitmq_auth_mechanism_ssl].
EOF

# no... this is something better to apply later and dynamicaally
# generate rabbitmq definitions file. this only creates / updates the config of things
# in this definitions file - other preexisting resources are not affected.
cat <<EOF > definitions.json
{
    "users": [
        {
            "name": "data-logger",
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
            "user": "data-logger",
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

# define config and secrets for rabbitmq
if kubectl get secret rabbitmq-config-secret &> /dev/null; then
    kubectl delete secret rabbitmq-config-secret
fi

kubectl create secret generic rabbitmq-config-secret \
    --from-file=rabbitmq.conf=rabbitmq.conf \
    --from-file=enabled_plugins=enabled_plugins \
    --from-file=definitions.json=definitions.json

# clean up all configs that should be in kubernetes
rm -f rabbitmq.conf enabled_plugins definitions.json

# ensure that rabbitmq is recreated with these credentials
kubectl delete -f rabbitmq.yaml
kubectl create -f rabbitmq.yaml

# technically, we could move the ca management into k8s and generate service secrets via a service account
