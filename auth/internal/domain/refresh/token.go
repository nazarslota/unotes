package refresh

import "encoding"

type Token string

var _ encoding.BinaryMarshaler = (*Token)(nil)

func (t Token) MarshalBinary() (data []byte, err error) {
	return []byte(t), nil
}
