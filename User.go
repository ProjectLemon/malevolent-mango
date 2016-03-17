package main

//User provides a datatype for gathering the user information
//for the current user.
type User struct {
	Email    string
	FullName string
	Password string
	Salt     string
}

//NewUser creates a new instance of a user type
//assigning the provided email address to the address field
//and zero value for other fields
func NewUser(email string) *User {
	user := new(User)
	user.Email = email
	user.FullName = ""
	user.Password = ""
	user.Salt = ""
	return user
}

func (u *User) InDatabase() bool {
	return (u.FullName != "" && u.Password != "" && u.Salt != "")
}
