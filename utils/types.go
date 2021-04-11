package utils

type Utf8StringPair struct {
	name  []byte // up to 65,535 bytes.
	value []byte // up to 65,535 bytes.
}
type BinaryData struct {
	length int16
	data   []byte
}
