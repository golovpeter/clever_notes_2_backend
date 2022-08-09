package signup

type SignUp struct {
	Email       string
	Password    string
	ConformPass string
}

type User struct {
	Email       string `db:"email"`
	Password    string `db:"password"`
	ConformPass string `db:"conform_pass"`
}
