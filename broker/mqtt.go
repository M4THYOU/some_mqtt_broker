package broker

import (
	"bufio"
	"encoding/binary"
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
	willPropsCode   = 0x00 // not defined by the spec, but we use this in getProps.
)

// Define all the packet structs.
type Connect struct {
	// Variable Header
	Flags     ConnectFlags
	KeepAlive uint16
	// Properties (Still in Variable Header)
	SessionExpiryInterval      uint32
	ReceiveMaximum             uint16
	MaximumPacketSize          uint32
	TopicAliasMaximum          uint16
	RequestResponseInformation uint8 // 0 or 1.
	RequestProblemInformation  uint8 // 0 or 1.
	UserProperty               utils.Utf8StringPair
	AuthMethod                 []byte // up to 65,535 bytes.
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

// getClientId gets the client ID from the next available bytes in the reader.
// If length is 0, assigns one randomly.
func getClientId(rdr *bufio.Reader) (string, error) {
	msb, err := rdr.ReadByte()
	if err != nil {
		return "", err
	}
	lsb, err := rdr.ReadByte()
	if err != nil {
		return "", err
	}
	len := int(binary.BigEndian.Uint16([]byte{msb, lsb}))
	if len == 0 {
		// TODO implement some auto assigning clientId method.
		return "one randomly", nil // LOL
	}
	fmt.Println(len)
	s := make([]byte, 0)
	for i := 0; i < len; i++ {
		b, err := rdr.ReadByte()
		if err != nil {
			return "", err
		}
		s = append(s, b)
	}
	fmt.Println(s)
	fmt.Println(string(s[:]))
	return string(s), nil
}

// getStringPropParams should only be called by getProps. Note that this function also works for Binary Data.
// Returns the number of bytes to be read, the new total of bytes read, and possibly an error.
func getStringPropParams(i int, rdr *bufio.Reader) (count, newI int, err error) {
	msb, err := rdr.ReadByte()
	if err != nil {
		return 0, 0, err
	}
	lsb, err := rdr.ReadByte()
	if err != nil {
		return 0, 0, err
	}
	buf := []byte{msb, lsb}
	count = int(binary.BigEndian.Uint16(buf))
	return count, i + 2, nil
}

// getBinaryDataPropParams should only be called by getProps. This function just calls getStringPropParams.
// returns the number of bytes to be read, the new total of bytes read, and possibly an error.
func getBinaryDataPropParams(i int, rdr *bufio.Reader) (count, newI int, err error) {
	return getStringPropParams(i, rdr)
}

// getStringPairProp reads the entire UTF-8 string pair into a byte array. Should only be called by getProps.
// Returns number of bytes read, the slice, and possibly an error.
func getStringPairProp(rdr *bufio.Reader) (int, []byte, error) {
	buf := make([]byte, 0)
	msb, err := rdr.ReadByte()
	if err != nil {
		return 0, nil, err
	}
	lsb, err := rdr.ReadByte()
	if err != nil {
		return 0, nil, err
	}
	buf = append(buf, msb)
	buf = append(buf, lsb)
	count := int(binary.BigEndian.Uint16(buf))
	for i := 0; i < count; i++ {
		b, err := rdr.ReadByte()
		if err != nil {
			return 0, nil, err
		}
		buf = append(buf, b)
	}
	bytesRead := count + 2

	// Do the exact same thing as above, for the second string.
	msb, err = rdr.ReadByte()
	if err != nil {
		return 0, nil, err
	}
	lsb, err = rdr.ReadByte()
	if err != nil {
		return 0, nil, err
	}
	buf = append(buf, msb)
	buf = append(buf, lsb)
	numBuf := []byte{msb, lsb}
	count = int(binary.BigEndian.Uint16(numBuf))
	for i := 0; i < count; i++ {
		b, err := rdr.ReadByte()
		if err != nil {
			return 0, nil, err
		}
		buf = append(buf, b)
	}
	bytesRead += (count + 2)

	return bytesRead, buf, nil
}

// getProps gets all the properties for this packet. Throws error if that prop is not valid for the specified packetCode.
// packetCode is one of the defined
func getProps(rdr *bufio.Reader, propLength, packetCode int) (map[int][]byte, [][]byte, error) {
	if (packetCode == pingreqCode) || (packetCode == pingrespCode) {
		return nil, nil, errors.New("getProps not valid for pingReq or pingResp packets")
	}

	m := make(map[int][]byte)
	userProps := make([][]byte, 0)

	for i := 0; i < propLength; {
		b, err := rdr.ReadByte()
		if err != nil {
			return nil, nil, err
		}
		count := 0
		var validCodes []int
		switch b {
		case 0x01:
			validCodes = []int{publishCode, willPropsCode}
			count = 1
		case 0x02:
			validCodes = []int{publishCode, willPropsCode}
			count = 4
		case 0x03, 0x08: // UTF-8 String
			validCodes = []int{publishCode, willPropsCode}
			count, i, err = getStringPropParams(i, rdr)
			if err != nil {
				return nil, nil, err
			}
		case 0x09: // Binary data
			validCodes = []int{publishCode, willPropsCode}
			count, i, err = getBinaryDataPropParams(i, rdr)
			if err != nil {
				return nil, nil, err
			}
		case 0x0B: // Variable Byte Integer
			validCodes = []int{publishCode, subscribeCode}
			count = 0
		case 0x11:
			validCodes = []int{connectCode, connackCode, disconnectCode}
			count = 4
		case 0x12:
			validCodes = []int{connackCode}
			count, i, err = getStringPropParams(i, rdr)
			if err != nil {
				return nil, nil, err
			}
		case 0x13:
			validCodes = []int{connackCode}
			count = 2
		case 0x15:
			validCodes = []int{connectCode, connackCode, authCode}
			count, i, err = getStringPropParams(i, rdr)
			if err != nil {
				return nil, nil, err
			}
		case 0x16:
			validCodes = []int{connectCode, connackCode, authCode}
			count, i, err = getBinaryDataPropParams(i, rdr)
			if err != nil {
				return nil, nil, err
			}
		case 0x17, 0x19:
			validCodes = []int{connectCode}
			count = 1
		case 0x18:
			validCodes = []int{willPropsCode}
			count = 4
		case 0x1A:
			validCodes = []int{connackCode}
			count, i, err = getStringPropParams(i, rdr)
			if err != nil {
				return nil, nil, err
			}
		case 0x1C:
			validCodes = []int{connackCode, disconnectCode}
			count, i, err = getStringPropParams(i, rdr)
			if err != nil {
				return nil, nil, err
			}
		case 0x1F:
			validCodes = []int{connackCode, pubackCode, pubrecCode, pubrelCode, pubcompCode, subackCode, unsubackCode, disconnectCode, authCode}
			count, i, err = getStringPropParams(i, rdr)
			if err != nil {
				return nil, nil, err
			}
		case 0x21, 0x22:
			validCodes = []int{connectCode, connackCode}
			count = 2
		case 0x23:
			validCodes = []int{publishCode}
			count = 2
		case 0x24, 0x25:
			validCodes = []int{connackCode}
			count = 1
		case 0x26:
			validCodes = []int{connectCode, connackCode, publishCode, willPropsCode, pubackCode, pubrecCode, pubrelCode, pubcompCode, subscribeCode, subackCode, unsubscribeCode, unsubackCode, disconnectCode, authCode}
		case 0x27:
			validCodes = []int{connectCode, connackCode}
			count = 4
		case 0x28, 0x29, 0x2A:
			validCodes = []int{connackCode}
			count = 1
		default:
			msg := fmt.Sprintf("No matching case for code: %d", b)
			return nil, nil, errors.New(msg)
		}

		isValid := utils.IsIntInSlice(packetCode, validCodes)
		if !isValid {
			msg := fmt.Sprintf("invalid property identifier %d for packet type %d", b, packetCode)
			return nil, nil, errors.New(msg)
		}

		var prop []byte
		if b == 0x26 { // get the string pair.
			count, prop, err = getStringPairProp(rdr)
			if err != nil {
				return nil, nil, err
			}
			userProps = append(userProps, prop)
		} else if count == 0 { // this indicates it must be a Variable Byte Integer!
			numRead, val, err := decodeVarByteInt(rdr)
			if err != nil {
				return nil, nil, err
			}
			count = numRead
			prop = make([]byte, 4)
			binary.BigEndian.PutUint32(prop, val)
			m[int(b)] = prop
		} else {
			prop, err = utils.ReadBytesToSlice(count, rdr)
			m[int(b)] = prop
		}
		if err != nil {
			return nil, nil, err
		}

		i += (count + 1)
	}

	return m, userProps, nil
}

// decodeVarByteInt returns the integer value of a decoded Variable Byte Int according to MQTT v5.0 Spec.
// Returns number of bytes read, the integer, and possibly an error.
func decodeVarByteInt(rdr *bufio.Reader) (int, uint32, error) {
	var multiplier, val uint32 = 1, 0
	var b byte
	var err error
	i := 0
	for {
		b, err = rdr.ReadByte()
		i++
		if err != nil {
			return 0, 0, err
		}
		val += uint32(b&0x7F) * multiplier
		if multiplier > 128*128*128 {
			return 0, 0, errors.New("malformed variable byte integer")
		}
		multiplier *= 128
		if (b & 0x80) == 0 {
			break
		}
	}
	return i, val, nil
}

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
	msb, err := rdr.ReadByte()
	if err != nil {
		return err
	}
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
	if msb != 0x00 || lsb != 0x04 || m != 'M' || q != 'Q' || t1 != 'T' || t2 != 'T' {
		msg := fmt.Sprintf("Got invalid protocol:\n%08b\n%08b\n%08b\n%08b\n%08b\n\nExpected:\n%08b\n%08b\n%08b\n%08b\n%08b", lsb, m, q, t1, t2, 0x04, 'M', 'Q', 'T', 'T')
		return errors.New(msg)
	}
	return nil
}

func getKeepAlive(rdr *bufio.Reader) (uint16, error) {
	// Read 2 bytes and turn it to an integer
	buf := make([]byte, 0)
	for i := 0; i < 2; i++ {
		b, err := rdr.ReadByte()
		if err != nil {
			return 0, err
		}
		buf = append(buf, b)
	}
	keepAlive := binary.BigEndian.Uint16(buf)
	return keepAlive, nil
}
