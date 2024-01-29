FROM alpine:latest as root-certs
RUN apk add -U --no-cache ca-certificates
RUN addgroup -g 1001 app
RUN adduser app -u 1001 -D -G app /home/app

FROM golang:latest as builder
WORKDIR /app
COPY --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin ./cmd/...

FROM scratch as final
COPY --from=root-certs /etc/passwd /etc/passwd
COPY --from=root-certs /etc/group /etc/group
COPY --chown=1001:10001 --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --chown=1001:10001 --from=builder /app/bin/ /app
USER app
ENTRYPOINT [ "/app" ]