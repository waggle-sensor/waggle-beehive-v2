# Beehive Upload Server

ssh+rsync based upload server for training data.

## Configuration

This image expects CA signed ssh keys. Please ensure that the following files exist in the expected locations:

```txt
TrustedUserCAKeys /etc/waggle/ca.pub
HostKey /etc/waggle/upload-server-key
HostCertificate /etc/waggle/upload-server-key-cert.pub
```
