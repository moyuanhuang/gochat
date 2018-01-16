package gochat

import (
    "fmt"
)

type Message struct {
    Sender string
    Receivers []string
    Text string
}

func (m Message) Print() {
    fmt.Printf("%s: %s\n", m.Sender, m.Text)
}
