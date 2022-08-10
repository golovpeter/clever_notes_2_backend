package signup

type SignUp struct {
	Username string
	Password string
}

type User struct {
	Username string `db:"email"`
	Password string `db:"password"`
}
