#!/bin/bash -e

# ensure working in pki-tools dir
cd $(dirname $0)

cn="$1"
secret_name="$2"

if [ -z "$cn" ] || [ -z "$secret_name" ]; then
    echo "usage: $0 cn secret-name"
    exit 1
fi

ssh_keyfile="ssh-host-key-$cn"

ssh-keygen -C "$cn ssh host key" -N "" -f "$ssh_keyfile"
ssh-keygen \
    -s ca \
    -t rsa-sha2-256 \
    -I "$cn ssh host key" \
    -n "$cn" \
    -V "-5m:+365d" \
    -h \
    "$ssh_keyfile"

# define rabbitmq credentials for beehive services
if kubectl get secret "$secret_name" &> /dev/null; then
    kubectl delete secret "$secret_name"
fi

kubectl create secret generic "$secret_name" \
    --from-file=ca.pub="ca.pub" \
    --from-file=ssh-host-key="$ssh_keyfile" \
    --from-file=ssh-host-key.pub="$ssh_keyfile.pub" \
    --from-file=ssh-host-key-cert.pub="$ssh_keyfile-cert.pub"

# clean up files which should now be in kubernetes
rm -f "$ssh_keyfile" "$ssh_keyfile.pub" "$ssh_keyfile-cert.pub"
