package main

import (
    "../../internal/app/server"
    "fmt"
    "os"
)

func main() {
    params := server.Params{Port: os.Getenv("PORT")}
    if params.Port == "" {
        params.Port = "8001"
    }

    err := server.StartApp(params)
    if err != nil {
        fmt.Println()
    }
}
