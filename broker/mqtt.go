package broker

import (
	"bufio"
	"errors"
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

func getRequestType(b byte) byte {
	return (b & 0xF0) >> 4
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
