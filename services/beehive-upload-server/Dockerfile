
# docker build -t waggle/beehive-upload-server .

# testing:
# docker build -t waggle/beehive-upload-server . && docker rm -f upload-server && docker run -d --name upload-server -p 8080:80  waggle/beehive-upload-server
# curl -X POST localhost:8080/user/1234567890abcdff

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

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]