package signin

type SignIn struct {
	Username string
	Password string
}

type Out struct {
	Token string
}

type User struct {
	User_id  int    `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
}
