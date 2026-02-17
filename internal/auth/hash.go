package auth

import "github.com/alexedwards/argon2id"

func HashPassword(password string) (string, error) {
	hashed, err := argon2id.CreateHash(password, &argon2id.Params{
		Memory:      64 * 1024, // 64MB
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16, // Important: Specify salt length
		KeyLength:   32, // Important: Specify key length
	})
	if err != nil {
		return "", err
	}
	return hashed, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	ok, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return ok, nil
}
