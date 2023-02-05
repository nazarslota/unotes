package user

type User struct {
	ID           string `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
}
