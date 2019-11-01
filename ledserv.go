package ledserv

import (
    "log"
    "errors"
    "github.com/franeklubi/ledgend"
)


var (
    server_started  bool = false
    websocket_send  chan []byte
)


func InitServer(port uint16) (chan<- []ledgend.Change, error) {
    if ( server_started ) {
        err := errors.New("Server already started!")
        return nil, err;
    }
    server_started = true

    websocket_send = make(chan []byte, 100)
    user_send := make(chan []ledgend.Change, 100)

    go changesReader(user_send, websocket_send)


    setupRoutes()
    go startServer(port)


    return user_send, nil
}


func changesReader(receive <-chan []ledgend.Change, send chan<- []byte) {
    for {
        change := <-receive
        for _, c := range change {
            send<- []byte{c.Pixel.R, c.Pixel.G, c.Pixel.B}
        }
    }
}


func standardErrorHandler(err error) {
    if ( err != nil ) {
        log.Fatal(err)
    }
}
