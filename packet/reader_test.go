package packet

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func readBytes(rdr *Reader, buf []byte, remainingLength, bytesToRead int, shouldPass bool) error {
	for i := 0; i < bytesToRead; i++ {
		b, err := rdr.ReadByte()
		fmt.Println(b, err)
		if err != nil && shouldPass {
			msg := fmt.Sprintf("ReadByte failed: %v", err.Error())
			return errors.New(msg)
		} else if b != buf[i] && shouldPass {
			msg := fmt.Sprintf("read %08b, expected %08b", b, buf[i])
			return errors.New(msg)
		} else if err != nil { // this is good => err != nil && !shouldPass. Which means the test actually does pass
			return nil
		}
	}
	if !shouldPass {
		msg := fmt.Sprintf("ReadVarByteInt should have failed.\nBuf: %v\nremainingLength: %v\nbytesToRead: %v", buf, remainingLength, bytesToRead)
		return errors.New(msg)
	}
	return nil
}

func checkReadByte(t *testing.T, buf []byte, remainingLength, bytesToRead int, shouldPass bool) {
	rdr := NewReader(bytes.NewReader(buf), remainingLength)
	err := readBytes(rdr, buf, remainingLength, bytesToRead, shouldPass)
	if err != nil {
		t.Fatal(err)
	}
}
func TestReadByte(t *testing.T) {
	// Cases => Make sure these are all within the range of buf. Otherwise, we're testing io.Reader. That's stupid:
	buf := []byte{4, 6, 'w', 0x32, 0, '.', 123, 'W'}
	checkReadByte(t, buf, 8, 8, true)
	checkReadByte(t, buf, 3, 3, true)
	checkReadByte(t, buf, 7, 3, true)
	checkReadByte(t, buf, 0, 0, true)
	checkReadByte(t, buf, 2, 6, false)
	checkReadByte(t, buf, 0, 1, false)
}

func checkRemainingLength(t *testing.T, buf []byte, initRemLen, readX1, remLen, readX2 int, shouldPass1, shouldPass2 bool) {
	rdr := NewReader(bytes.NewReader(buf), initRemLen)
	if rdr.remainingLength != initRemLen {
		t.Fatalf("remainingLength got %d, expected %d", rdr.remainingLength, initRemLen)
	}
	err := readBytes(rdr, buf, initRemLen, readX1, shouldPass1)
	if err != nil {
		t.Fatal(err)
	}

	if !shouldPass1 {
		return
	}

	rdr.SetRemainingLength(remLen)
	if rdr.remainingLength != remLen {
		t.Fatalf("remainingLength got %d, expected %d", rdr.remainingLength, remLen)
	}
	err = readBytes(rdr, buf[readX1:], initRemLen, readX2, shouldPass2)
	if err != nil {
		t.Fatal(err)
	}
}
func TestSetRemainingLength(t *testing.T) {
	buf := []byte{1, 2, 3, 4, 5, 6, 'A', 'B', 2, 'D', 0x01, 'y', '1', 7, 8, 9, 10, 111}
	checkRemainingLength(t, buf, 2, 2, 4, 4, true, true)
	checkRemainingLength(t, buf, 2, 3, 4, 4, false, true)
	checkRemainingLength(t, buf, 2, 2, 2, 5, true, false)
	checkRemainingLength(t, buf, 2, 2, 11, 5, true, true)
}

func checkReadVarByteInt(t *testing.T, buf []byte, expectedByteCount int, expected uint32, shouldPass bool) {
	rdr := NewReader(bytes.NewReader(buf), 50)
	// Get the value via testing!
	i, val, err := rdr.ReadVarByteInt()
	if err != nil && shouldPass {
		t.Fatalf("ReadVarByteInt failed: %v", err.Error())
	} else if err == nil && !shouldPass {
		t.Fatalf("ReadVarByteInt should have failed: %v", val)
	} else if val != expected && shouldPass {
		t.Fatalf("Got:\n%v\nExpected:\n%v", val, expected)
	} else if i != expectedByteCount && shouldPass {
		t.Fatalf("read %d bytes, expected %d", i, expectedByteCount)
	}
}
func TestReadVarByteInt(t *testing.T) {
	buf := []byte{0xFF, 0x64}
	var expected uint32 = 12927
	checkReadVarByteInt(t, buf, 2, expected, true)
	buf = []byte{0x76}
	expected = 118
	checkReadVarByteInt(t, buf, 1, expected, true)
	buf = []byte{0x7F}
	expected = 127
	checkReadVarByteInt(t, buf, 1, expected, true)
	buf = []byte{0x80, 0x01}
	expected = 128
	checkReadVarByteInt(t, buf, 2, expected, true)
	buf = []byte{0x00}
	expected = 0
	checkReadVarByteInt(t, buf, 1, expected, true)
	buf = []byte{0x80, 0x80, 0x01}
	expected = 16384
	checkReadVarByteInt(t, buf, 3, expected, true)
	buf = []byte{0xFF, 0xFF, 0x7F}
	expected = 2097151
	checkReadVarByteInt(t, buf, 3, expected, true)
	buf = []byte{0x80, 0x80, 0x80, 0x01}
	expected = 2097152
	checkReadVarByteInt(t, buf, 4, expected, true)
	buf = []byte{0xFF, 0xFF, 0xFF, 0x7F}
	expected = 268435455
	checkReadVarByteInt(t, buf, 4, expected, true)
	buf = []byte{0xFF, 0xFF, 0xFF, 0x7F, 0x01}
	expected = 268435455
	checkReadVarByteInt(t, buf, 4, expected, true)
	buf = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x01}
	expected = 268435456
	checkReadVarByteInt(t, buf, 5, expected, false)
}
