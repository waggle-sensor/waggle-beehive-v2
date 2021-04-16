# Installing Beehive in a Vagrant development box

This document is primarily intended for developers who are working on Beehive!

Most users should instead [deploy Beehive to a Kubernetes cluster](installing-beehive-in-a-kubernetes-cluster.md).

## Installation Steps

1. Install [Vagrant](https://www.vagrantup.com).
2. Run the Vagrant box (`vagrant up`) and connect (`vagrant ssh`).
3. Inside of the Vaggrant box, `sudo -s` and `cd /vagrant`
4. Run the `./create-beehive.sh` script.

_TODO(sean) Sync notes on setting up Beehive certs / keys._
