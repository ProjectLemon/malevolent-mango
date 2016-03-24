package main

import (
	"testing"
)

func TestValidateEmailWorking(t *testing.T) {
	user := User{Email: "test@exemple.com"}
	err := validateEmail(&user)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestValidateEmailNonWorking(t *testing.T) {
	user := User{Email: "test@"}
	err := validateEmail(&user)
	if err == nil {
		t.Fatalf("Should not be allowed to omit domain")
	}
}

func TestValidateEmailEmptyString(t *testing.T) {
	user := User{Email: ""}
	err := validateEmail(&user)
	if err == nil {
		t.Fatalf("Email strings should not be allowed to be empty")
	}
}

func TestValidateEmailAgainstSQLInject(t *testing.T) {
	user := User{Email: "test@exemple.com'"}
	err := validateEmail(&user)
	if err == nil {
		t.Fatalf("Emails should not be allowed to end on a '")
	}
}

func TestValidateEmailSQLInject(t *testing.T) {
	user := User{Email: "test@user.com'select * from Users;"}
	err := validateEmail(&user)
	if err == nil {
		t.Fatalf("We really should really not allow sql injections")
	}
}
