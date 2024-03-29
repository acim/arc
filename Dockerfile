FROM golang:1.17.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-s -w" -o /go/bin/arc ./cmd/arc/.

FROM alpine:3.14.3

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/arc /usr/bin/arc

EXPOSE 3000 3001

USER 65534:65534

ENTRYPOINT ["/usr/bin/arc"]
