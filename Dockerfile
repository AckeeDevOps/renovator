FROM golang:1.12.6-alpine3.10 AS builder
WORKDIR /root
COPY . .
RUN apk add git ca-certificates --no-cache
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o renovator .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /root/renovator .
ENV INSECURE="false"
ENV DEBUG="false"
ENV TTL_THRESHOLD="15206400"
ENV TTL_INCREMENT="5184000"
ENV VAULT_ADDRESS=""
ENV CONFIG_FILE_PATH="/etc/renovator/config.json"
ENV SLACK_WEBHOOK_URL=""
CMD ["./renovator"]
