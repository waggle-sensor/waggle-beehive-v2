#!/bin/sh

fatal() {
    echo $*
    exit 1
}

if ! test -e "${SSH_CA_PUBKEY}"; then
    fatal "Error: CA public key ${SSH_CA_PUBKEY} does not exist!"
fi
echo "Using CA public key at ${SSH_CA_PUBKEY}."

if ! test -e "${SSH_HOST_KEY}"; then
    fatal "Error: Upload server host key key ${SSH_HOST_KEY} does not exist!"
fi
echo "Using upload server host key at ${SSH_HOST_KEY}."

if ! test -e "${SSH_HOST_CERT}"; then
    fatal "Error: Upload server signed host key key ${SSH_HOST_CERT} does not exist!"
fi
echo "Using upload server signed host key at ${SSH_HOST_CERT}."

# generate sshd_config from env vars
cat > /etc/ssh/sshd_config <<EOF
Port 22
ListenAddress 0.0.0.0
ListenAddress ::

TrustedUserCAKeys ${SSH_CA_PUBKEY}
HostKey ${SSH_HOST_KEY}
HostCertificate ${SSH_HOST_CERT}

LogLevel VERBOSE

LoginGraceTime 60
PermitRootLogin prohibit-password
StrictModes yes
MaxAuthTries 3
MaxSessions 3

PubkeyAuthentication yes
AuthorizedKeysFile	.ssh/authorized_keys
AuthorizedPrincipalsFile none

PasswordAuthentication no
ChallengeResponseAuthentication no

AllowAgentForwarding no
AllowTcpForwarding no
GatewayPorts no
X11Forwarding no
PermitTTY no
PrintMotd no
TCPKeepAlive yes
PermitUserEnvironment no
#Compression delayed
#ClientAliveInterval 0
#ClientAliveCountMax 3
UseDNS no
PidFile /var/run/sshd.pid
#MaxStartups 10:30:100
PermitTunnel no
AcceptEnv LANG LC_*
EOF

# setup users and directories for existing items in /home
echo "Initializing existing users."
for username in $(ls /home); do
    echo "Found $username."
    adduser -D -g "" "$username"
    passwd -u "$username"
    chown -R "$username:$username" "/home/$username"
done

for username in $(curl -s https://api.sagecontinuum.org/api/state | jq -r .data[].id | awk '{print tolower($0)}') ; do
    if [ -e "/home/$username" ] ; then
        echo "/home/$username already exists"
        continue
    fi

    echo "Adding $username."
    adduser -D -g "" "$username"
    passwd -u "$username"
    chown -R "$username:$username" "/home/$username"
done


exec /usr/sbin/sshd -D -e
