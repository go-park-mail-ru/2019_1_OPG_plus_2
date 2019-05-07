FROM golang:latest as builder
ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN mkdir colors-cookie-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o colors-cookie-service/main -i cmd/chat/chat_server.go
RUN cp config.json colors-cookie-service/config.json


FROM alpine:latest
ENV COLORS_SERVICE_USE_MODE=IN_DOCKER_NET
ENV COLORS_CONFIG_PATH="/root/colors-auth-service"

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/colors-cookie-service ./colors-cookie-service
EXPOSE 50243
ENTRYPOINT ["/root/colors-cookie-service/main"]
