# Installing Beehive in a Kubernetes cluster

## Install Kubernetes

The following instructions use [k3s](https://k3s.io), but most of the steps should apply to any Kubernetes cluster.

1. Install [k3s](https://k3s.io).

## Install Beehive Credentials

1. Clone the [Waggle PKI Tools repo](https://github.com/waggle-sensor/waggle-pki-tools)

```sh
git clone https://github.com/waggle-sensor/waggle-pki-tools
```

2. Create credentials for Beehive.

```sh
cd waggle-pki-tools
./create-credentials-for-beehive.sh
```

3. Install credentials in Kubernetes cluster.

```sh
kubectl apply -f credentials/beehive.yaml
```

This will provide everything our cluster needs to authenticate and secure connections.

## Install Beehive

1. Clone the [Beehive repo](https://github.com/waggle-sensor/waggle-beehive-v2).

```sh
git clone https://github.com/waggle-sensor/waggle-beehive-v2
```

2. Install Beehive to the Kubernetes cluster.

```sh
cd waggle-beehive-v2
./create-beehive.sh
```

3. Confirm Beehive is running.

We'll confirm that Beehive is up and running using the following.

```sh
kubectl get pod
```

If everything was installed correctly, we should see the following pods with status `Running`.

```sh
NAME                                         READY   STATUS    RESTARTS   AGE
beehive-rabbitmq-0                           1/1     Running   0          6d4h
beehive-upload-server-99fc4c499-b5gn7        1/1     Running   0          6d3h
beehive-influxdb-0                           1/1     Running   0          6d3h
beehive-message-generator-5d94bdd587-wvtmm   1/1     Running   0          6d3h
beehive-message-logger-5585fc77b9-fxtwh      1/1     Running   0          6h18m
beehive-influxdb-loader-d5f7b856d-wbr4v      1/1     Running   0          6h7m
beehive-data-api-555c968656-dfmcp            1/1     Running   0          4h25m
```

Note: Beehive needs a couple minutes to prepare its databases, so some of these commands may initially fail. Please check again after 2-3 minutes to see if all the pods have stablized.
