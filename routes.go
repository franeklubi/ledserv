package ledserv

import (
    "fmt"
    "log"
    "strconv"
    "net/http"
    "github.com/gorilla/websocket"
)


var (
    upgrader = websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
    }

    clients = make(map[*websocket.Conn]bool)

    ledservServeMux = http.NewServeMux()
)


func setupRoutes() {
    ledservServeMux.HandleFunc("/", mainEndpoint)
    ledservServeMux.HandleFunc("/ws", wsEndpoint)
}


func startServer(port uint16) {
    port_string := strconv.FormatUint(uint64(port), 10)
    log.Printf("Starting server on %s\n", port_string)

    err := http.ListenAndServe(":"+port_string, ledservServeMux)
    standardErrorHandler(err)
}


func mainEndpoint(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "/ws is the websocket endpoint")
}


func wsEndpoint(w http.ResponseWriter, r *http.Request) {
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }

    ws, err := upgrader.Upgrade(w, r, nil)
    if ( err != nil ) {
        log.Println("err")
        return
    }

    log.Println("Client successfully connected")

    clients[ws] = true
}


func sender() {
    for {
        bytes := <-websocket_send
        for c := range clients {
            err := c.WriteMessage(websocket.BinaryMessage, bytes)
            if ( err != nil ) {
                log.Println(err, "Deleting connection.")
                c.Close()
                delete(clients, c)
            }
        }
    }
}
