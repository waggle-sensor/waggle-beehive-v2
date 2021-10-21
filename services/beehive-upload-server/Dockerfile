
# docker build -t waggle/beehive-upload-server .

FROM alpine:3.12
RUN apk add --no-cache openssh-server rsync curl jq
RUN mkdir -p /run/sshd
COPY entrypoint.sh /entrypoint.sh
CMD [ "/entrypoint.sh" ]
