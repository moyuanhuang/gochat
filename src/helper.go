package gochat

import (
    "fmt"
    "os"
    "strings"
)

func HandleError(err error) {
    if err != nil {
        fmt.Println("Fatal Error: ", err.Error())
        os.Exit(0)
    }
}

func isCommand(message Message) bool {
    text := message.Text
    return strings.HasPrefix(text, "/quit") || strings.HasPrefix(text, "/name")
}
