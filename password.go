package model

import "golang.org/x/crypto/bcrypt"

// Password описывает тип для пароля, хранящегося в виде хеш с использованием
// алгоритма bcrypt.
type Password []byte

// NewPassword возвращает пароль в виде хеш.
func NewPassword(password string) Password {
	passwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return Password(passwd)
}

// Compare сравнивает сохраненный в виде хеш пароль с указанным в параметре и
// возвращает true, если указанный пароль с очень большой степенью вероятности и
// является оригинальным паролем.
func (p Password) Compare(password string) bool {
	return bcrypt.CompareHashAndPassword(p, []byte(password)) == nil
}
