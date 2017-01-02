FROM alpine:3.4

RUN apk --no-cache --update add ca-certificates && \
    update-ca-certificates

WORKDIR /root

COPY build/app /root/

CMD ["/root/app"]
