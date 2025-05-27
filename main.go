package main

import (
    "fmt"
    "net/http"
    "os"
)

func main() {
mux := http.NewServeMux()

s := &http.Server{
Addr: ":8080",
Handler: mux,
}

mux.Handle("/", http.FileServer(http.Dir(".")))

err := s.ListenAndServe()
if err!=nil{
    fmt.Println("Couldn't start server")
    os.Exit(1)
    }
}

