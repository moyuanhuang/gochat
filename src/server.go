package gochat

import (
    "os"
    "net"
    "fmt"
    "strings"
    "encoding/gob"
)

const SERVER = "Server"

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
        message.Print()
        if len(message.Receivers) == 0 {  // To All
            for key, b := range s.Broadcasters {
                if key != message.Sender {
                    err := b.Encode(message)
                    HandleError(err)
                }
            }
        } else {  // direct message
            for _, r := range message.Receivers {
                if b, ok := s.Broadcasters[r]; ok {
                    err := b.Encode(message)
                    HandleError(err)
                } else {
                    reply := NewMessage(SERVER, r + " doesn't exist!")
                    err := s.Broadcasters[message.Sender].Encode(reply)
                    HandleError(err)
                }
            }
        }
    }
}

func (s *Server) ListenClient(conn net.Conn){
    decoder :=  gob.NewDecoder(conn)
    for {
        var message Message
        err := decoder.Decode(&message)
        if err != nil {
            s.disconectClient(conn.RemoteAddr().String())
            fmt.Printf("Lost connection with %s\n", conn.RemoteAddr().String())
            return
        }
        if isCommand(message) {
            cmds := strings.Fields(message.Text)
            switch cmds[0] {
            case "/name":
                s.Broadcasters[cmds[1]] = gob.NewEncoder(conn)
                fmt.Printf("%s named himself %s\n", conn.RemoteAddr().String(), cmds[1])
            case "/quit":
                s.disconectClient(message.Sender)
                return
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

func (s *Server) disconectClient(clientName string) {
    quitMessage := NewMessage(SERVER, clientName + " quited the chatroom!")
    s.MessageCh <- *quitMessage
    delete(s.Broadcasters, clientName)
}
