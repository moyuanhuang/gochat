package main

import (
    "fmt"
    "os"
    "./src"
)

func main() {
    if len(os.Args) != 3 {
        printUsage()
        os.Exit(0)
    }
    servType := os.Args[1]
    serviceAddr := os.Args[2]

    if servType == "client" {
        c, err := gochat.NewClient(serviceAddr)
        handleError(err)
        c.StartClient()
    } else if servType == "server" {
        s, err := gochat.NewServer(serviceAddr)
        handleError(err)
        s.StartServer()
    } else {
        printUsage()
    }
    os.Exit(0)
}

func printUsage() {
    fmt.Printf("Usage %s [client/server] host:port\n", os.Args[0])
}

func handleError(err error) {
    if err != nil {
        fmt.Println("Fatal Error: ", err.Error())
        os.Exit(0)
    }
}
