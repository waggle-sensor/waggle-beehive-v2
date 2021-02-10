# Beehive V2

*These notes are a work in progress, as Beehive V2 is under rapid development.*

Beehive server provides the data and management pipeline for Waggle nodes. This repo contains
tools and configs to help deploy a Beehive server.

The quickest way to deploy a Beehive server for dev/test is use the provided Vagrantbox:
1. Install [Vagrant](https://www.vagrantup.com).
2. Run the Vagrant box (`vagrant up`) and connect (`vagrant ssh`).
3. Inside of the Vaggrant box, `sudo -s` and `cd /vagrant`
4. Run the `./create-beehive.sh` script.

*These steps are likely to change. The initial plan is to spin up a dev/test Beehive at Vagrant provision time.*