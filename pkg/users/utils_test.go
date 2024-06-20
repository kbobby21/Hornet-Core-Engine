package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {

	password := "mySecretPassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Error("Hashed password is empty")
	}
}

func TestVerifyPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("mySecretPassword"), bcrypt.DefaultCost)

	//verifying the password
	isValid := VerifyPassword(string(hashedPassword), "mySecretPassword")
	assert.True(t, isValid, "Expected password to be valid")

	isValid = VerifyPassword(string(hashedPassword), "incorrectPassword")
	assert.False(t, isValid, "Expected password to be invalid")
}

func TestCreateToken(t *testing.T) {
	email := "gauravk@hornet.technology"
	tokenString, err := CreateToken(email)

	if err != nil {
		t.Fatalf("CreateToken returned an error: %v", err)
	}

	if len(tokenString) == 0 {
		t.Error("Token string is empty")
	}
}
