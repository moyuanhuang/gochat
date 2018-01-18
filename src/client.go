package gochat

import (
    "fmt"
    "net"
    "os"
    "bufio"
    "encoding/gob"
    "strings"
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
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Buffer([]byte{}, bufio.MaxScanTokenSize * 10)
    // enlarge the scanner's buffer size, or may not able to read a complete line, otherwise use bufio.Reader.ReadString('\n')
    // detailed explanation here: https://github.com/ma6174/blog/issues/10
    encoder := gob.NewEncoder(c.Conn)

    renameClient(c, encoder)

    for scanner.Scan() {
        // DONNOT use scanln, because it treats space as delimiter:
        // https://stackoverflow.com/questions/43843477/scanln-in-golang-doesnt-accept-whitespace
        // fmt.Scanln(&message)
        rawMsg := scanner.Text()

        if isQuit(rawMsg) {
            message := Message{Sender: c.UserName, Text: rawMsg}
            err := encoder.Encode(message)
            HandleError(err)
            fmt.Println("Exiting chat room... Bye.")
            os.Exit(0)
        } else {
            message := generateMessage(rawMsg, c)
            err := encoder.Encode(message)
            HandleError(err)
        }
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

func isQuit(rawMsg string) bool {
    quit := "/quit"
    return strings.HasPrefix(rawMsg, quit)
}

func generateMessage(rawMsg string, c *Client) *Message {
    message := &Message{Sender: c.UserName, Text: rawMsg}

    // handling different commands
    cmds := strings.Fields(rawMsg)
    switch cmds[0] {
    case "/dm":  // direct message
        message.Receivers = []string{cmds[1]}
        message.Text = strings.TrimPrefix(rawMsg[len(cmds[0] + cmds[1]):], " ")
    // TODO: /group (group chat)
    }

    return message
}
