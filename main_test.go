package main

import (
	"testing"
)

//Correctness tests
func TestValidateEmailWorking(t *testing.T) {
	email := "test@exemple.com"
	err := validateEmail(email)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestValidateEmailNonWorking(t *testing.T) {
	email := "test@"
	err := validateEmail(email)
	if err == nil {
		t.Fatalf("Should not be allowed to omit domain")
	}
}

func TestValidateEmailEmptyString(t *testing.T) {
	email := ""
	err := validateEmail(email)
	if err == nil {
		t.Fatalf("Email strings should not be allowed to be empty")
	}
}

func TestValidateEmailAgainstSQLInject(t *testing.T) {
	email := "test@exemple.com'"
	err := validateEmail(email)
	if err == nil {
		t.Fatalf("Emails should not be allowed to end on a '")
	}
}

func TestValidateEmailSQLInject(t *testing.T) {
	email := "test@user.com'select * from Users;"
	err := validateEmail(email)
	if err == nil {
		t.Fatalf("We really should not allow sql injections")
	}
}

func TestGenerateSalt(t *testing.T) {
	bytes := randBase64String(128)
	if len(bytes) < 100 {
		t.Fatalf("Should have read 100b, only got: ", len(bytes))
	}
}

func TestValidateToken(t *testing.T) {
	userId := randBase64String(64)
	token, _ := generateToken(userId)
	session := UserSession{SessionKey: token}
	user := User{UserId: userId, Session: &session}
	valid, _ := validateToken(&user)
	if !valid {
		t.Fatalf("Token should be valid")
	}
}

//Benchmark tests
func BenchmarkGenerateToken(b *testing.B) {
	userId := randBase64String(64)
	for n := 0; n < b.N; n++ {
		generateToken(userId)
	}
}

func BenchmarkRandString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		randBase64String(128)
	}
}

func BenchmarkValidateEmail(b *testing.B) {
	for n := 0; n < b.N; n++ {
		validateEmail("linus.lagerhjelm@gmail.com")
	}
}

func BenchmarkValidateToken(b *testing.B) {
	userId := randBase64String(64)
	token, _ := generateToken(userId)
	session := UserSession{SessionKey: token}
	user := User{UserId: userId, Session: &session}
	for n := 0; n < b.N; n++ {
		validateToken(&user)
	}
}
