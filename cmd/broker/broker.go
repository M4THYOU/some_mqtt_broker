package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/M4THYOU/some_mqtt_broker/internal/client"
	"github.com/M4THYOU/some_mqtt_broker/internal/defaults"
	"github.com/M4THYOU/some_mqtt_broker/pkg/packet"
)

func listen(c *client.Client) {
	defer c.Conn.Close()
	// to timeout hanging/broken connections.
	c.Conn.SetDeadline(time.Now().Add(time.Second * 60))
	for {
		err := c.ProcessPacket()
		if err != nil {
			fmt.Println("Error processing:", err.Error())
			break
		}
	}

}

func main() {
	fmt.Println("Starting the server...")
	host := defaults.Host + ":" + defaults.Port
	l, err := net.Listen(defaults.ConType, host)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	fmt.Printf("Listening on %v\n\n", host)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(2)
		}
		client := &client.Client{Conn: conn, Rdr: packet.NewReader(conn, 0)}
		go listen(client)
	}

}
