FROM alpine:3.10

RUN apk --no-cache add ca-certificates

ADD gomods /gomods

CMD ["/gomods"]