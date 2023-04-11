# syntax=docker/dockerfile:1
FROM golang:1.20.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o proxy .


FROM scratch

WORKDIR /app

ARG API_VERSION
ENV VERSION ${API_VERSION}

# copy ca certificate.crt from the build stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/proxy /app

EXPOSE 3030

# Run
ENTRYPOINT ["/app/proxy"]