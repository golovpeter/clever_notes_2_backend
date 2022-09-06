package sign_in

type SignIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInOut struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	User_id  int    `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
}
