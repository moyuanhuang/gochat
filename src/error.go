package gochat

import (
    "fmt"
    "os"
)

func HandleError(err error) {
    if err != nil {
        fmt.Println("Fatal Error: ", err.Error())
        os.Exit(0)
    }
}
