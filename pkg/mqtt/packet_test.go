package mqtt

import (
	"testing"
)

func checkRequestType(t *testing.T, firstByte, expected byte, shouldPass bool) {
	b := GetRequestType(firstByte)
	if b != expected && shouldPass {
		t.Fatalf("GetRequestType: bytes should be equal: %08b, %08b\n", b, expected)
	} else if b == expected && !shouldPass {
		t.Fatalf("GetRequestType: bytes should NOT be equal: %08b, %08b\n", b, expected)
	}
}
func TestGetRequestType(t *testing.T) {
	checkRequestType(t, connectFirstByte, ConnackCode, false)
	checkRequestType(t, connectFirstByte, ConnectCode, true)
	checkRequestType(t, connackFirstByte, ConnackCode, true)
	checkRequestType(t, publishFirstByte1, PublishCode, true)
	checkRequestType(t, publishFirstByte2, PublishCode, true)
	checkRequestType(t, publishFirstByte3, PublishCode, true)
	checkRequestType(t, publishFirstByte4, PublishCode, true)
	checkRequestType(t, publishFirstByte5, PublishCode, true)
	checkRequestType(t, publishFirstByte6, PublishCode, true)
	checkRequestType(t, publishFirstByte7, PublishCode, true)
	checkRequestType(t, publishFirstByte8, PublishCode, true)
	checkRequestType(t, publishFirstByte9, PublishCode, true)
	checkRequestType(t, publishFirstByte10, PublishCode, true)
	checkRequestType(t, publishFirstByte11, PublishCode, true)
	checkRequestType(t, publishFirstByte12, PublishCode, true)
	checkRequestType(t, pubackFirstByte, PubackCode, true)
	checkRequestType(t, pubrecFirstByte, PubrecCode, true)
	checkRequestType(t, pubrelFirstByte, PubrelCode, true)
	checkRequestType(t, pubcompFirstByte, PubcompCode, true)
	checkRequestType(t, subscribeFirstByte, SubscribeCode, true)
	checkRequestType(t, subackFirstByte, SubackCode, true)
	checkRequestType(t, unsubscribeFirstByte, UnsubscribeCode, true)
	checkRequestType(t, unsubackFirstByte, UnsubackCode, true)
	checkRequestType(t, pingreqFirstByte, PingreqCode, true)
	checkRequestType(t, pingrespFirstByte, PingrespCode, true)
	checkRequestType(t, disconnectFirstByte, DisconnectCode, true)
	checkRequestType(t, authFirstByte, AuthCode, true)
}

func checkSetRequestType(t *testing.T, packetCode uint8, expected byte, dupFlag, retainFlag bool, qos int, shouldPass bool) {
	b := SetRequestType(packetCode, dupFlag, retainFlag, qos)
	if b != expected && shouldPass {
		t.Fatalf("SetRequestType: bytes should be equal: %08b, %08b\n", b, expected)
	} else if b == expected && !shouldPass {
		t.Fatalf("SetRequestType: bytes should NOT be equal: %08b, %08b\n", b, expected)
	}
}
func TestSetRequestType(t *testing.T) {
	checkSetRequestType(t, ConnectCode, connackFirstByte, false, false, 0, false)
	checkSetRequestType(t, ConnectCode, connectFirstByte, false, false, 0, true)
	checkSetRequestType(t, ConnackCode, connackFirstByte, false, false, 0, true)
	checkSetRequestType(t, PublishCode, publishFirstByte1, false, false, 0, true)
	checkSetRequestType(t, PublishCode, publishFirstByte2, false, true, 0, true)
	checkSetRequestType(t, PublishCode, publishFirstByte3, false, false, 1, true)
	checkSetRequestType(t, PublishCode, publishFirstByte4, false, true, 1, true)
	checkSetRequestType(t, PublishCode, publishFirstByte5, false, false, 2, true)
	checkSetRequestType(t, PublishCode, publishFirstByte6, false, true, 2, true)
	checkSetRequestType(t, PublishCode, publishFirstByte7, true, false, 0, true)
	checkSetRequestType(t, PublishCode, publishFirstByte8, true, true, 0, true)
	checkSetRequestType(t, PublishCode, publishFirstByte9, true, false, 1, true)
	checkSetRequestType(t, PublishCode, publishFirstByte10, true, true, 1, true)
	checkSetRequestType(t, PublishCode, publishFirstByte11, true, false, 2, true)
	checkSetRequestType(t, PublishCode, publishFirstByte12, true, true, 2, true)
	checkSetRequestType(t, PubackCode, pubackFirstByte, false, false, 0, true)
	checkSetRequestType(t, PubrecCode, pubrecFirstByte, false, false, 0, true)
	checkSetRequestType(t, PubrelCode, pubrelFirstByte, false, false, 0, true)
	checkSetRequestType(t, PubcompCode, pubcompFirstByte, false, false, 0, true)
	checkSetRequestType(t, SubscribeCode, subscribeFirstByte, false, false, 0, true)
	checkSetRequestType(t, SubackCode, subackFirstByte, false, false, 0, true)
	checkSetRequestType(t, UnsubscribeCode, unsubscribeFirstByte, false, false, 0, true)
	checkSetRequestType(t, UnsubackCode, unsubackFirstByte, false, false, 0, true)
	checkSetRequestType(t, PingreqCode, pingreqFirstByte, false, false, 0, true)
	checkSetRequestType(t, PingrespCode, pingrespFirstByte, false, false, 0, true)
	checkSetRequestType(t, DisconnectCode, disconnectFirstByte, false, false, 0, true)
	checkSetRequestType(t, AuthCode, authFirstByte, false, false, 0, true)
}
