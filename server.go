package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/M4THYOU/some_mqtt_broker/client"
	"github.com/M4THYOU/some_mqtt_broker/config"
	"github.com/M4THYOU/some_mqtt_broker/packet"
)

func listen(c *client.Client) {
	defer c.Conn.Close()
	for {
		// to timeout hanging/broken connections.
		c.Conn.SetDeadline(time.Now().Add(time.Second * 60))
		err := c.ProcessPacket()
		if err != nil {
			fmt.Println("Error processing:", err.Error())
			break
		}
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
		client := &client.Client{Conn: conn, Rdr: packet.NewReader(conn)}
		go listen(client)
	}

}
