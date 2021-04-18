package utils

import "github.com/M4THYOU/some_mqtt_broker/packet"

// IsIntInSlice checks if the provided integer is in the provided slice of integers.
func IsIntInSlice(i int, arr []int) bool {
	for _, v := range arr {
		if i == v {
			return true
		}
	}
	return false
}

// ReadBytesToSlice reads count bytes into a slice of bytes.
// On error, the current state of the slice, prior to error, is still returned.
func ReadBytesToSlice(count int, rdr *packet.Reader) ([]byte, error) {
	buf := make([]byte, 0)
	for i := 0; i < count; i++ {
		b, err := rdr.ReadByte()
		if err != nil {
			return buf, err
		}
		buf = append(buf, b)
	}
	return buf, nil
}

func Btoi(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}
