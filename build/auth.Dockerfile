FROM golang:latest as builder
ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN mkdir static
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -i cmd/auth/auth_server.go


FROM alpine:latest
ENV COLORS_SERVICE_USE_MODE=IN_DOCKER_NET
ENV COLORS_DB=IN_DOCKER_NET
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app .
EXPOSE 50242
ENTRYPOINT ["/root/auth_server"]
