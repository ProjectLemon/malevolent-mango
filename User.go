package main

import "time"

//User provides a datatype for gathering the user information
//for the current user.
type User struct {
	Email    string
	FullName string
	Password string
	Salt     string
	Session  *UserSession
}

//UserSession holds the current session information for a specific user
type UserSession struct {
	SessionKey string
	LoginTime  time.Time
	LastSeen   time.Time
}
