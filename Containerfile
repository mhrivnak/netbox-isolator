FROM registry.access.redhat.com/ubi9/go-toolset:latest AS builder

COPY . .

RUN go build -o netbox-isolator ./cmd/main.go

FROM registry.access.redhat.com/ubi9/ubi-micro:latest

COPY --from=builder /opt/app-root/src/netbox-isolator /usr/bin/netbox-isolator

USER 1001

EXPOSE 8080

CMD ["/usr/bin/netbox-isolator"]