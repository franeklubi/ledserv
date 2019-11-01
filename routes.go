package ledserv

import (
    "fmt"
    "log"
    "strconv"
    "net/http"
    "github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}


func setupRoutes() {
    http.HandleFunc("/", mainEndpoint)
    http.HandleFunc("/ws", wsEndpoint)
}


func startServer(port uint16) {
    port_string := strconv.FormatUint(uint64(port), 10)
    log.Printf("Starting server on %s\n", port_string)
    err := http.ListenAndServe(":"+port_string, nil)
    standardErrorHandler(err)
}


func mainEndpoint(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "/ws is the websocket endpoint")
}


func wsEndpoint(w http.ResponseWriter, r *http.Request) {
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }

    ws, err := upgrader.Upgrade(w, r, nil)
    standardErrorHandler(err)

    log.Println("Client successfully connected")

    sender(ws)
}


func sender(conn *websocket.Conn) {
    for {
        bytes := <-websocket_send
        err := conn.WriteMessage(websocket.BinaryMessage, bytes)
        if ( err != nil ) {
            log.Println(err)
        }
    }
}
