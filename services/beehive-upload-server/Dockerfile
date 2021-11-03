
# docker build -t waggle/beehive-upload-server .

#FROM alpine:3.12
FROM golang:1.17-alpine
RUN apk add --no-cache openssh-server rsync curl jq supervisor
RUN mkdir -p /run/sshd
COPY entrypoint.sh /entrypoint.sh
COPY supervisord.conf /etc/supervisor/conf.d/

COPY user-service /go/src/app

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

#CMD [ "/entrypoint.sh" ]
#CMD ["/go/bin/user-service"]
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]