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
	Conn                  net.Conn
	Rdr                   *packet.Reader
	connectFlags          *mqtt.ConnectFlags
	KeepAlive             uint16
	SessionExpiryInterval uint32
	ReceiveMaximum        uint16
}

func (client *Client) handleConnect(remainingLength int) (int, error) {
	fmt.Println("Handle Connect")

	// verify the protocol is set to 'MQTT'
	err := mqtt.VerifyProtocol(client.Rdr)
	if err != nil {
		return 0, err
	}

	// Check protocol version. Currently only supports v5.0
	b, err := client.Rdr.ReadByte()
	if err != nil {
		return 0, err
	} else if b != 5 {
		msg := fmt.Sprintf("This broker currently only supports MQTT v5.0. You specified: %d", b)
		return 0, errors.New(msg)
	}

	// Check the connect flags!
	b, err = client.Rdr.ReadByte()
	if err != nil {
		return 0, err
	}
	flags, err := mqtt.GetConnectFlags(b)
	if err != nil {
		return 0, err
	}
	client.connectFlags = flags

	keepAlive, err := mqtt.GetKeepAlive(client.Rdr)
	if err != nil {
		return 0, err
	}
	client.KeepAlive = keepAlive

	fmt.Println(client)

	// And now for the properties!
	// TODO
	_, propLength, err := mqtt.DecodeVarByteInt(client.Rdr)
	if err != nil {
		return 0, err
	}
	fmt.Println(propLength)

	props, userProps, err := mqtt.GetProps(client.Rdr, int(propLength), mqtt.ConnectCode)
	if err != nil {
		return 0, err
	}
	// TODO
	// Do things with these props. Use a default value if not found.
	// Put them all in a map, then assign
	fmt.Println(props)
	fmt.Println(userProps)

	// Process the payload.
	// first thing is a utf8 encoded string for the clientId.
	// If zero length, assign one from the server.
	clientId, err := mqtt.GetClientId(client.Rdr)
	if err != nil {
		return 0, err
	}
	fmt.Printf("Client ID: %v\n", clientId)

	return remainingLength, nil
}
func (client *Client) handleConnack(remainingLength int) (int, error) {
	fmt.Println("Handle Connack")
	log.Fatalln("Not yet implemented.")
	return remainingLength, nil
}
func (client *Client) handlePublish(remainingLength int) (int, error) {
	fmt.Println("Handle Publish")
	log.Fatalln("Not yet implemented.")
	return remainingLength, nil
}
func (client *Client) handlePuback(remainingLength int) (int, error) {
	fmt.Println("Handle Puback")
	log.Fatalln("Not yet implemented.")
	return remainingLength, nil
}
func (client *Client) handlePubrec(remainingLength int) (int, error) {
	fmt.Println("Handle Pubrec")
	log.Fatalln("Not yet implemented.")
	return remainingLength, nil
}
func (client *Client) handlePubrel(remainingLength int) (int, error) {
	fmt.Println("Handle Pubrel")
	log.Fatalln("Not yet implemented.")
	return remainingLength, nil
}
func (client *Client) handlePubcomp(remainingLength int) (int, error) {
	fmt.Println("Handle Pubcomp")
	log.Fatalln("Not yet implemented.")
	return remainingLength, nil
}
func (client *Client) handleSubscribe(remainingLength int) (int, error) {
	fmt.Println("Handle Subscribe")
	log.Fatalln("Not yet implemented.")
	return remainingLength, nil
}
func (client *Client) handleSuback(remainingLength int) (int, error) {
	fmt.Println("Handle Suback")
	log.Fatalln("Not yet implemented.")
	return remainingLength, nil
}
func (client *Client) handleUnsubscribe(remainingLength int) (int, error) {
	fmt.Println("Handle Unsubscribe")
	log.Fatalln("Not yet implemented.")
	return remainingLength, nil
}
func (client *Client) handleUnsuback(remainingLength int) (int, error) {
	fmt.Println("Handle Unsuback")
	log.Fatalln("Not yet implemented.")
	return remainingLength, nil
}
func (client *Client) handlePingreq(remainingLength int) (int, error) {
	fmt.Println("Handle Pingreq")
	log.Fatalln("Not yet implemented.")
	return remainingLength, nil
}
func (client *Client) handlePingresp(remainingLength int) (int, error) {
	fmt.Println("Handle Pingresp")
	log.Fatalln("Not yet implemented.")
	return remainingLength, nil
}
func (client *Client) handleDisconnect(remainingLength int) (int, error) {
	fmt.Println("Handle Disconnect")
	log.Fatalln("Not yet implemented.")
	return remainingLength, nil
}
func (client *Client) handleAuth(remainingLength int) (int, error) {
	fmt.Println("Handle Auth")
	log.Fatalln("Not yet implemented.")
	return remainingLength, nil
}

// processFixedHeader processes the fixed header.
// Returns the request type code, remaining length of the packet, and maybe an error.
func (client *Client) processFixedHeader() (byte, int, error) {
	fmt.Println("Fixed header:")
	b1, err := client.Rdr.ReadByte()
	if err != nil {
		return 0x00, 0, err
	}

	reqType := mqtt.GetRequestType(b1)
	if reqType == mqtt.PublishCode {
		// TODO Do thing with flags!
		// flags := (b1 & 0xF)
		log.Fatalln("Not yet implemented.")
	}

	_, remainingLength, err := mqtt.DecodeVarByteInt(client.Rdr)
	if err != nil {
		return 0x00, 0, err
	}

	return reqType, int(remainingLength), nil
}

func (client *Client) processVarHeader(reqType byte, remainingLength int) (remLength int, err error) {
	fmt.Println("The rest:")

	switch reqType {
	case mqtt.ConnectCode:
		remLength, err = client.handleConnect(remainingLength)
	case mqtt.ConnackCode:
		remLength, err = client.handleConnack(remainingLength)
	case mqtt.PublishCode:
		remLength, err = client.handlePublish(remainingLength)
	case mqtt.PubackCode:
		remLength, err = client.handlePuback(remainingLength)
	case mqtt.PubrecCode:
		remLength, err = client.handlePubrec(remainingLength)
	case mqtt.PubrelCode:
		remLength, err = client.handlePubrel(remainingLength)
	case mqtt.PubcompCode:
		remLength, err = client.handlePubcomp(remainingLength)
	case mqtt.SubscribeCode:
		remLength, err = client.handleSubscribe(remainingLength)
	case mqtt.SubackCode:
		remLength, err = client.handleSuback(remainingLength)
	case mqtt.UnsubscribeCode:
		remLength, err = client.handleUnsubscribe(remainingLength)
	case mqtt.UnsubackCode:
		remLength, err = client.handleUnsuback(remainingLength)
	case mqtt.PingreqCode:
		remLength, err = client.handlePingreq(remainingLength)
	case mqtt.PingrespCode:
		remLength, err = client.handlePingresp(remainingLength)
	case mqtt.DisconnectCode:
		remLength, err = client.handleDisconnect(remainingLength)
	case mqtt.AuthCode:
		remLength, err = client.handleAuth(remainingLength)
	default:
		msg := fmt.Sprintf("No matching case for request type: %d", reqType)
		return remLength, errors.New(msg)
	}
	return remLength, err

}

func (client *Client) ProcessPacket() error {
	// process the fixed header.
	reqType, remLen, err := client.processFixedHeader() // make this guy return remaining length!
	if err != nil {
		return err
	}
	// check the protocol is correct then run the switch statement currently in processFixedHeader
	remLen, err = client.processVarHeader(reqType, remLen)
	return err
}
