package refresh

import (
	"encoding"
	"errors"
)

type Token string

var _ encoding.BinaryMarshaler = (*Token)(nil)

func (t Token) MarshalBinary() ([]byte, error) { return []byte(t), nil }

var ErrTokenNotFound = errors.New("token not found")
