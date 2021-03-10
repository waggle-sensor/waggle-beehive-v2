# Installing Beehive in a Kubernetes cluster

The following instructions use [k3s](https://k3s.io), but most of the steps should apply to any Kubernetes cluster.

1. Install and start [k3s](https://k3s.io).

2. Clone the [Beehive repo](https://github.com/waggle-sensor/waggle-beehive-v2).

```sh
git clone https://github.com/waggle-sensor/waggle-beehive-v2
```

3. Install Beehive to the Kubernetes cluster.

```sh
cd waggle-beehive-v2
./create-beehive.sh
```

4. Confirm Beehive is running.

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
