package user

// User represents a user.
type User struct {
	// ID is the user's unique identifier.
	ID string `json:"id" bson:"_id" db:"id"`
	// Username is the user's username.
	Username string `json:"username" bson:"username" db:"username"`
	// PasswordHash is the hashed version of the user's password.
	PasswordHash string `json:"password_hash" bson:"password_hash" db:"password_hash"`
}
