package user

type User struct {
	ID       int    `db:"id"`
	UserID   string `db:"userID"`
	Fullname string `db:"fullname"`
}
