package ledserv

import (
    "net"
    "time"
    // "fmt"
)


var (
    is_broadcasting bool = false
)


func Broadcast(ms float64) {
    if ( is_broadcasting ) {
        return
    }

    is_broadcasting = true
    go broadcast(ms)
}


func StopBroadcast() {
    is_broadcasting = false
}


func broadcast(ms float64) {
    connection, err := net.Dial("udp", "255.255.255.255:10107")
    if ( err != nil ) {
        panic(err)
    }
    defer connection.Close()

    local_address_json := `{"address":"`+connection.LocalAddr().String()+`"}`

    for {
        // fmt.Println(local_address_json)
        connection.Write([]byte(local_address_json))
        time.Sleep(time.Millisecond*time.Duration(ms))


        if ( !is_broadcasting ) {
            return
        }
    }
}
