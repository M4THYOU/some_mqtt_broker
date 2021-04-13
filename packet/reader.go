package packet

import (
	"bufio"
	"errors"
	"io"
)

// bytesRead <= remainingLength, always.
type Reader struct {
	rdr             *bufio.Reader
	remainingLength int
}

// NewReader returns a new Reader whose buffer has the default size.
func NewReader(rdr io.Reader, remainingLength int) *Reader {
	return &Reader{bufio.NewReader(rdr), remainingLength}
}

// SetRemainingLength sets the remaining length of the packet currently being read.
func (rdr *Reader) SetRemainingLength(v int) {
	rdr.remainingLength = v
}

// ReadByte reads and returns a single byte.
// If no byte is available, returns an error.
func (rdr *Reader) ReadByte() (byte, error) {
	if rdr.remainingLength == 0 {
		return 0, errors.New("packet.Reader ReadByte failed to read byte: remainingLength == 0")
	}
	b, err := rdr.rdr.ReadByte()
	if err != nil {
		return 0, err
	}
	rdr.remainingLength--
	return b, nil
}

// ReadVarByteInt returns the integer value of a decoded Variable Byte Int according to MQTT v5.0 Spec.
// Returns number of bytes read, the integer, and possibly an error.
func (rdr *Reader) ReadVarByteInt() (int, uint32, error) {
	var multiplier, val uint32 = 1, 0
	var b byte
	var err error
	bytesRead := 0
	for {
		b, err = rdr.ReadByte()
		bytesRead++
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
	return bytesRead, val, nil
}
