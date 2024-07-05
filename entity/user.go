package entity

import (
	"fashion-api/infra/config"
	"fashion-api/pkg/exception"

	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int       `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (u *User) ValidateToken(bearerToken string) exception.Exception {

	isBearer := strings.HasPrefix(bearerToken, "Bearer")

	if !isBearer {
		return exception.NewUnauthenticationError("invalid token")
	}

	tokenFields := strings.Fields(bearerToken)

	if len(tokenFields) != 2 {
		return exception.NewUnauthenticationError("invalid token")
	}

	tokenString := tokenFields[1]

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, exception.NewUnauthenticationError("invalid token")
		}

		return []byte(config.NewAppConfig().JWTSecretKey), nil
	})

	if err != nil {
		return exception.NewUnauthenticationError("invalid token")
	}

	mapClaims := jwt.MapClaims{}

	if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return exception.NewUnauthenticationError("invalid token")
	} else {
		mapClaims = claims
	}

	email, ok := mapClaims["email"].(string)

	if !ok {
		return exception.NewUnauthenticationError("invalid token")
	}

	u.Email = email

	role, ok := mapClaims["role"].(string)

	if !ok {
		return exception.NewUnauthenticationError("invalid token")
	}

	u.Role = role

	_, ok = mapClaims["expired_at"]

	if !ok {
		return exception.NewUnauthenticationError("invalid token")
	}

	return nil
}

func (u *User) GenereateTokenString() string {

	claims := jwt.MapClaims{
		"email":      u.Email,
		"role":       u.Role,
		"expired_at": time.Now().Add(8 * time.Hour).UnixMilli(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(config.NewAppConfig().JWTSecretKey))

	return tokenString
}

func (u *User) CompareHashPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) GenerateHashPassword() {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashPassword)
}
