package pwdutils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytePwd, err := bcrypt.GenerateFromPassword([]byte(password), 8)

	return string(bytePwd), err
}

func ComparePassword(password, hashedPwd string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password))
	return err
}
