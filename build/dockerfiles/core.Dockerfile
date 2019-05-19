FROM golang:latest as builder
ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN mkdir colors-core-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o colors-core-service/main -i cmd/core/core_server.go
RUN cp config.json colors-core-service/config.json


FROM alpine:latest
ENV COLORS_SERVICE_USE_MODE=IN_DOCKER_NET
ENV COLORS_DB=IN_DOCKER_NET
ENV COLORS_CONFIG_PATH="/root/colors-core-service"

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/colors-core-service ./colors-core-service
EXPOSE 8002
ENTRYPOINT ["/root/colors-core-service/main"]
