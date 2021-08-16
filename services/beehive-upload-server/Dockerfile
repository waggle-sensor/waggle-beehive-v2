FROM alpine:3.12
RUN apk add --no-cache openssh-server rsync
RUN mkdir -p /run/sshd
COPY entrypoint.sh /entrypoint.sh
CMD [ "/entrypoint.sh" ]
