#!/bin/bash
set -e

rmqctl() {
    kubectl exec svc/beehive-rabbitmq -n ${namespace} -- rabbitmqctl "$@"
}

namespace="$1"
secretname="$2"
username="$3"
confperm="$4"
writeperm="$5"
readperm="$6"
tags="$7"
password="$(openssl rand -hex 20)"

echo "updating kubernetes config ${secretname}..."
kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: ${secretname}
  namespace: ${namespace}
type: kubernetes.io/basic-auth
stringData:
  username: ${username}
  password: ${password}
EOF

echo "updating rabbitmq user ${username}..."
(
while ! rmqctl authenticate_user "$username" "$password"; do
    while ! (rmqctl add_user "$username" "$password" || rmqctl change_password "$username" "$password"); do
      sleep 3
    done
done

while ! rmqctl set_permissions "$username" "$confperm" "$writeperm" "$readperm"; do
  sleep 3
done

while ! rmqctl set_user_tags "$username" "$tags"; do
  sleep 3
done
) &> /dev/null
echo "done"
