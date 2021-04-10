package broker

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var _ = fmt.Printf // For debugging; delete when done.

const (
	connectFirstByte = byte(16) // 00010000
	connackFirstByte = byte(32) // 00100000

	// Publish Cases			  // byte: (DUP, QoS, Retain)
	publishFirstByte1  = byte(48) // 00110000: (0, 0, 0)
	publishFirstByte2  = byte(49) // 00110001: (0, 0, 1)
	publishFirstByte3  = byte(50) // 00110010: (0, 1, 0)
	publishFirstByte4  = byte(51) // 00110011: (0, 1, 1)
	publishFirstByte5  = byte(52) // 00110100: (0, 2, 0)
	publishFirstByte6  = byte(53) // 00110101: (0, 2, 1)
	publishFirstByte7  = byte(56) // 00111000: (1, 0, 0)
	publishFirstByte8  = byte(57) // 00111001: (1, 0, 1)
	publishFirstByte9  = byte(58) // 00111010: (1, 1, 0)
	publishFirstByte10 = byte(59) // 00111011: (1, 1, 1)
	publishFirstByte11 = byte(60) // 00111100: (1, 2, 0)
	publishFirstByte12 = byte(61) // 00111101: (1, 2, 1)

	pubackFirstByte      = byte(64)  // 01000000
	pubrecFirstByte      = byte(80)  // 01010000
	pubrelFirstByte      = byte(98)  // 01100010
	pubcompFirstByte     = byte(112) // 01110000
	subscribeFirstByte   = byte(130) // 10000010
	subackFirstByte      = byte(144) // 10010000
	unsubscribeFirstByte = byte(162) // 10100010
	unsubackFirstByte    = byte(176) // 10110000
	pingreqFirstByte     = byte(192) // 11000000
	pingrespFirstByte    = byte(208) // 11010000
	disconnectFirstByte  = byte(224) // 11100000
	authFirstByte        = byte(240) // 11110000
)

// func printRawBuffer(buf []byte, len int) {
// 	for i := 0; i < len; i++ {
// 		fmt.Printf("%d: %08b\n", i, buf[i])
// 	}
// }

func TestGetRequestType(t *testing.T) {
	b := getRequestType(connectFirstByte)
	if b == connackCode {
		t.Fatalf("Connect byte (%08b) and connackCode (%08b) should NOT be equal.", b, connackCode)
	}
	if b != connectCode {
		t.Fatalf("Connect bytes should be equal: %08b, %08b", b, connectCode)
	}

	b = getRequestType(connackFirstByte)
	if b != connackCode {
		t.Fatalf("Connack bytes should be equal: %08b, %08b", b, connackCode)
	}
	b = getRequestType(publishFirstByte1)
	if b != publishCode {
		t.Fatalf("Publish1 bytes should be equal: %08b, %08b", b, publishCode)
	}
	b = getRequestType(publishFirstByte2)
	if b != publishCode {
		t.Fatalf("Publish2 bytes should be equal: %08b, %08b", b, publishCode)
	}
	b = getRequestType(publishFirstByte3)
	if b != publishCode {
		t.Fatalf("Publish3 bytes should be equal: %08b, %08b", b, publishCode)
	}
	b = getRequestType(publishFirstByte4)
	if b != publishCode {
		t.Fatalf("Publish4 bytes should be equal: %08b, %08b", b, publishCode)
	}
	b = getRequestType(publishFirstByte5)
	if b != publishCode {
		t.Fatalf("Publish5 bytes should be equal: %08b, %08b", b, publishCode)
	}
	b = getRequestType(publishFirstByte6)
	if b != publishCode {
		t.Fatalf("Publish6 bytes should be equal: %08b, %08b", b, publishCode)
	}
	b = getRequestType(publishFirstByte7)
	if b != publishCode {
		t.Fatalf("Publish7 bytes should be equal: %08b, %08b", b, publishCode)
	}
	b = getRequestType(publishFirstByte8)
	if b != publishCode {
		t.Fatalf("Publish8 bytes should be equal: %08b, %08b", b, publishCode)
	}
	b = getRequestType(publishFirstByte9)
	if b != publishCode {
		t.Fatalf("Publish9 bytes should be equal: %08b, %08b", b, publishCode)
	}
	b = getRequestType(publishFirstByte10)
	if b != publishCode {
		t.Fatalf("Publish10 bytes should be equal: %08b, %08b", b, publishCode)
	}
	b = getRequestType(publishFirstByte11)
	if b != publishCode {
		t.Fatalf("Publish11 bytes should be equal: %08b, %08b", b, publishCode)
	}
	b = getRequestType(publishFirstByte12)
	if b != publishCode {
		t.Fatalf("Publish12 bytes should be equal: %08b, %08b", b, publishCode)
	}
	b = getRequestType(pubackFirstByte)
	if b != pubackCode {
		t.Fatalf("Puback bytes should be equal: %08b, %08b", b, pubackCode)
	}
	b = getRequestType(pubrecFirstByte)
	if b != pubrecCode {
		t.Fatalf("Pubrec bytes should be equal: %08b, %08b", b, pubrecCode)
	}
	b = getRequestType(pubrelFirstByte)
	if b != pubrelCode {
		t.Fatalf("Pubrel bytes should be equal: %08b, %08b", b, pubrelCode)
	}
	b = getRequestType(pubcompFirstByte)
	if b != pubcompCode {
		t.Fatalf("Pubcomp bytes should be equal: %08b, %08b", b, pubcompCode)
	}
	b = getRequestType(subscribeFirstByte)
	if b != subscribeCode {
		t.Fatalf("Subscribe bytes should be equal: %08b, %08b", b, subscribeCode)
	}
	b = getRequestType(subackFirstByte)
	if b != subackCode {
		t.Fatalf("Suback bytes should be equal: %08b, %08b", b, subackCode)
	}
	b = getRequestType(unsubscribeFirstByte)
	if b != unsubscribeCode {
		t.Fatalf("Unsubscribe bytes should be equal: %08b, %08b", b, unsubscribeCode)
	}
	b = getRequestType(unsubackFirstByte)
	if b != unsubackCode {
		t.Fatalf("Unsuback bytes should be equal: %08b, %08b", b, unsubackCode)
	}
	b = getRequestType(pingreqFirstByte)
	if b != pingreqCode {
		t.Fatalf("Pingreq bytes should be equal: %08b, %08b", b, pingreqCode)
	}
	b = getRequestType(pingrespFirstByte)
	if b != pingrespCode {
		t.Fatalf("Pingresp bytes should be equal: %08b, %08b", b, pingrespCode)
	}
	b = getRequestType(disconnectFirstByte)
	if b != disconnectCode {
		t.Fatalf("Disconnect bytes should be equal: %08b, %08b", b, disconnectCode)
	}
	b = getRequestType(authFirstByte)
	if b != authCode {
		t.Fatalf("Auth bytes should be equal: %08b, %08b", b, authCode)
	}
}

func checkSliceProtocol(t *testing.T, buf []byte, shouldPass bool) {
	rdr := bufio.NewReader(bytes.NewReader(buf))
	err := verifyProtocol(rdr)
	if err != nil && shouldPass {
		t.Fatalf("Invalid protocol: %v", err.Error())
	} else if err == nil && !shouldPass {
		t.Fatalf("Should have been an invalid protocol: %v", buf)
	}
}

func TestVerifyProtocol(t *testing.T) {
	buf := []byte{0, 4, 'M', 'Q', 'T', 'T'}
	checkSliceProtocol(t, buf, true)
	buf = []byte{0, 4, 'M', 'Q', 'T', 'T', 'T', 'T', 0x2d}
	checkSliceProtocol(t, buf, true)
	buf = []byte{0, 4, 'm', 'Q', 'T', 'T'}
	checkSliceProtocol(t, buf, false)
	buf = []byte{0, 5, 'M', 'Q', 'T', 'T', 'T'}
	checkSliceProtocol(t, buf, false)
	buf = []byte{0, 4, 'm', 'q', 't', 't'}
	checkSliceProtocol(t, buf, false)
	buf = []byte{0, 1, 'M', 'Q', 'T', 'T'}
	checkSliceProtocol(t, buf, false)
	buf = []byte{0, 0, 4, 'M', 'Q', 'T', 'T'} // expects first byte to be LSB of the protocol. SHOULD BE 4!
	checkSliceProtocol(t, buf, false)
	buf = []byte{1, 4, 'M', 'Q', 'T', 'T'}
	checkSliceProtocol(t, buf, false)
}

func checkConnectFlags(t *testing.T, b byte, expected *ConnectFlags, shouldPass bool) {
	flags, err := getConnectFlags(b)
	if err != nil && shouldPass {
		t.Fatalf("Invalid flags byte: %v", err.Error())
	} else if err == nil && !shouldPass {
		t.Fatalf("Should have been an invalid flags byte: %08b", b)
	} else if !cmp.Equal(flags, expected) && shouldPass {
		t.Fatalf("Got:\n%v\nExpected:\n%v", flags, expected)
	}
}
func TestGetConnectFlags(t *testing.T) {
	// Other rules:
	// If Will Flag == 0, Will QoS must be 0. Otherwise, Will QoS can be 0, 1, or 2
	// If Will Flag == 0, Will Retain must be 0. Otherwise, can be 0 or 1.

	// VALID
	checkConnectFlags(t, 0, &ConnectFlags{false, false, false, 0, false, false}, true) // 00000000: pass
	checkConnectFlags(t, 130, &ConnectFlags{true, false, false, 0, false, true}, true) // 10000010: pass
	checkConnectFlags(t, 134, &ConnectFlags{true, false, false, 0, true, true}, true)  // 10000110: pass
	checkConnectFlags(t, 70, &ConnectFlags{false, true, false, 0, true, true}, true)   // 01000110: pass
	checkConnectFlags(t, 44, &ConnectFlags{false, false, true, 1, true, false}, true)  // 00101100: pass
	checkConnectFlags(t, 52, &ConnectFlags{false, false, true, 2, true, false}, true)  // 00110100: pass
	checkConnectFlags(t, 214, &ConnectFlags{true, true, false, 2, true, true}, true)   // 11010110: pass
	// ERROR
	checkConnectFlags(t, 255, &ConnectFlags{true, true, true, 3, true, true}, false)    // 11111111: fail => invalid QoS && reserved bit is 1
	checkConnectFlags(t, 112, &ConnectFlags{false, true, true, 2, false, false}, false) // 01110000: fail => will flag is zero!
	checkConnectFlags(t, 120, &ConnectFlags{false, true, true, 3, false, false}, false) // 01111000: fail => invalid QoS && will flag is zero!
	checkConnectFlags(t, 124, &ConnectFlags{false, true, true, 3, true, false}, false)  // 01111100: fail => invalid QoS
	checkConnectFlags(t, 125, &ConnectFlags{false, true, true, 3, true, false}, false)  // 01111101: fail => invalid QoS && reserved bit is 1
	checkConnectFlags(t, 253, &ConnectFlags{true, true, true, 3, true, false}, false)   // 11111101: fail => invalid QoS && reserved bit is 1
	checkConnectFlags(t, 1, &ConnectFlags{false, false, false, 0, false, false}, false) // 00000001: fail => reserved bit is 1
}

func checkKeepAlive(t *testing.T, i, expected uint16) {
	// Convert the int to two bytes.
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, i)
	rdr := bufio.NewReader(bytes.NewReader(buf))
	// Get the value via testing!
	keepAlive, err := getKeepAlive(rdr)
	if err != nil {
		t.Fatalf("getKeepAlive failed: %v", err.Error())
	} else if keepAlive != expected {
		t.Fatalf("Got:\n%v\nExpected:\n%v", keepAlive, expected)
	}
}
func TestGetKeepAlive(t *testing.T) {
	var val uint16 = 4
	checkKeepAlive(t, val, val)
	val = 65535
	checkKeepAlive(t, val, val)
	val = 0
	checkKeepAlive(t, val, val)
}

func checkDecodeVarByteInt(t *testing.T, buf []byte, expected uint32, shouldPass bool) {
	// Convert the int to two bytes.
	rdr := bufio.NewReader(bytes.NewReader(buf))
	// Get the value via testing!
	val, err := decodeVarByteInt(rdr)
	if err != nil && shouldPass {
		t.Fatalf("decodeVarByteInt failed: %v", err.Error())
	} else if err == nil && !shouldPass {
		t.Fatalf("decodeVarByteInt should have failed: %v", val)
	} else if val != expected && shouldPass {
		t.Fatalf("Got:\n%v\nExpected:\n%v", val, expected)
	}
}
func TestDecodeVarByteInt(t *testing.T) {
	buf := []byte{0xFF, 0x64}
	var expected uint32 = 12927
	checkDecodeVarByteInt(t, buf, expected, true)
	buf = []byte{0x76}
	expected = 118
	checkDecodeVarByteInt(t, buf, expected, true)
	buf = []byte{0x7F}
	expected = 127
	checkDecodeVarByteInt(t, buf, expected, true)
	buf = []byte{0x80, 0x01}
	expected = 128
	checkDecodeVarByteInt(t, buf, expected, true)
	buf = []byte{0x00}
	expected = 0
	checkDecodeVarByteInt(t, buf, expected, true)
	buf = []byte{0x80, 0x80, 0x01}
	expected = 16384
	checkDecodeVarByteInt(t, buf, expected, true)
	buf = []byte{0xFF, 0xFF, 0x7F}
	expected = 2097151
	checkDecodeVarByteInt(t, buf, expected, true)
	buf = []byte{0x80, 0x80, 0x80, 0x01}
	expected = 2097152
	checkDecodeVarByteInt(t, buf, expected, true)
	buf = []byte{0xFF, 0xFF, 0xFF, 0x7F}
	expected = 268435455
	checkDecodeVarByteInt(t, buf, expected, true)
	buf = []byte{0xFF, 0xFF, 0xFF, 0x7F, 0x01}
	expected = 268435455
	checkDecodeVarByteInt(t, buf, expected, true)
	buf = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x01}
	expected = 268435456
	checkDecodeVarByteInt(t, buf, expected, false)
}
