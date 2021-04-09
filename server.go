package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	mqtt "github.com/M4THYOU/some_mqtt_broker/broker"
	"github.com/M4THYOU/some_mqtt_broker/config"
)

func listen(conn net.Conn) {
	defer conn.Close()
	rdr := bufio.NewReader(conn)
	for {
		// to get rid of hanging/broken connections.
		conn.SetDeadline(time.Now().Add(time.Second * 60))
		err := mqtt.ProcessPacket(rdr)
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
		go listen(conn)
	}

}
