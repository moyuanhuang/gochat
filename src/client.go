package gochat

import (
    "fmt"
    "net"
    "os"
    "bufio"
    "encoding/gob"
)

const PROTOCOL = "tcp4"

type Client struct {
    UserName string
    Conn net.Conn
    ServiceAddr string
}

func NewClient(serviceAddr string) (*Client, error) {
    var userName string
    fmt.Println("Please enter your user name:")
    _, err := fmt.Scanln(&userName)
    if err != nil {
        return nil, err
    }
    return &Client{UserName: userName, ServiceAddr: serviceAddr}, nil
}

func sendMessage(c *Client) {
    scanner := bufio.NewScanner(os.Stdin, bufio.MaxScanTokenSize * 10)
    // enlarge the scanner's buffer size, or may not able to read a complete line, otherwise use bufio.Reader.ReadString('\n')
    // detailed explanation here: https://github.com/ma6174/blog/issues/10
    encoder := gob.NewEncoder(c.Conn)

    renameClient(c, encoder)

    for scanner.Scan() {
        // DONNOT use scanln, because it treats space as delimiter:
        // https://stackoverflow.com/questions/43843477/scanln-in-golang-doesnt-accept-whitespace
        // fmt.Scanln(&message)
        message := Message{Sender: c.UserName, Text: scanner.Text()}
        if message.Text == "/quit" {
            fmt.Println("Exiting chat room... Bye.")
            os.Exit(0)
        }

        err := encoder.Encode(message)
        HandleError(err)
    }
}

func receiveMsgFromServer(c *Client) {
    decoder := gob.NewDecoder(c.Conn)
    var message Message
    for {
        err := decoder.Decode(&message)
        HandleError(err)
        message.Print()
    }
}

func (c *Client) StartClient() {
    addr, err := net.ResolveTCPAddr(PROTOCOL, c.ServiceAddr)
    HandleError(err)
    conn, err := net.DialTCP(PROTOCOL, nil, addr)
    HandleError(err)
    c.Conn = conn

    go sendMessage(c)
    go receiveMsgFromServer(c)

    for { }
    os.Exit(0)
}

func renameClient(c *Client, enc *gob.Encoder){
    text := "/name " + c.UserName
    message := Message{Text: text}
    err := enc.Encode(message)
    HandleError(err)
}
