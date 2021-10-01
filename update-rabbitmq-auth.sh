#!/bin/bash
set -e

rmqctl() {
    kubectl exec svc/beehive-rabbitmq -n ${rmq_namespace} -- rabbitmqctl "$@"
}

secret_namespace="$1"
secretname="$2"
rmq_namespace="$3"
username="$4"
confperm="$5"
writeperm="$6"
readperm="$7"
tags="$8"
password="$(openssl rand -hex 20)"

echo "updating kubernetes config ${secretname}..."
kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: ${secretname}
  namespace: ${secret_namespace}
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
