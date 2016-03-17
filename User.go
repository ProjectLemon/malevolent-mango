package main

//User provides a datatype for gathering the user information
//for the current user.
type User struct {
	Email    string
	FullName string
	Password string
	Salt     string
}

func (u *User) InDatabase() bool {
	return (u.FullName != "" && u.Password != "" && u.Salt != "")
}
