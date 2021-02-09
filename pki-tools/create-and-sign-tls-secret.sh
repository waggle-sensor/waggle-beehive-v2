#!/bin/bash -e

secret_name="$1"

if [ -z "$secret_name" ]; then
    echo "please provide a secret name"
    exit 1
fi

keyfile="$secret_name.key.pem"
csrfile="$secret_name.csr.pem"
certfile="$secret_name.cert.pem"

openssl genrsa -out "$keyfile" 2048

openssl req -new -key "$keyfile" -out "$csrfile" -config csr.conf

openssl x509 -req -in "$csrfile" -CA cacert.pem -CAkey cakey.pem \
    -CAcreateserial -out "$certfile" -days 365 \
    -extensions v3_ext -extfile csr.conf

# define rabbitmq credentials for beehive services
if kubectl get secret "$secret_name" &> /dev/null; then
    kubectl delete secret "$secret_name"
fi

kubectl create secret generic "$secret_name" \
    --from-file=cacert.pem="cacert.pem" \
    --from-file=cert.pem="$certfile" \
    --from-file=key.pem="$keyfile"
