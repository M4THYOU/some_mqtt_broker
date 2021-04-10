package broker

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/M4THYOU/some_mqtt_broker/utils"
)

var _ = fmt.Printf // For debugging; delete when done.

const (
	connectCode     = 0x01
	connackCode     = 0x02
	publishCode     = 0x03
	pubackCode      = 0x04
	pubrecCode      = 0x05
	pubrelCode      = 0x06
	pubcompCode     = 0x07
	subscribeCode   = 0x08
	subackCode      = 0x09
	unsubscribeCode = 0x0A
	unsubackCode    = 0x0B
	pingreqCode     = 0x0C
	pingrespCode    = 0x0D
	disconnectCode  = 0x0E
	authCode        = 0x0F
)

type Client struct {
	Conn         net.Conn
	Rdr          *bufio.Reader
	connectFlags ConnectFlags
}

// Define all the packet structs.
type Connect struct {
	// Variable Header
	Flags     ConnectFlags
	KeepAlive utils.TwoByteInt
	// Properties (Still in Variable Header)
	PropertyLength             utils.VariableByteInt
	SessionExpiryInterval      utils.FourByteInt
	ReceiveMaximum             utils.TwoByteInt
	MaximumPacketSize          utils.FourByteInt
	TopicAliasMaximum          utils.TwoByteInt
	RequestResponseInformation utils.OneByteInt // 0 or 1.
	RequestProblemInformation  utils.OneByteInt // 0 or 1.
	UserProperty               utils.Utf8StrPair
	AuthMethod                 utils.Utf8Str
	AuthData                   utils.BinaryData
}
type ConnectFlags struct {
	UserNameFlag bool
	PasswordFlag bool
	WillRetain   bool
	WillQos      uint8 // consisting only of 2 bits. Valid values are 0, 1, 2. Not 3!
	WillFlag     bool
	CleanStart   bool
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

func getConnectFlags(b byte) (*ConnectFlags, error) {
	userNameFlag := ((b & 0x80) >> 7) == 1
	passwordFlag := ((b & 0x40) >> 6) == 1
	willRetain := ((b & 0x20) >> 5) == 1
	willQoS := ((b & 0x18) >> 3)
	willFlag := ((b & 0x04) >> 2) == 1
	cleanStart := ((b & 0x02) >> 1) == 1
	reserved := (b & 0x01) == 1
	if reserved {
		return nil, errors.New("invalid reserved bit")
	} else if willQoS > 2 {
		return nil, errors.New("invalid QoS")
	} else if !willFlag && (willRetain || (willQoS > 0)) {
		msg := fmt.Sprintf("Will Flag is: %t but Will Retain is %t and Will QoS is %d", willFlag, willRetain, willQoS)
		return nil, errors.New(msg)
	}
	flags := &ConnectFlags{userNameFlag, passwordFlag, willRetain, willQoS, willFlag, cleanStart}
	return flags, nil
}

func (client *Client) handleConnect(remainingLength uint64) error {
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
	fmt.Println(flags)

	return nil
}
func (client *Client) handleConnack(remainingLength uint64) error {
	fmt.Println("Handle Connack")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePublish(remainingLength uint64) error {
	fmt.Println("Handle Publish")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePuback(remainingLength uint64) error {
	fmt.Println("Handle Puback")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePubrec(remainingLength uint64) error {
	fmt.Println("Handle Pubrec")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePubrel(remainingLength uint64) error {
	fmt.Println("Handle Pubrel")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePubcomp(remainingLength uint64) error {
	fmt.Println("Handle Pubcomp")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleSubscribe(remainingLength uint64) error {
	fmt.Println("Handle Subscribe")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleSuback(remainingLength uint64) error {
	fmt.Println("Handle Suback")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleUnsubscribe(remainingLength uint64) error {
	fmt.Println("Handle Unsubscribe")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleUnsuback(remainingLength uint64) error {
	fmt.Println("Handle Unsuback")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePingreq(remainingLength uint64) error {
	fmt.Println("Handle Pingreq")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePingresp(remainingLength uint64) error {
	fmt.Println("Handle Pingresp")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleDisconnect(remainingLength uint64) error {
	fmt.Println("Handle Disconnect")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleAuth(remainingLength uint64) error {
	fmt.Println("Handle Auth")
	log.Fatalln("Not yet implemented.")
	return nil
}

func getRequestType(b byte) byte {
	return (b & 0xF0) >> 4
}

func (client *Client) processFixedHeader() (byte, uint64, error) {
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

	// Read all the bytes until 0x00 byte. This means it's about to write what protocol it is.
	bSlice := make([]byte, 0)
	b, err := client.Rdr.ReadByte()
	if err != nil {
		return 0x00, 0, err
	}
	for b != 0x00 {
		bSlice = append(bSlice, b)
		b, err = client.Rdr.ReadByte()
		if err != nil {
			return 0x00, 0, err
		}
	}
	remainingLength, n := binary.Uvarint(bSlice)
	if n == 0 {
		msg := fmt.Sprintf("Empty bSlice: %v\n", bSlice)
		return 0x00, 0, errors.New(msg)
	} else if n < 0 {
		msg := fmt.Sprintf("Overflow on: %v\n", bSlice)
		return 0x00, 0, errors.New(msg)
	}

	return reqType, remainingLength, nil
}

func verifyProtocol(rdr *bufio.Reader) error {
	lsb, err := rdr.ReadByte() // should be 00000100, i.e. 4
	if err != nil {
		return err
	}
	m, err := rdr.ReadByte()
	if err != nil {
		return err
	}
	q, err := rdr.ReadByte()
	if err != nil {
		return err
	}
	t1, err := rdr.ReadByte()
	if err != nil {
		return err
	}
	t2, err := rdr.ReadByte()
	if err != nil {
		return err
	}
	if lsb != 0x04 || m != 'M' || q != 'Q' || t1 != 'T' || t2 != 'T' {
		msg := fmt.Sprintf("Got invalid protocol:\n%08b\n%08b\n%08b\n%08b\n%08b\n\nExpected:\n%08b\n%08b\n%08b\n%08b\n%08b", lsb, m, q, t1, t2, 0x04, 'M', 'Q', 'T', 'T')
		return errors.New(msg)
	}
	return nil
}

func (client *Client) processVarHeader(reqType byte, remainingLength uint64) (err error) {
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
