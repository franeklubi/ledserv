package ledserv

import (
    "net"
    "time"
    "strings"
)


var (
    is_broadcasting bool = false
)


func Broadcast(ms float64) {
    if ( is_broadcasting ) {
        return
    }

    is_broadcasting = true
    broadcast(ms)
}


func StopBroadcast() {
    is_broadcasting = false
}


func broadcast(ms float64) {
    networks := getLocalNetworks()

    var broadcasts []net.IP
    for _, n := range networks {
        broadcasts = append(broadcasts, getBroadcastAddress(*n))
    }

    var connections []net.Conn
    for _, a := range broadcasts {
        c, err := net.Dial("udp", a.String()+":10107")
        standardErrorHandler(err)

        connections = append(connections, c)
    }

    for _, c := range connections {
        go func(c net.Conn) {
            defer c.Close()

            local_address := strings.Split(c.LocalAddr().String(), ":")[0]
            local_address_json := `{"address":"`+local_address+`"}`

            for {

                c.Write([]byte(local_address_json))
                time.Sleep(time.Millisecond*time.Duration(ms))

                if ( !is_broadcasting ) {
                    return
                }
            }
        }(c)
    }
}


func getLocalNetworks() ([]*net.IPNet) {
    interfaces, err := net.Interfaces()
    standardErrorHandler(err)

    var found_addresses []*net.IPNet

    for _, i := range interfaces {
        addresses, err := i.Addrs()
        standardErrorHandler(err)

        for _, a := range addresses {
            network, cast_successful := a.(*net.IPNet)
            if ( !cast_successful ) {
                break
            }

            if ( verifyAddress(network) ) {
                found_addresses = append(found_addresses, network)
            }
        }
    }

    return found_addresses
}


func getBroadcastAddress(n net.IPNet) (net.IP) {
    broadcast_address := n.IP.To4()
    for x := range broadcast_address {
        broadcast_address[x] = broadcast_address[x] | ^n.Mask[x]
    }

    return broadcast_address
}


func verifyAddress(network *net.IPNet) (bool) {
    return !network.IP.IsLoopback() && !network.IP.IsUnspecified() &&
        !strings.ContainsRune(network.IP.String(), ':')
}


func standardErrorHandler(err error) {
    if ( err != nil ) {
        panic(err)
    }
}
