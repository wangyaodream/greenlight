package data

import (
	"errors"
	"time"

	"github.com/wangyaodream/greenlight/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

// 在使用`json:"-"`标签时，表示在序列化和反序列化时忽略该字段
type User struct {
    ID int64 `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    Name string `json:"name"`
    Email string `json:"email"`
    Password password `json:"-"`
    Activated bool `json:"activated"`
    Version int `json:"-"`
}

type password struct {
    plaintext *string
    hash []byte
}

func (p *password) Set(plaintextPassword string) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
    if err != nil {
        return err
    }

    p.plaintext = &plaintextPassword
    p.hash = hash

    return nil
}

func (p *password) Matchs(plaintextPassword string) (bool, error) {
    err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
    if err != nil {
        switch {
        case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
            // 密码不匹配
            return false, nil
        default:
            return false, err
        }
    }

    return true, nil

}

func ValidateEmail(v *validator.Validator, email string) {
    v.Check(email != "", "email", "must be provided")
    v.Check(validator.Matchs(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
    v.Check(password != "", "password", "must be provided")
    v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
    v.Check(len(password) <= 72, "password", "must be at most 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
    v.Check(user.Name != "", "name", "must be provided")
    v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

    ValidateEmail(v, user.Email)

    if user.Password.plaintext != nil {
        ValidatePasswordPlaintext(v, *user.Password.plaintext)
    }

    if user.Password.hash == nil {
        panic("missing password hash for user")
    }
}