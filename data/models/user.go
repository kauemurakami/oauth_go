package models

import (
	"errors"
	sec "oauth2/core/security"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Pass      string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (user *User) Prepare(step string) error {
	if err := user.validate(step); err != nil {
		return err
	}
	if err := user.formatter(step); err != nil {
		return err
	}
	return nil
}

func (user *User) validate(step string) error {
	if user.Name == "" {
		return errors.New("Insert a name")
	}
	if user.Email == "" {
		return errors.New("Insert an E-mail")
	}
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New("E-mail is invalid")
	}
	if step == "signup" && user.Pass == "" {
		return errors.New("Insert a password")
	}

	return nil
}

func (user *User) formatter(step string) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
	user.Pass = strings.TrimSpace(user.Pass)
	if step == "signup" {
		hashPass, err := sec.Hash(user.Pass)
		if err != nil {
			return err
		}
		user.Pass = string(hashPass)
	}
	return nil
}
