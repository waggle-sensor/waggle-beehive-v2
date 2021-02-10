#!/bin/bash -e

# ensure working in pki-tools dir
cd $(dirname $0)

CN=beekeeper

# generate signing config
cat <<EOF > csr.conf
[ req ]
default_bits = 2048
prompt = no
default_md = sha256
distinguished_name = dn

[ dn ]
C = US
CN = ${CN}

[ v3_ext ]
authorityKeyIdentifier=keyid,issuer:always
basicConstraints=CA:FALSE
keyUsage=keyEncipherment,dataEncipherment
extendedKeyUsage=serverAuth,clientAuth
EOF

# generate tls ca key and cert
openssl genrsa -out cakey.pem 2048
openssl req -x509 -new -nodes -key cakey.pem -subj "/CN=${CN}" -days 3650 -out cacert.pem
chmod 600 cakey.pem

# genrate ssh ca
rm -f ca
ssh-keygen -C "${CN}" -N "" -f ca
