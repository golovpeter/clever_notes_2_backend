package sign_in

type SignIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	User_id  int    `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
}
