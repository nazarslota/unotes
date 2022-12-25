package user

type User struct {
	ID           string `json:"id" bson:"_id" db:"id"`
	Username     string `json:"username" bson:"username" db:"username"`
	PasswordHash string `json:"password_hash" bson:"password_hash" db:"password_hash"`
}
