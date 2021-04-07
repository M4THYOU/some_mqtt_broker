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

func printRawBuffer(buf []byte, len int) {
	for i := 0; i < len; i++ {
		fmt.Printf("%d: %08b\n", i, buf[i])
	}
}

func listen(conn net.Conn) {
	defer conn.Close()
	rdr := bufio.NewReader(conn)
	for {
		conn.SetDeadline(time.Now().Add(time.Second * 120))
		err := mqtt.ProcessPacket(rdr)
		if err != nil {
			fmt.Println("Error processing:", err.Error())
			break
		}

		// // Make buffer to hold incoming data.
		// fmt.Println("a")
		// buf := make([]byte, config.MaxPacketSize)
		// // Read incoming connection into the buffer.
		// fmt.Println("b")
		// reqLen, err := conn.Read(buf)
		// fmt.Println("c")
		// printRawBuffer(buf, reqLen)
		// fmt.Println("d")
		// if err != nil {
		// 	fmt.Println("Error reading:", err.Error())
		// 	break
		// }
		// err = mqtt.ProcessPacket(rdr)
		// if err != nil {
		// 	fmt.Println("Error processing:", err.Error())
		// 	break
		// }
		// // Send response back to contacting device.
		// result := fmt.Sprintf("Message received: %d", reqLen)
		// fmt.Println(result)
		// break
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
