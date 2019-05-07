FROM golang:latest as builder
ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN mkdir colors-auth-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o colors-auth-service/main -i cmd/auth/auth_server.go
RUN cp config.json colors-auth-service/config.json


FROM alpine:latest
ENV COLORS_SERVICE_USE_MODE=IN_DOCKER_NET
ENV COLORS_DB=IN_DOCKER_NET
ENV COLORS_CONFIG_PATH="/root/colors-auth-service"

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/colors-auth-service ./colors-auth-service
EXPOSE 50242
ENTRYPOINT ["/root/colors-auth-service/main"]
