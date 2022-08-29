package user

type User struct {
	ID           string `json:"id" bson:"_id"`
	Username     string `json:"username" bson:"username"`
	PasswordHash string `json:"password_hash" bson:"password_hash"`
}
