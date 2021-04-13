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
