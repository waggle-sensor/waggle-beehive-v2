#!/bin/sh

# setup users and directories for existing items in /home
echo "Initializing existing users."
for username in $(ls /home); do
    echo "Found $username."
    adduser -D -g "" "$username"
    passwd -u "$username"
    chown -R "$username:$username" "/home/$username"
done

exec /usr/sbin/sshd -D -e
