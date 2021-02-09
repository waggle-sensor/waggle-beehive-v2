# PKI Tools

This directory contains some tools for bootstrapping a PKI for _development and testing_. In production,
the Beekeeper should be managing the PKI.

If you are working on Beehive services without any nodes, then you shouldn't have to use this as the
deployment scripts will automatically use these tools.

If you are using an external node, then you can use the `create-and-sign-tls-archive.sh` after deploying
Beehive to create a simple .tar.gz archive you can copy to your node. It contains:
* `cacert.pem`
* `cert.pem`
* `key.pem`
