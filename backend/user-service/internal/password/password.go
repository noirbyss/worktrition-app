package password

import "golang.org/x/crypto/bcrypt"

const Cost = bcrypt.DefaultCost

func Hash(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), Cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func Compare(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
