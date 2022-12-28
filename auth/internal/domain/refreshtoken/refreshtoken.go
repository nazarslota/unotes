package refreshtoken

import "errors"

type Token struct {
	Token string `json:"token" bson:"token"`
}

var (
	ErrTokenNotFound  = errors.New("token not found")
	ErrTokensNotFound = errors.New("tokens not found")
)
