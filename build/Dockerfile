FROM golang:1.12 as builder

COPY / /go/src/example-server
WORKDIR /go/src/example-server
RUN make example-server

FROM alpine:3.9
COPY --from=builder /go/src/example-server/bin/example-server /bin/example-server
CMD ["/bin/example-server"]