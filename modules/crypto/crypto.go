package crypto

import "golang.org/x/crypto/bcrypt"

func HashPassword(
	pw string,
) (
	string,
	error,
) {

	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func Compare(
	pw string,
	hash string,
) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))

	return err == nil
}
