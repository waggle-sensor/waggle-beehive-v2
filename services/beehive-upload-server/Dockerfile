FROM alpine:3.12
RUN apk add --no-cache openssh-server rsync
RUN mkdir -p /run/sshd
COPY sshd_config /etc/ssh/sshd_config
CMD [ "/usr/sbin/sshd", "-D", "-e" ]
