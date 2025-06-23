package kafka

import "errors"

var (
	ErrNoSASLMechanismFound = errors.New("invalid SHA algorithm: can be either \"sha256\" or \"sha512\"")
)
