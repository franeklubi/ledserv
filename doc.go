/*
Package ledserv provides an easy way to send github/franeklubi/ledgend Change
information to any websocket client.

Information is sent in the following form:
    The type of payload is []byte

    payload[0] (preambule) contains information about options, that change the
    way the payload should be read:
    Reserved options:
        - (payload[0] & 0x0F) - indicates how many bytes long the Changes'
            addresses will be
            (takes on the value of 2 when at least one of the Changes'
            address is above 255)

        - The rest of preambule's space is reserved for the situation
            there is ever a need to send over more options

    Let's assume the preambule tells us address length is == 1:
    payload[1] - Change's address
    payload[2] - Change's R value
    payload[3] - Change's G value
    payload[4] - Change's B value


The broadcast of ledserv's IP address looks like this:
    "ledgend;"+address
*/
package ledserv
