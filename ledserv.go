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

// InitServer initializes the ledserv's server on provided port
//
// Two and only endpoints are "/" and "/ws"
//
// "/ws" lets You establish a websocket connection
//
// "/" - informs the user user about the "/ws" endpoint
//
// Returns a channel You can send []ledgend.Change to,
// that will be sent to all connected clients
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
    go sender()
    go startServer(port)


    return user_send, nil
}


func changesReader(receive <-chan []ledgend.Change, send chan<- []byte) {
    for {
        changes := <-receive
        if ( len(changes) == 0 ) {
            continue;
        }

        var bytes []byte

        var address_bytes byte = 1
        for _, c := range changes {
            if ( c.Index > 255 ) {
                address_bytes = 2
                break;
            }
        }
        bytes = append(bytes, preambuleSetter(address_bytes))


        for _, c := range changes {
            var addr_a, addr_b byte
            addr_b = byte(c.Index & 0x00FF)

            if ( address_bytes == 2) {
                addr_a = byte(c.Index >> 8)
                bytes = append(bytes, addr_a)
            }

            bytes = append(bytes, addr_b, c.Pixel.R, c.Pixel.G, c.Pixel.B)
        }

        send<- bytes
    }
}


func preambuleSetter(address_bytes byte) (byte) {
    return address_bytes
}


func standardErrorHandler(err error) {
    if ( err != nil ) {
        log.Fatal(err)
    }
}
