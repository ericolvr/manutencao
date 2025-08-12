package domain

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	// ROLE
	// 1 Admin
	// 2 Finaceiro
	// 3 Suporte
	// 4 Técnicos
	Role   int64 `json:"role"`
	Status bool  `json:"status"`
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) error {
	hashedPassword = strings.TrimSpace(hashedPassword)
	password = strings.TrimSpace(password)

	if !strings.HasPrefix(hashedPassword, "$2a$") && !strings.HasPrefix(hashedPassword, "$2b$") {
		return fmt.Errorf("formato de hash inválido")
	}

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
