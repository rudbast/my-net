FROM golang:1.11.0-stretch AS builder

WORKDIR /my-net
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o app

FROM alpine:latest AS runtime

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /root
COPY --from=builder /my-net/app .
COPY --from=builder /my-net/files /
CMD ["./app"]
