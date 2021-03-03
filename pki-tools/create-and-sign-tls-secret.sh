#!/bin/bash -e

# ensure working in pki-tools dir
cd $(dirname $0)

cn="$1"
secret_name="$2"

if [ -z "$cn" ] || [ -z "$secret_name" ]; then
    echo "usage: $0 cn secret-name"
    exit 1
fi

keyfile="$cn.key.pem"
csrfile="$cn.csr.pem"
certfile="$cn.cert.pem"

openssl genrsa -out "$keyfile" 2048
openssl req -new -key "$keyfile" -out "$csrfile" -config csr.conf -subj "/CN=$cn"
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

# clean up files which should now be in kubernetes
rm -f "$keyfile" "$csrfile" "$certfile"
