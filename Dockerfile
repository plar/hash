FROM golang:1.14 as BUILDER

ARG VERSION

WORKDIR /build
COPY . .

RUN make test
RUN CGO_ENABLED=0 make service

FROM alpine:3

EXPOSE 8080

WORKDIR /app

COPY --from=BUILDER /build/_bin/hashsvc .

ENTRYPOINT ["./hashsvc"]
