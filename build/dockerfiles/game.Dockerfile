FROM golang:latest as builder
ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN mkdir colors-game-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o colors-game-service/main -i cmd/game/game_server.go
RUN cp config.json colors-game-service/config.json


FROM alpine:latest
ENV COLORS_SERVICE_USE_MODE=IN_DOCKER_NET
ENV COLORS_DB=IN_DOCKER_NET
ENV COLORS_CONFIG_PATH="/root/colors-game-service"

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/colors-game-service ./colors-game-service
EXPOSE 8004
ENTRYPOINT ["/root/colors-game-service/main"]
