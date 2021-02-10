#!/bin/bash -e

# ensure working in pki-tools dir
cd $(dirname $0)

cn="$1"

if [ -z "$cn" ]; then
    echo "usage: $0 cn archive-name.tar.gz"
    exit 1
fi

keyfile="key.pem"
csrfile="csr.pem"
certfile="cert.pem"

openssl genrsa -out "$keyfile" 2048

openssl req -new -key "$keyfile" -out "$csrfile" -config csr.conf -subj="/CN=$cn"

openssl x509 -req -in "$csrfile" -CA cacert.pem -CAkey cakey.pem \
    -CAcreateserial -out "$certfile" -days 365 \
    -extensions v3_ext -extfile csr.conf

tar -czf "$cn.tar.gz" "cacert.pem" "$keyfile" "$certfile"

# clean up files which should now be in archive
rm -f "$keyfile" "$csrfile" "$certfile"
