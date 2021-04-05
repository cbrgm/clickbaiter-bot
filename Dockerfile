FROM alpine:latest
RUN apk add --update ca-certificates

ADD ./clickbaiter-bot /usr/bin/clickbaiter-bot

ENTRYPOINT ["/usr/bin/clickbaiter-bot"]