package sign_up

type SignUpIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username string `db:"email"`
	Password string `db:"password"`
}
