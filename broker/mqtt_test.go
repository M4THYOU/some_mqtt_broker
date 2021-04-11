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

func checkDecodeVarByteInt(t *testing.T, buf []byte, expectedByteCount int, expected uint32, shouldPass bool) {
	// Convert the int to two bytes.
	rdr := bufio.NewReader(bytes.NewReader(buf))
	// Get the value via testing!
	i, val, err := decodeVarByteInt(rdr)
	if err != nil && shouldPass {
		t.Fatalf("decodeVarByteInt failed: %v", err.Error())
	} else if err == nil && !shouldPass {
		t.Fatalf("decodeVarByteInt should have failed: %v", val)
	} else if val != expected && shouldPass {
		t.Fatalf("Got:\n%v\nExpected:\n%v", val, expected)
	} else if i != expectedByteCount && shouldPass {
		t.Fatalf("read %d bytes, expected %d", i, expectedByteCount)
	}
}
func TestDecodeVarByteInt(t *testing.T) {
	buf := []byte{0xFF, 0x64}
	var expected uint32 = 12927
	checkDecodeVarByteInt(t, buf, 2, expected, true)
	buf = []byte{0x76}
	expected = 118
	checkDecodeVarByteInt(t, buf, 1, expected, true)
	buf = []byte{0x7F}
	expected = 127
	checkDecodeVarByteInt(t, buf, 1, expected, true)
	buf = []byte{0x80, 0x01}
	expected = 128
	checkDecodeVarByteInt(t, buf, 2, expected, true)
	buf = []byte{0x00}
	expected = 0
	checkDecodeVarByteInt(t, buf, 1, expected, true)
	buf = []byte{0x80, 0x80, 0x01}
	expected = 16384
	checkDecodeVarByteInt(t, buf, 3, expected, true)
	buf = []byte{0xFF, 0xFF, 0x7F}
	expected = 2097151
	checkDecodeVarByteInt(t, buf, 3, expected, true)
	buf = []byte{0x80, 0x80, 0x80, 0x01}
	expected = 2097152
	checkDecodeVarByteInt(t, buf, 4, expected, true)
	buf = []byte{0xFF, 0xFF, 0xFF, 0x7F}
	expected = 268435455
	checkDecodeVarByteInt(t, buf, 4, expected, true)
	buf = []byte{0xFF, 0xFF, 0xFF, 0x7F, 0x01}
	expected = 268435455
	checkDecodeVarByteInt(t, buf, 4, expected, true)
	buf = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x01}
	expected = 268435456
	checkDecodeVarByteInt(t, buf, 5, expected, false)
}

func checkStringPropParams(t *testing.T, i int, buf []byte, expectedCount int, shouldPass bool) {
	expectedI := i + 2 // since the function reads 2 bytes, the newI should be 2 greater than i.
	rdr := bufio.NewReader(bytes.NewReader(buf))
	count, newI, err := getStringPropParams(i, rdr)
	if err != nil && shouldPass {
		t.Fatalf("getStringPropParams failed: %v", err.Error())
	} else if err == nil && !shouldPass {
		t.Fatalf("getStringPropParams should have failed: %v", buf)
	} else if newI != expectedI && shouldPass {
		t.Fatalf("incorrect i. Got:\n%v\nExpected:\n%v", newI, expectedI)
	} else if count != expectedCount && shouldPass {
		t.Fatalf("incorrect count. Got:\n%v\nExpected:\n%v", count, expectedCount)
	}
}
func TestGetStringPropParams(t *testing.T) {
	example := []byte{0x00, 0x05, 0x41, 0xF0, 0xAA, 0x9B, 0x94}
	expectedCount := 5
	checkStringPropParams(t, 99, example, expectedCount, true)
	checkStringPropParams(t, 2, example, expectedCount, true)
	example = []byte{0x00, 0x04, 0x41, 0xF0, 0xAA, 0x9B, 0x94}
	expectedCount = 4
	checkStringPropParams(t, 13, example, expectedCount, true)
	example = []byte{0x00, 0x00}
	expectedCount = 0
	checkStringPropParams(t, 7, example, expectedCount, true)
	example = []byte{0x00}
	checkStringPropParams(t, 7, example, expectedCount, false)
	example = []byte{0x05}
	checkStringPropParams(t, 7, example, expectedCount, false)
}

func checkStringPairProp(t *testing.T, buf, expected []byte, expectedRead int, shouldPass bool) {
	rdr := bufio.NewReader(bytes.NewReader(buf))
	numRead, res, err := getStringPairProp(rdr)
	if err != nil && shouldPass {
		t.Fatalf("getStringPropParams failed: %v", err.Error())
	} else if err == nil && !shouldPass {
		t.Fatalf("getStringPropParams should have failed: %v", buf)
	} else if !cmp.Equal(res, expected) && shouldPass {
		t.Fatalf("Got:\n%v\nExpected:\n%v", res, expected)
	} else if numRead != expectedRead && shouldPass {
		t.Fatalf("incorrect numRead. Got:\n%v\nExpected:\n%v", numRead, expectedRead)
	}
}
func TestGetStringPairProp(t *testing.T) {
	buf := []byte{0x00, 0x05, 0x41, 0xF0, 0xAA, 0x9B, 0x94, 0x00, 0x05, 0x41, 0xF0, 0xAA, 0x9B, 0x94}
	expected := buf
	expectedRead := 14
	checkStringPairProp(t, buf, expected, expectedRead, true)
	buf = []byte{0x00, 0x05, 0x41, 0xF0, 0xAA, 0x9B, 0x94}
	checkStringPairProp(t, buf, expected, expectedRead, false)
	buf = []byte{0x00, 0x05, 0x41, 0xF0, 0xAA, 0x9B, 0x94, 0x00, 0x05}
	checkStringPairProp(t, buf, expected, expectedRead, false)
	buf = []byte{0x00, 0x05, 0x41, 0xF0, 0xAA, 0x9B, 0x94, 0x00}
	checkStringPairProp(t, buf, expected, expectedRead, false)
	buf = []byte{0x00, 0x04, 0xF0, 0xAA, 0x9B, 0x94, 0x00, 0x05, 0x41, 0xF0, 0xAA, 0x9B, 0x94, 0x05}
	expected = []byte{0x00, 0x04, 0xF0, 0xAA, 0x9B, 0x94, 0x00, 0x05, 0x41, 0xF0, 0xAA, 0x9B, 0x94}
	expectedRead = 13
	checkStringPairProp(t, buf, expected, expectedRead, true)
	buf = []byte{0xF0, 0x00, 0x04, 0xF0, 0xAA, 0x9B, 0x94, 0x00, 0x05, 0x41, 0xF0, 0xAA, 0x9B, 0x94, 0x05}
	checkStringPairProp(t, buf, expected, expectedRead, false)
	buf = []byte{0x00, 0x02, 0xF0, 0xF0, 0x00, 0x00, 0x41, 0xF0, 0xAA, 0x9B, 0x94, 0x05}
	expected = []byte{0x00, 0x02, 0xF0, 0xF0, 0x00, 0x00}
	expectedRead = 6
	checkStringPairProp(t, buf, expected, expectedRead, true)
}

func checkClientId(t *testing.T, buf []byte, expected string, shouldPass bool) {
	rdr := bufio.NewReader(bytes.NewReader(buf))
	clientId, err := getClientId(rdr)
	if err != nil && shouldPass {
		t.Fatalf("getClientId failed: %v", err.Error())
	} else if err == nil && !shouldPass {
		t.Fatalf("getClientId should have failed: %v", buf)
	} else if clientId != expected && shouldPass {
		t.Fatalf("Got:\n%v\nExpected:\n%v", clientId, expected)
	}
}
func TestGetClientId(t *testing.T) {
	buf := []byte{0x00, 0x00}
	expected := "one randomly"
	checkClientId(t, buf, expected, true)
	buf = []byte{0x00, 0x01, 0x32}
	expected = "2"
	checkClientId(t, buf, expected, true)
	buf = []byte{0x00, 0x02, 0x32, 0x61}
	expected = "2a"
	checkClientId(t, buf, expected, true)
	buf = []byte{0x00}
	checkClientId(t, buf, expected, false)
}

// Below are all the tests for the getProps function. They are split into 15 different functions, 1 for each packet type.
// Each one tests on every property identifier at least once, plus some extra cases that may be unique to that packet.
var (
	payloadFormatIndicator []byte = []byte{0x01, 0x01}                                                                         // 1
	messageExpiryInterval  []byte = []byte{0x02, 0x00, 0x00, 0x00, 0x3C}                                                       // 60
	contentType            []byte = []byte{0x03, 0x00, 0x04, 0x6a, 0x73, 0x6F, 0x6E}                                           // "json"
	responseTopic          []byte = []byte{0x08, 0x00, 0x0A, 0x73, 0x6f, 0x6d, 0x65, 0x2f, 0x74, 0x6f, 0x70, 0x69, 0x63}       // "some/topic"
	correlationData        []byte = []byte{0x09, 0x00, 0x04, 0x00, 0x00, 0x01, 0x00}                                           // 4 useless bytes of data.
	subscriptionId         []byte = []byte{0x0B, 0x01}                                                                         // 1
	sessionExpiryInterval  []byte = []byte{0x11, 0x00, 0x00, 0x00, 0x3C}                                                       // 60
	assignedClientId       []byte = []byte{0x12, 0x00, 0x03, 0x6f, 0x6e, 0x65}                                                 // "one"
	serverKeepAlive        []byte = []byte{0x13, 0x01, 0x90}                                                                   // 400
	authenticationMethod   []byte = []byte{0x15, 0x00, 0x0B, 0x53, 0x43, 0x52, 0x41, 0x4d, 0x2d, 0x53, 0x48, 0x41, 0x2d, 0x31} // "SCRAM-SHA-1"
	authenticationData     []byte = []byte{0x16, 0x00, 0x02, 0x03, 0x0FF}                                                      // 2 useless bytes of data
	requestProblemInfo     []byte = []byte{0x17, 0x01}                                                                         // 1
	willDelayInterval      []byte = []byte{0x18, 0x00, 0x00, 0x00, 0x3C}                                                       // 60
	requestResponseInfo    []byte = []byte{0x19, 0x01}                                                                         // 1
	responseInfo           []byte = []byte{0x1A, 0x00, 0x0B, 0x73, 0x6f, 0x6d, 0x65, 0x20, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67} // "some string"
	serverReference        []byte = []byte{0x1C, 0x00, 0x0B, 0x31, 0x39, 0x32, 0x2e, 0x31, 0x36, 0x38, 0x2e, 0x32, 0x2e, 0x31} // "192.168.2.1"
	reasonString           []byte = []byte{0x1F, 0x00, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e}                               // "reason"
	receiveMax             []byte = []byte{0x21, 0x00, 0x63}                                                                   // 99
	topicAliasMax          []byte = []byte{0x22, 0x01, 0x2D}                                                                   // 301
	topicAlias             []byte = []byte{0x23, 0x00, 0x05}                                                                   // 5
	maxQoS                 []byte = []byte{0x24, 0x01}                                                                         // 1
	retainAvailable        []byte = []byte{0x25, 0x01}                                                                         // 1
	userProperty           []byte = []byte{0x26, 0x00, 0x01, 0x41, 0x00, 0x02, 0x63, 0x64}                                     // "A" "cd"
	userProperty2          []byte = []byte{0x26, 0x00, 0x02, 0x63, 0x64, 0x00, 0x01, 0x41}                                     // "cd" "A"
	maxPacketSize          []byte = []byte{0x27, 0xFF, 0xFF, 0xFF, 0xFF}                                                       // 4294967296
	wildcardSubAvailable   []byte = []byte{0x28, 0x01}                                                                         // 1
	subIdAvailable         []byte = []byte{0x29, 0x00}                                                                         // 0
	sharedSubAvailable     []byte = []byte{0x2A, 0x01}                                                                         // 1
)

func checkProps(t *testing.T, propLen, packetCode int, buf []byte, expectedM map[int][]byte, expectedUserProps [][]byte, shouldPass bool) {
	rdr := bufio.NewReader(bytes.NewReader(buf))
	m, userProps, err := getProps(rdr, propLen, packetCode)
	if err != nil && shouldPass {
		t.Fatalf("getProps failed: %v", err.Error())
	} else if err == nil && !shouldPass {
		t.Fatalf("getProps should have failed: %v", buf)
	} else if !cmp.Equal(m, expectedM) && shouldPass {
		t.Fatalf("incorrect map Got:\n%v\nExpected:\n%v", m, expectedM)
	} else if !cmp.Equal(userProps, expectedUserProps) && shouldPass {
		t.Fatalf("incorrect userProps. Got:\n%v\nExpected:\n%v", userProps, expectedUserProps)
	}
}
func TestGetPropsCONNECT(t *testing.T) {
	// The basics
	packetCode := connectCode
	checkProps(t, 2, packetCode, payloadFormatIndicator, nil, nil, false)
	checkProps(t, 5, packetCode, messageExpiryInterval, nil, nil, false)
	checkProps(t, 5, packetCode, contentType, nil, nil, false)
	checkProps(t, 11, packetCode, responseTopic, nil, nil, false)
	checkProps(t, 7, packetCode, correlationData, nil, nil, false)
	checkProps(t, 2, packetCode, subscriptionId, nil, nil, false)
	expected := map[int][]byte{0x11: sessionExpiryInterval[1:]}
	checkProps(t, 2, packetCode, sessionExpiryInterval, expected, [][]byte{}, true)
	checkProps(t, 4, packetCode, assignedClientId, nil, nil, false)
	checkProps(t, 3, packetCode, serverKeepAlive, nil, nil, false)
	expected = map[int][]byte{0x15: authenticationMethod[3:]}
	checkProps(t, 14, packetCode, authenticationMethod, expected, [][]byte{}, true)
	expected = map[int][]byte{0x16: authenticationData[3:]}
	checkProps(t, 5, packetCode, authenticationData, expected, [][]byte{}, true)
	expected = map[int][]byte{0x17: requestProblemInfo[1:]}
	checkProps(t, 2, packetCode, requestProblemInfo, expected, [][]byte{}, true)
	checkProps(t, 5, packetCode, willDelayInterval, nil, nil, false)
	expected = map[int][]byte{0x19: requestResponseInfo[1:]}
	checkProps(t, 2, packetCode, requestResponseInfo, expected, [][]byte{}, true)
	checkProps(t, 14, packetCode, responseInfo, nil, nil, false)
	checkProps(t, 14, packetCode, serverReference, nil, nil, false)
	checkProps(t, 9, packetCode, reasonString, nil, nil, false)
	expected = map[int][]byte{0x21: receiveMax[1:]}
	checkProps(t, 3, packetCode, receiveMax, expected, [][]byte{}, true)
	expected = map[int][]byte{0x22: topicAliasMax[1:]}
	checkProps(t, 3, packetCode, topicAliasMax, expected, [][]byte{}, true)
	checkProps(t, 3, packetCode, topicAlias, nil, nil, false)
	checkProps(t, 2, packetCode, maxQoS, nil, nil, false)
	checkProps(t, 2, packetCode, retainAvailable, nil, nil, false)
	expected = map[int][]byte{}
	expectedUProps := [][]byte{userProperty[1:]}
	checkProps(t, 8, packetCode, userProperty, expected, expectedUProps, true)
	expected = map[int][]byte{0x27: maxPacketSize[1:]}
	checkProps(t, 5, packetCode, maxPacketSize, expected, [][]byte{}, true)
	checkProps(t, 2, packetCode, wildcardSubAvailable, nil, nil, false)
	checkProps(t, 2, packetCode, subIdAvailable, nil, nil, false)
	checkProps(t, 2, packetCode, sharedSubAvailable, nil, nil, false)
	// Special ones.
	// multiple userProps, using the same one.
	payload := append(userProperty, userProperty...)
	expected = map[int][]byte{}
	expectedUProps = [][]byte{userProperty[1:], userProperty[1:]}
	checkProps(t, 16, packetCode, payload, expected, expectedUProps, true)
	// multiple userProps, using different ones.
	payload = append(userProperty, userProperty2...)
	expected = map[int][]byte{}
	expectedUProps = [][]byte{userProperty[1:], userProperty2[1:]}
	checkProps(t, 16, packetCode, payload, expected, expectedUProps, true)
	// multiple different valid properties
	payload = append(authenticationData, append(sessionExpiryInterval, authenticationMethod...)...)
	expected = map[int][]byte{0x15: authenticationMethod[3:], 0x16: authenticationData[3:], 0x11: sessionExpiryInterval[1:]}
	checkProps(t, 24, packetCode, payload, expected, [][]byte{}, true)
	// multiple different properties with one invalid
	payload = append(authenticationData, append(sessionExpiryInterval, append(authenticationMethod, serverKeepAlive...)...)...)
	checkProps(t, 27, packetCode, payload, nil, nil, false)
	// userProps + other props in some random order.
	payload = append(authenticationData, append(userProperty, append(sessionExpiryInterval, append(userProperty2, authenticationMethod...)...)...)...)
	expected = map[int][]byte{0x15: authenticationMethod[3:], 0x16: authenticationData[3:], 0x11: sessionExpiryInterval[1:]}
	expectedUProps = [][]byte{userProperty[1:], userProperty2[1:]}
	checkProps(t, 40, packetCode, payload, expected, expectedUProps, true)
	// Empty props
	checkProps(t, 0, packetCode, []byte{}, map[int][]byte{}, [][]byte{}, true)
}
