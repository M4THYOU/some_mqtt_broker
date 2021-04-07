package broker

import (
	"bufio"
	"fmt"

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

func getRequestType(b byte) byte {
	return (b & 0xF0) >> 4
}

func processFixedHeader(rdr *bufio.Reader) error {
	fmt.Println("Fixed header:")
	b1, err := rdr.ReadByte()
	if err != nil {
		return err
	}
	reqType := getRequestType(b1)
	flags := (b1 & 0xF)
	fmt.Printf("%08b\n", b1)
	fmt.Printf("%08b\n", reqType)
	fmt.Printf("%08b\n", flags)
	return nil
}

func processVarHeader(rdr *bufio.Reader) error {
	fmt.Println("The rest:")
	for {
		b, err := rdr.ReadByte()
		if err != nil {
			return err
		}
		fmt.Printf("%08b\n", b)
	}
}

func ProcessPacket(rdr *bufio.Reader) error {
	// process the fixed header.
	err := processFixedHeader(rdr)
	if err != nil {
		return err
	}
	err = processVarHeader(rdr)
	return err
}
