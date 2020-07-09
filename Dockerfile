# Build steps
FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ercole-services .

# Run steps
FROM alpine:latest
#RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /app/ercole-services .
# Default config file, at the bottom of the hierarchy
COPY --from=builder /app/config.toml /opt/ercole/config.toml

RUN mkdir /app/distributed_files/
COPY --from=builder /app/resources/ /app/resources/

# Mount this volume to add your config file as ercole.toml and to insert ssh keys
VOLUME [ "/etc/ercole" ]

EXPOSE 11111 11112 11113 11114 11115 11116
CMD ["./ercole-services", "serve"]