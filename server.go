package main

import (
	"fmt"
	"net"
	"os"

	"github.com/M4THYOU/some_mqtt_broker/config"
)

type utf8Str = []byte
type utf8StrPair = utf8Str
type variableByteInt = []uint8
type oneByteInt = uint8
type twoByteInt = uint16
type fourByteInt = uint32
type binaryData struct {
	length twoByteInt
	data   []byte // of the specified length.
}

// Define all the packet structs.
type Connect struct {
	// Variable Header
	ProtocolName    utf8Str // UTF-8 encoded string, must be 'MQTT'
	ProtocolVersion oneByteInt
	Flags           ConnectFlags
	KeepAlive       twoByteInt
	// Properties (Still in Variable Header)
	PropertyLength             variableByteInt
	SessionExpiryInterval      fourByteInt
	ReceiveMaximum             twoByteInt
	MaximumPacketSize          fourByteInt
	TopicAliasMaximum          twoByteInt
	RequestResponseInformation oneByteInt // 0 or 1.
	RequestProblemInformation  oneByteInt // 0 or 1.
	UserProperty               utf8StrPair
	AuthMethod                 utf8Str
	AuthData                   binaryData
}
type ConnectFlags struct {
	CleanStart   bool
	WillFlag     bool
	WillQos      uint8 // consisting only of 2 bits. Valid values are 0, 1, 2. Not 3!
	WillRetain   bool
	PasswordFlag bool
	UserNameFlag bool
}
type Connack struct{}
type Publish struct{}
type Puback struct{}
type Pubrec struct{}
type Pubrel struct{}
type Pubcomp struct{}
type Subscribe struct{}
type Suback struct{}
type Unsubscribe struct{}
type Unsuback struct{}
type Pingreq struct{}
type Pingresp struct{}
type Disconnect struct{}
type Auth struct{}

func printRawBuffer(buf []byte, len int) {
	for i := 0; i < len; i++ {
		fmt.Printf("%d: %08b\n", i, buf[i])
	}
}

func handleRequest(conn net.Conn) {
	for {
		// Make buffer to hold incoming data.
		buf := make([]byte, config.MaxPacketSize)
		// Read incoming connection into the buffer.
		reqLen, err := conn.Read(buf)
		printRawBuffer(buf, reqLen)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			conn.Close()
			return
		}
		// Send response back to contacting device.
		result := fmt.Sprintf("Message received: %d", reqLen)
		fmt.Println(result)
		// conn.Write([]byte(result))
		// Close when done.
		// conn.Close()
	}

}

func main() {
	fmt.Println("Starting the server...")
	host := config.Host + ":" + config.Port
	l, err := net.Listen(config.ConType, host)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Listening on " + host)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(2)
		}
		go handleRequest(conn)
	}

}
