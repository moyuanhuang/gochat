package gochat

import (
    "os"
    "net"
    "fmt"
    "strings"
    "encoding/gob"
)

type Server struct {
    Listener net.Listener
    Broadcasters map[string]*gob.Encoder
    MessageCh chan Message
}

func NewServer(serviceAddr string) (*Server, error) {
    addr, err := net.ResolveTCPAddr(PROTOCOL, serviceAddr)
    if err != nil {
        return nil, err
    }
    listener, err := net.ListenTCP(PROTOCOL, addr)
    if err != nil {
        return nil, err
    }
    server := &Server{Listener: listener}
    server.MessageCh = make(chan Message)
    server.Broadcasters = make(map[string]*gob.Encoder)
    return server, nil
}

func handleBroadcast(s *Server) {
    for {
        message := <-s.MessageCh
        fmt.Printf("Broadcasting message from %s...\n", message.Sender)
        if len(message.Receivers) == 0 {
            for key, b := range s.Broadcasters {
                if key != message.Sender {
                    err := b.Encode(message)
                    HandleError(err)
                }
            }
        }
        // else if { @userName }
        // be aware of non-existing userName
    }
}

func (s *Server) ListenClient(conn net.Conn){
    decoder :=  gob.NewDecoder(conn)
    for {
        var message Message
        err := decoder.Decode(&message)
        HandleError(err)
        if message.Sender == "" {
            cmds := strings.Fields(message.Text)
            switch cmds[0] {
            case "/name":
                s.Broadcasters[cmds[1]] = gob.NewEncoder(conn)
                fmt.Printf("%s named himself %s\n", conn.RemoteAddr().String(), cmds[1])
            default:
                fmt.Println("Unrecognized command: ", message.Text)
            }
            continue
        }
        s.MessageCh <- message
    }
}

func (s *Server) StartServer() {
    go handleBroadcast(s)
    for {
        conn, err := s.Listener.Accept()
        HandleError(err)
        fmt.Printf("%s entered the room!\n", conn.RemoteAddr())
        go s.ListenClient(conn)
    }
    os.Exit(0)
}
