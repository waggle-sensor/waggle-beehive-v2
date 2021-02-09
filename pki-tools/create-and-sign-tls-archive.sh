#!/bin/bash -e

# ensure working in pki-tools dir
cd $(dirname $0)

cn="$1"

if [ -z "$cn" ]; then
    echo "usage: $0 cn archive-name.tar.gz"
    exit 1
fi

keyfile="$cn.key.pem"
csrfile="$cn.csr.pem"
certfile="$cn.cert.pem"

openssl genrsa -out "$keyfile" 2048

openssl req -new -key "$keyfile" -out "$csrfile" -config csr.conf -subj="/CN=$cn"

openssl x509 -req -in "$csrfile" -CA cacert.pem -CAkey cakey.pem \
    -CAcreateserial -out "$certfile" -days 365 \
    -extensions v3_ext -extfile csr.conf

tar -czf "$cn.tar.gz" cacert.pem cert.pem key.pem

# clean up files which should now be in kubernetes
rm -f "$keyfile" "$csrfile" "$certfile"
