package main

import "database/sql"

//User provides a datatype for gathering the user information
//for the current user.
type User struct {
	Email        sql.NullString
	FullName     sql.NullString
	PasswordHash sql.NullString
	Salt         sql.NullString
}

//NewUser creates a new instance of a user type
//assigning the provided email address to the address field
//and zero value for other fields
func NewUser(email string) *User {
	user := new(User)
	user.Email.String = email
	user.FullName.String = ""
	user.PasswordHash.String = ""
	user.Salt.String = ""
	return user
}
