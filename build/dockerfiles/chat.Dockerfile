FROM golang:latest as builder
ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN mkdir colors-chat-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o colors-chat-service/main -i cmd/chat/chat_server.go
RUN cp config.json colors-chat-service/config.json


FROM alpine:latest
ENV COLORS_SERVICE_USE_MODE=IN_DOCKER_NET
ENV COLORS_DB=IN_DOCKER_NET
ENV COLORS_CONFIG_PATH="/root/colors-chat-service"

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/colors-chat-service ./colors-chat-service
EXPOSE 8003
ENTRYPOINT ["/root/colors-chat-service/main"]
