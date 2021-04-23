package mqtt

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/M4THYOU/some_mqtt_broker/pkg/packet"
	"github.com/M4THYOU/some_mqtt_broker/pkg/utils"
)

var _ = fmt.Printf // For debugging; delete when done.

const (
	ConnectCode     = 0x01
	ConnackCode     = 0x02
	PublishCode     = 0x03
	PubackCode      = 0x04
	PubrecCode      = 0x05
	PubrelCode      = 0x06
	PubcompCode     = 0x07
	SubscribeCode   = 0x08
	SubackCode      = 0x09
	UnsubscribeCode = 0x0A
	UnsubackCode    = 0x0B
	PingreqCode     = 0x0C
	PingrespCode    = 0x0D
	DisconnectCode  = 0x0E
	AuthCode        = 0x0F
	WillPropsCode   = 0x00 // not defined by the spec, but we use this in getProps.
)

// All the properties!
const (
	PayloadFormatIndicatorCode = 0x01
	MessageExpiryIntervalCode  = 0x02
	ContentTypeCode            = 0x03
	ResponseTopicCode          = 0x08
	CorrelationDataCode        = 0x09
	SubscriptionIdCode         = 0x0B
	SessionExpiryIntervalCode  = 0x11
	AssignedClientIdCode       = 0x12
	ServerKeepAliveCode        = 0x13
	AuthenticationMethodCode   = 0x15
	AuthenticationDataCode     = 0x16
	RequestProblemInfoCode     = 0x17
	WillDelayIntervalCode      = 0x18
	RequestResponseInfoCode    = 0x19
	ResponseInfoCode           = 0x1A
	ServerReferenceCode        = 0x1C
	ReasonStringCode           = 0x1F
	ReceiveMaxCode             = 0x21
	TopicAliasMaxCode          = 0x22
	TopicAliasCode             = 0x23
	MaxQoSCode                 = 0x24
	RetainAvailableCode        = 0x25
	UserPropertyCode           = 0x26
	MaxPacketSizeCode          = 0x27
	WildcardSubAvailableCode   = 0x28
	SubIdAvailableCode         = 0x29
	SharedSubAvailableCode     = 0x2A
)

type ConnectFlags struct {
	UserNameFlag bool
	PasswordFlag bool
	WillRetain   bool
	WillQos      uint8 // consisting only of 2 bits. Valid values are 0, 1, 2. Not 3!
	WillFlag     bool
	CleanStart   bool
}
type WillProps struct {
	WillDelayInterval      uint32
	PayloadFormatIndicator uint8 // 0 or 1.
	MessageExpiryInterval  uint32
	ContentType            string
	ResponseTopic          string
	CorrelationData        []byte
	UserProperty           []byte
	Topic                  string
	Payload                []byte
}

// GetClientId gets the client ID from the next available bytes in the reader.
// If length is 0, assigns one randomly.
func GetClientId(rdr *packet.Reader) (string, error) {
	_, clientId, err := rdr.ReadUtf8Str()
	if err != nil {
		return "", err
	}
	if clientId == "" {
		// TODO implement some auto assigning clientId method.
		clientId = "one randomly" // LOL
	}
	return clientId, nil
}

// getStringPropParams should only be called by getProps. Note that this function also works for Binary Data.
// Returns the number of bytes to be read, the new total of bytes read, and possibly an error.
func getStringPropParams(i int, rdr *packet.Reader) (count, newI int, err error) {
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
func getBinaryDataPropParams(i int, rdr *packet.Reader) (count, newI int, err error) {
	return getStringPropParams(i, rdr)
}

// getStringPairProp reads the entire UTF-8 string pair into a byte array. Should only be called by getProps.
// Returns number of bytes read, the slice, and possibly an error.
func getStringPairProp(rdr *packet.Reader) (int, []byte, error) {
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

// GetProps gets all the properties for this packet. Throws error if that prop is not valid for the specified packetCode.
// packetCode is one of the defined
func GetProps(rdr *packet.Reader, propLength, packetCode int) (map[int][]byte, [][]byte, error) {
	if (packetCode == PingreqCode) || (packetCode == PingrespCode) {
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
		case PayloadFormatIndicatorCode:
			validCodes = []int{PublishCode, WillPropsCode}
			count = 1
		case MessageExpiryIntervalCode:
			validCodes = []int{PublishCode, WillPropsCode}
			count = 4
		case ContentTypeCode, ResponseTopicCode: // UTF-8 String
			validCodes = []int{PublishCode, WillPropsCode}
			count, i, err = getStringPropParams(i, rdr)
			if err != nil {
				return nil, nil, err
			}
		case CorrelationDataCode: // Binary data
			validCodes = []int{PublishCode, WillPropsCode}
			count, i, err = getBinaryDataPropParams(i, rdr)
			if err != nil {
				return nil, nil, err
			}
		case SubscriptionIdCode: // Variable Byte Integer
			validCodes = []int{PublishCode, SubscribeCode}
			count = 0
		case SessionExpiryIntervalCode:
			validCodes = []int{ConnectCode, ConnackCode, DisconnectCode}
			count = 4
		case AssignedClientIdCode:
			validCodes = []int{ConnackCode}
			count, i, err = getStringPropParams(i, rdr)
			if err != nil {
				return nil, nil, err
			}
		case ServerKeepAliveCode:
			validCodes = []int{ConnackCode}
			count = 2
		case AuthenticationMethodCode:
			validCodes = []int{ConnectCode, ConnackCode, AuthCode}
			count, i, err = getStringPropParams(i, rdr)
			if err != nil {
				return nil, nil, err
			}
		case AuthenticationDataCode:
			validCodes = []int{ConnectCode, ConnackCode, AuthCode}
			count, i, err = getBinaryDataPropParams(i, rdr)
			if err != nil {
				return nil, nil, err
			}
		case RequestProblemInfoCode, RequestResponseInfoCode:
			validCodes = []int{ConnectCode}
			count = 1
		case WillDelayIntervalCode:
			validCodes = []int{WillPropsCode}
			count = 4
		case ResponseInfoCode:
			validCodes = []int{ConnackCode}
			count, i, err = getStringPropParams(i, rdr)
			if err != nil {
				return nil, nil, err
			}
		case ServerReferenceCode:
			validCodes = []int{ConnackCode, DisconnectCode}
			count, i, err = getStringPropParams(i, rdr)
			if err != nil {
				return nil, nil, err
			}
		case ReasonStringCode:
			validCodes = []int{ConnackCode, PubackCode, PubrecCode, PubrelCode, PubcompCode, SubackCode, UnsubackCode, DisconnectCode, AuthCode}
			count, i, err = getStringPropParams(i, rdr)
			if err != nil {
				return nil, nil, err
			}
		case ReceiveMaxCode, TopicAliasMaxCode:
			validCodes = []int{ConnectCode, ConnackCode}
			count = 2
		case TopicAliasCode:
			validCodes = []int{PublishCode}
			count = 2
		case MaxQoSCode, RetainAvailableCode:
			validCodes = []int{ConnackCode}
			count = 1
		case UserPropertyCode:
			validCodes = []int{ConnectCode, ConnackCode, PublishCode, WillPropsCode, PubackCode, PubrecCode, PubrelCode, PubcompCode, SubscribeCode, SubackCode, UnsubscribeCode, UnsubackCode, DisconnectCode, AuthCode}
		case MaxPacketSizeCode:
			validCodes = []int{ConnectCode, ConnackCode}
			count = 4
		case WildcardSubAvailableCode, SubIdAvailableCode, SharedSubAvailableCode:
			validCodes = []int{ConnackCode}
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
			numRead, val, err := rdr.ReadVarByteInt()
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

// GetConnectFlags parses the given byte into flags for the connect packet.
func GetConnectFlags(b byte) (*ConnectFlags, error) {
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

// VerifyProtocol verifies that the following bytes from the reader represent the correct protocol. Hint: it must be MQTT.
// Assumes there are enough bytes to process the request.
func VerifyProtocol(rdr *packet.Reader) (err error) {
	_, s, err := rdr.ReadUtf8Str()
	if err != nil {
		return err
	}
	if s != "MQTT" {
		msg := fmt.Sprintf("Got invalid protocol `%v` expected `MQTT`", s)
		return errors.New(msg)
	}
	return nil
}

// GetKeepAlive reads the following two bytes and turns it into a 2 bytes integer
func GetKeepAlive(rdr *packet.Reader) (uint16, error) {
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
