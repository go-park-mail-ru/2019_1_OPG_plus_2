FROM golang:latest
ENV GO111MODULE=on
ENV IN_DOCKER=on
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN mkdir static
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -i cmd/server/main.go

EXPOSE 8002
ENTRYPOINT ["/app/main"]
