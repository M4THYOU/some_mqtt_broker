package utils

type Utf8Str = []byte
type Utf8StrPair = Utf8Str
type VariableByteInt = []uint8
type OneByteInt = uint8
type TwoByteInt = uint16
type FourByteInt = uint32
type BinaryData struct {
	length TwoByteInt
	data   []byte // of the specified length.
}
