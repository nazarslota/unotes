package refreshtoken

// Token represents a token.
type Token struct {
	// Token is the token string.
	Token string `json:"token" bson:"token"`
}
