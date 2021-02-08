#!/bin/bash

openssl req -newkey rsa:2048 -nodes -keyout rabbitmq-key.pem -x509 -days 365 -out rabbitmq-cert.pem -subj "/CN=beehive"

# ensure secret is recreated
kubectl delete secret rabbitmq-tls-secret
kubectl create secret tls rabbitmq-tls-secret \
  --cert=rabbitmq-cert.pem \
  --key=rabbitmq-key.pem
