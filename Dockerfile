FROM golang as builder

WORKDIR /build
COPY . . 

RUN make build

FROM gcr.io/distroless/base:nonroot

COPY --from=builder /build/bin/conf-sync  /usr/local/bin/conf-sync 

ENTRYPOINT [ "/usr/local/bin/conf-sync" ]