package note

type Note struct {
	ID      string `bson:"_id"`
	Title   string `bson:"title"`
	Content string `bson:"content"`
	UserID  string `bson:"user_id"`
}
