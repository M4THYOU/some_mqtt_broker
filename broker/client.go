package broker

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
)

type Client struct {
	Conn                  net.Conn
	Rdr                   *bufio.Reader
	connectFlags          *ConnectFlags
	KeepAlive             uint16
	SessionExpiryInterval uint32
	ReceiveMaximum        uint16
}

func (client *Client) handleConnect(remainingLength uint32) error {
	fmt.Println("Handle Connect")

	// verify the protocol is set to 'MQTT'
	err := verifyProtocol(client.Rdr)
	if err != nil {
		return err
	}

	// Check protocol version. Currently only supports v5.0
	b, err := client.Rdr.ReadByte()
	if err != nil {
		return err
	} else if b != 5 {
		msg := fmt.Sprintf("This broker currently only supports MQTT v5.0. You specified: %d", b)
		return errors.New(msg)
	}

	// Check the connect flags!
	b, err = client.Rdr.ReadByte()
	if err != nil {
		return err
	}
	flags, err := getConnectFlags(b)
	if err != nil {
		return err
	}
	client.connectFlags = flags

	keepAlive, err := getKeepAlive(client.Rdr)
	if err != nil {
		return err
	}
	client.KeepAlive = keepAlive

	fmt.Println(client)

	// And now for the properties!
	// TODO
	_, propLength, err := decodeVarByteInt(client.Rdr)
	if err != nil {
		return err
	}
	fmt.Println(propLength)

	// TODO
	// Set all the props. Use a defualt value if not found.
	// Put them all in a map, then assign

	return nil
}
func (client *Client) handleConnack(remainingLength uint32) error {
	fmt.Println("Handle Connack")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePublish(remainingLength uint32) error {
	fmt.Println("Handle Publish")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePuback(remainingLength uint32) error {
	fmt.Println("Handle Puback")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePubrec(remainingLength uint32) error {
	fmt.Println("Handle Pubrec")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePubrel(remainingLength uint32) error {
	fmt.Println("Handle Pubrel")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePubcomp(remainingLength uint32) error {
	fmt.Println("Handle Pubcomp")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleSubscribe(remainingLength uint32) error {
	fmt.Println("Handle Subscribe")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleSuback(remainingLength uint32) error {
	fmt.Println("Handle Suback")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleUnsubscribe(remainingLength uint32) error {
	fmt.Println("Handle Unsubscribe")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleUnsuback(remainingLength uint32) error {
	fmt.Println("Handle Unsuback")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePingreq(remainingLength uint32) error {
	fmt.Println("Handle Pingreq")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePingresp(remainingLength uint32) error {
	fmt.Println("Handle Pingresp")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleDisconnect(remainingLength uint32) error {
	fmt.Println("Handle Disconnect")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleAuth(remainingLength uint32) error {
	fmt.Println("Handle Auth")
	log.Fatalln("Not yet implemented.")
	return nil
}

func (client *Client) processFixedHeader() (byte, uint32, error) {
	fmt.Println("Fixed header:")
	b1, err := client.Rdr.ReadByte()
	if err != nil {
		return 0x00, 0, err
	}

	reqType := getRequestType(b1)
	if reqType == publishCode {
		// Do thing with flags!
		// flags := (b1 & 0xF)
		log.Fatalln("Not yet implemented.")
	}

	_, remainingLength, err := decodeVarByteInt(client.Rdr)
	if err != nil {
		return 0x00, 0, err
	}

	return reqType, remainingLength, nil
}

func (client *Client) processVarHeader(reqType byte, remainingLength uint32) (err error) {
	fmt.Println("The rest:")

	switch reqType {
	case connectCode:
		err = client.handleConnect(remainingLength)
	case connackCode:
		err = client.handleConnack(remainingLength)
	case publishCode:
		err = client.handlePublish(remainingLength)
	case pubackCode:
		err = client.handlePuback(remainingLength)
	case pubrecCode:
		err = client.handlePubrec(remainingLength)
	case pubrelCode:
		err = client.handlePubrel(remainingLength)
	case pubcompCode:
		err = client.handlePubcomp(remainingLength)
	case subscribeCode:
		err = client.handleSubscribe(remainingLength)
	case subackCode:
		err = client.handleSuback(remainingLength)
	case unsubscribeCode:
		err = client.handleUnsubscribe(remainingLength)
	case unsubackCode:
		err = client.handleUnsuback(remainingLength)
	case pingreqCode:
		err = client.handlePingreq(remainingLength)
	case pingrespCode:
		err = client.handlePingresp(remainingLength)
	case disconnectCode:
		err = client.handleDisconnect(remainingLength)
	case authCode:
		err = client.handleAuth(remainingLength)
	}
	return err

}

func (client *Client) ProcessPacket() error {
	// process the fixed header.
	reqType, remLen, err := client.processFixedHeader() // make this guy return remaining length!
	if err != nil {
		return err
	}
	// check the protocol is correct then run the switch statement currently in processFixedHeader
	err = client.processVarHeader(reqType, remLen)
	return err
}