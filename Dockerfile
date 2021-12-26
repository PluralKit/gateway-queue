FROM alpine:latest as builder

RUN apk add go

WORKDIR /build
COPY . /build/

RUN go build .

FROM alpine:latest

COPY --from=builder /build/gateway-queue /bin/gateway-queue

ENTRYPOINT [ "/bin/gateway-queue" ]