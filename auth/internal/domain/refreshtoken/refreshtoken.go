package refreshtoken

type RefreshToken struct {
	Token string `json:"token" bson:"token"`
}
