package config

const (
	Host          = "localhost"
	Port          = "1883"
	ConType       = "tcp"
	MaxPacketSize = 65536 // bytes
)

// Default values as defined in the spec.
const (
	DefaultSessionExpiryInterval = 0
	DefaultReceiveMaximum        = 65535 // max uint16
	DefaultTopicAliasMaximum     = 0
	DefaultRequestResponseInfo   = false
	DefaultRequestProblemInfo    = true
)
