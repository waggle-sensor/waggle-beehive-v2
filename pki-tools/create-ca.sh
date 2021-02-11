#!/bin/bash -e

cd $(dirname $0)

CN=beekeeper

echo "generating tls ca"
if test -e cakey.pem; then
    echo "$PWD/cakey.pem already exists - skipping."
else
    openssl genrsa -out cakey.pem 2048
    chmod 600 cakey.pem
    openssl req -x509 -new -nodes -key cakey.pem -subj "/CN=${CN}" -days 3650 -out cacert.pem
fi

echo "generating ssh ca"
if test -e ca; then
    echo "$PWD/ca already exists - skipping."
else
    ssh-keygen -C "${CN}" -N "" -f ca
fi
