FROM golang:latest
ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/server/main.go

EXPOSE 8001
ENTRYPOINT ["/app/main"]