package broker

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"log"

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

// Define all the packet structs.
type Connect struct {
	// Variable Header
	ProtocolName    utils.Utf8Str // UTF-8 encoded string, must be 'MQTT'
	ProtocolVersion utils.OneByteInt
	Flags           ConnectFlags
	KeepAlive       utils.TwoByteInt
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

func handleConnect(rdr *bufio.Reader, remainingLength uint64) {
	fmt.Println("Handle Connect")
	// log.Fatalln("Not yet implemented.")
}
func handleConnack(rdr *bufio.Reader, remainingLength uint64) {
	fmt.Println("Handle Connack")
	log.Fatalln("Not yet implemented.")
}
func handlePublish(rdr *bufio.Reader, remainingLength uint64) {
	fmt.Println("Handle Publish")
	log.Fatalln("Not yet implemented.")
}
func handlePuback(rdr *bufio.Reader, remainingLength uint64) {
	fmt.Println("Handle Puback")
	log.Fatalln("Not yet implemented.")
}
func handlePubrec(rdr *bufio.Reader, remainingLength uint64) {
	fmt.Println("Handle Pubrec")
	log.Fatalln("Not yet implemented.")
}
func handlePubrel(rdr *bufio.Reader, remainingLength uint64) {
	fmt.Println("Handle Pubrel")
	log.Fatalln("Not yet implemented.")
}
func handlePubcomp(rdr *bufio.Reader, remainingLength uint64) {
	fmt.Println("Handle Pubcomp")
	log.Fatalln("Not yet implemented.")
}
func handleSubscribe(rdr *bufio.Reader, remainingLength uint64) {
	fmt.Println("Handle Subscribe")
	log.Fatalln("Not yet implemented.")
}
func handleSuback(rdr *bufio.Reader, remainingLength uint64) {
	fmt.Println("Handle Suback")
	log.Fatalln("Not yet implemented.")
}
func handleUnsubscribe(rdr *bufio.Reader, remainingLength uint64) {
	fmt.Println("Handle Unsubscribe")
	log.Fatalln("Not yet implemented.")
}
func handleUnsuback(rdr *bufio.Reader, remainingLength uint64) {
	fmt.Println("Handle Unsuback")
	log.Fatalln("Not yet implemented.")
}
func handlePingreq(rdr *bufio.Reader, remainingLength uint64) {
	fmt.Println("Handle Pingreq")
	log.Fatalln("Not yet implemented.")
}
func handlePingresp(rdr *bufio.Reader, remainingLength uint64) {
	fmt.Println("Handle Pingresp")
	log.Fatalln("Not yet implemented.")
}
func handleDisconnect(rdr *bufio.Reader, remainingLength uint64) {
	fmt.Println("Handle Disconnect")
	log.Fatalln("Not yet implemented.")
}
func handleAuth(rdr *bufio.Reader, remainingLength uint64) {
	fmt.Println("Handle Auth")
	log.Fatalln("Not yet implemented.")
}

func getRequestType(b byte) byte {
	return (b & 0xF0) >> 4
}

func processFixedHeader(rdr *bufio.Reader) (byte, uint64, error) {
	fmt.Println("Fixed header:")
	b1, err := rdr.ReadByte()
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
	b, err := rdr.ReadByte()
	if err != nil {
		return 0x00, 0, err
	}
	for b != 0x00 {
		bSlice = append(bSlice, b)
		b, err = rdr.ReadByte()
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

func processVarHeader(rdr *bufio.Reader, reqType byte, remainingLength uint64) error {
	fmt.Println("The rest:")

	// TODO: verify the protocol is set to 'MQTT'
	// lsb, err := rdr.ReadByte() // should be 00000100, 4
	// if err != nil {
	// 	return err
	// }
	// m, err := rdr.ReadByte()
	// if err != nil {
	// 	return err
	// }
	// q, err := rdr.ReadByte()
	// if err != nil {
	// 	return err
	// }
	// t, err := rdr.ReadByte()
	// if err != nil {
	// 	return err
	// }
	// t, err = rdr.ReadByte()
	// if err != nil {
	// 	return err
	// }

	switch reqType {
	case connectCode:
		handleConnect(rdr, remainingLength)
	case connackCode:
		handleConnack(rdr, remainingLength)
	case publishCode:
		handlePublish(rdr, remainingLength)
	case pubackCode:
		handlePuback(rdr, remainingLength)
	case pubrecCode:
		handlePubrec(rdr, remainingLength)
	case pubrelCode:
		handlePubrel(rdr, remainingLength)
	case pubcompCode:
		handlePubcomp(rdr, remainingLength)
	case subscribeCode:
		handleSubscribe(rdr, remainingLength)
	case subackCode:
		handleSuback(rdr, remainingLength)
	case unsubscribeCode:
		handleUnsubscribe(rdr, remainingLength)
	case unsubackCode:
		handleUnsuback(rdr, remainingLength)
	case pingreqCode:
		handlePingreq(rdr, remainingLength)
	case pingrespCode:
		handlePingresp(rdr, remainingLength)
	case disconnectCode:
		handleDisconnect(rdr, remainingLength)
	case authCode:
		handleAuth(rdr, remainingLength)
	}

	// for {
	// 	b, err := rdr.ReadByte()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	fmt.Printf("%08b\n", b)
	// }

}

func ProcessPacket(rdr *bufio.Reader) error {
	// process the fixed header.
	reqType, remLen, err := processFixedHeader(rdr) // make this guy return remaining length!
	if err != nil {
		return err
	}
	// check the protocol is correct then run the switch statement currently in processFixedHeader
	err = processVarHeader(rdr, reqType, remLen)
	return err
}
