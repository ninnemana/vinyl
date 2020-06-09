package users

import (
	"errors"

	firestore "cloud.google.com/go/firestore"
	"golang.org/x/crypto/bcrypt"
)

type ValidationError error

var (
	ErrNotFound                        = errors.New("user not found")
	ErrUserExists                      = errors.New("user exists")
	ErrNotAuthorized                   = errors.New("user is not authorized")
	ErrInvalidName     ValidationError = errors.New("user name was blank")
	ErrInvalidEmail    ValidationError = errors.New("user email was blank")
	ErrInvalidPassword ValidationError = errors.New("user password was not valid")
)

func (a AuthenticatedAccount) Where(coll *firestore.CollectionRef) {
	coll.Where("authenticatedAccounts.id", "=", a.GetId())
}

func (u User) Validate() error {
	if u.GetName() == "" {
		return ErrInvalidName
	}

	if u.GetEmail() == "" {
		return ErrInvalidEmail
	}

	if u.GetPassword() == "" {
		return ErrInvalidPassword
	}

	return nil
}

func HashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePasswords(hashedPwd, plainPwd []byte) bool {
	if err := bcrypt.CompareHashAndPassword(
		hashedPwd,
		plainPwd,
	); err != nil {
		return false
	}

	return true
}
