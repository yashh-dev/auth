package security

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/matthewhartstonge/argon2"
)

func EncryptPassword(password string) (string, error) {
	argon := argon2.MemoryConstrainedDefaults()
	encoded, err := argon.HashEncoded([]byte(password))
	return string(encoded), err
}

func VerifyPassword(hash string, password string) (bool, error) {
	ok, err := argon2.VerifyEncoded([]byte(password), []byte(hash))
	return ok, err
}

func GenerateJWT(uid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "auth.miauw.social",
		"sub": uid,
		"iat": time.Now().Unix(),
		"nbf": time.Now().Add(time.Millisecond).Unix(),
		"exp": time.Now().Add(10 * time.Minute).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func VerifyJWT(token string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(
		token,
		func(to *jwt.Token) (interface{}, error) {
			if _, ok := to.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("abc")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)
	if err != nil {
		fmt.Println(err)
	}
	if !t.Valid {
		return jwt.MapClaims{}, err
	}

	return t.Claims.(jwt.MapClaims), nil
}
