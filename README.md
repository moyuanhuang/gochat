# gochat
Simple chatroom application written in Go

## Install
```bash
git clone https://github.com/moyuanhuang/gochat.git
cd gochat
make
```

## Usage
- For server, type `./gochat server host:port` to start the server. If running the server locally, you can omit the host part.
- For client, type `./gochat client host:port` to enter an existing chatroom. Then type your chat name(required) to start chatting with others. If running the client locally, you can omit the host part.

## TODO list
- how to distinguish different connections. Right now I use the conn.RemoteAddr() as key of the connection, not sure whether this is appropriate though.

- Currently the server can't handle the case when two clients have the same `UserName`. This is because in `Server.handleBroadcast()`, I use `UserName` to distinguish different clients. The solution is to use `conn.RemoteAddr().String()` as the key of the map, however, this would require an extra map of `RemoteAddr -> UserName` in order to get the clients' name. It will require another map of if `UserName -> RemoteAddr` if we are to support *mention(@)* functionality(consider the case `@dahuang hello`, how to direct this message to dahuang only). The current solution thus seems to be the most cost-efficient one.

- when hitting `ctrl+c`, the server and all clients will stop
