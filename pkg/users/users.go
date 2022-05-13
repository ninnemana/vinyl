package users

//go:generate protoc --proto_path=. --go_out=. --go_opt=paths=source_relative ./users.proto --proto_path=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.2/ --proto_path=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis/ --proto_path=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/

import (
	"errors"
	"net/url"

	"cloud.google.com/go/firestore"
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

func GetFromQueryString(values url.Values) (GetParams, error) {
	var params GetParams

	if vals, ok := values["id"]; ok && len(vals) == 1 {
		params.Id = vals[0]
	}

	if vals, ok := values["email"]; ok && len(vals) == 1 {
		params.Email = vals[0]
	}

	return GetParams{
		Id:    params.Id,
		Email: params.Email,
	}, nil
}

func (x *User) Validate() error {
	if x.GetName() == "" {
		return ErrInvalidName
	}

	if x.GetEmail() == "" {
		return ErrInvalidEmail
	}

	if x.GetPassword() == "" {
		return ErrInvalidPassword
	}

	return nil
}

func (x *GetParams) Where(query firestore.Query) firestore.Query {
	if x.Id != "" {
		query = query.Where("id", "==", x.GetId())
	}

	if x.Email != "" {
		query = query.Where("email", "==", x.GetEmail())
	}

	for _, acct := range x.GetAuthenticatedAccounts() {
		query = query.Where("authenticatedAccounts.id", "==", acct.GetId())
	}

	return query
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
