package packet

import (
	"bytes"
	"testing"
)

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
