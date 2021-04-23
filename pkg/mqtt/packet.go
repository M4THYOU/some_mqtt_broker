package mqtt

import (
	"errors"
	"fmt"
	"net"

	"github.com/M4THYOU/some_mqtt_broker/pkg/utils"
)

// GetRequestType converts the given byte into another byte of the appropriate request type format.
func GetRequestType(b byte) byte {
	return (b & 0xF0) >> 4
}

// SetRequestType creates the first byte of an MQTT packet from the provided info.
// dupFlag, retainFlag, and qos are only used when packetCode is the PUBLISH code
func SetRequestType(packetCode uint8, dupFlag, retainFlag bool, qos int) byte {
	b := (packetCode & 0x0F) << 4
	if packetCode == PublishCode {
		bDupFlag := uint8(utils.Btoi(dupFlag)) << 3
		bRetainFlag := uint8(utils.Btoi(retainFlag))
		bQos := uint8(qos * 2)
		b = b | bDupFlag | bRetainFlag | bQos
	} else if packetCode == PubrelCode || packetCode == SubscribeCode || packetCode == UnsubscribeCode {
		b = (b | 0x02)
	}
	return b
}

func SendPacket(conn net.Conn, packet []byte) error {
	n, err := conn.Write(packet)
	packetLen := len(packet)
	if err != nil {
		return err
	} else if n != packetLen {
		msg := fmt.Sprintf("Connection did not write the expected length. expected %v got %v", packetLen, n)
		return errors.New(msg)
	}
	return nil
}
