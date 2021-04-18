package client

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/M4THYOU/some_mqtt_broker/mqtt"
	"github.com/M4THYOU/some_mqtt_broker/packet"
)

type Client struct {
	Conn         net.Conn
	Rdr          *packet.Reader
	connectFlags *mqtt.ConnectFlags
	KeepAlive    uint16
	UserName     string
	Password     []byte

	// Connect Properties
	SessionExpiryInterval uint32
	ReceiveMaximum        uint16
	MaxPacketSize         uint32 // because go is go, the default value is 0. However, according to spec the client must not set this to zero. Therefore, 0 => no limit.
	TopicAliasMaximum     uint16
	ReturnResponseInfo    bool // if true, return response info to the client in connack.
	ReturnProblemInfo     bool // if true, maybe send proper reason string to client (CONNACK & DISCONNECT). Depends on log level though.
	AuthMethod            string
	AuthData              []byte

	WillProps *mqtt.WillProps
}

// processFixedHeader processes the fixed header.
// Returns the request type code, remaining length of the packet, and maybe an error.
func (client *Client) processFixedHeader() (byte, int, error) {
	b1, err := client.Rdr.ReadByte()
	if err != nil {
		return 0x00, 0, err
	}
	fmt.Println("(fixed header)")

	reqType := mqtt.GetRequestType(b1)
	if reqType == mqtt.PublishCode {
		// TODO Do thing with flags!
		// flags := (b1 & 0xF)
		log.Fatalln("Not yet implemented.")
	}

	_, remainingLength, err := client.Rdr.ReadVarByteInt()
	if err != nil {
		return 0x00, 0, err
	}

	return reqType, int(remainingLength), nil
}

// processVarHeader processes the variable header and payload of the packet (if payload exists)
func (client *Client) processVarHeader(reqType byte) (err error) {
	fmt.Println("(var header and payload)")

	switch reqType {
	case mqtt.ConnectCode:
		err = client.handleConnect()
	case mqtt.ConnackCode:
		err = client.handleConnack()
	case mqtt.PublishCode:
		err = client.handlePublish()
	case mqtt.PubackCode:
		err = client.handlePuback()
	case mqtt.PubrecCode:
		err = client.handlePubrec()
	case mqtt.PubrelCode:
		err = client.handlePubrel()
	case mqtt.PubcompCode:
		err = client.handlePubcomp()
	case mqtt.SubscribeCode:
		err = client.handleSubscribe()
	case mqtt.SubackCode:
		err = client.handleSuback()
	case mqtt.UnsubscribeCode:
		err = client.handleUnsubscribe()
	case mqtt.UnsubackCode:
		err = client.handleUnsuback()
	case mqtt.PingreqCode:
		err = client.handlePingreq()
	case mqtt.PingrespCode:
		err = client.handlePingresp()
	case mqtt.DisconnectCode:
		err = client.handleDisconnect()
	case mqtt.AuthCode:
		err = client.handleAuth()
	default:
		msg := fmt.Sprintf("No matching case for request type: %d", reqType)
		return errors.New(msg)
	}
	return err

}

func (client *Client) ProcessPacket() error {
	fmt.Println("Waiting for packet...")
	// fixed header can be up to 5 bytes, so set that as the limit.
	client.Rdr.SetRemainingLength(5)
	reqType, remLen, err := client.processFixedHeader() // make this guy return remaining length!
	if err != nil {
		return err
	}
	client.Rdr.SetRemainingLength(remLen)
	err = client.processVarHeader(reqType)
	fmt.Printf("Packet processed.\n\n")
	return err
}
