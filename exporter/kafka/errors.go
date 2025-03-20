package kafka

import "errors"

var (
	NoSASLMechanismFound = errors.New("invalid SHA algorithm: can be either \"sha256\" or \"sha512\"")
)
