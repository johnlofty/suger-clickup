package models

type User struct {
	UserId    int32
	Username  string
	Password  string
	Email     string
	CreatedAt int64
	OrgId     int32
}

func CreateUser(u *User) error {
	_, err := db.Exec("INSERT INTO users(username, password) VALUES ($1, $2)")
	return err
}

func GetUser(email string) (User, error) {
	// user := User{}

	// query := `SELECT * FROM users WHERE email=$1`

	// err :=
	// return
	return User{}, nil
}
