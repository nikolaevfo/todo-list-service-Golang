package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"
	todo "to-do-list"
	"to-do-list/pkg/repository"

	"github.com/dgrijalva/jwt-go"
)

// объявляем соль для хэширования пароля,
const (
	salt       = "adfj;f894ojfeoaf93"
	signingKey = "dfa3fef4breg43f43"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

// описываем структуру сервиса авторизации, котороя принимает в контструктор
// репозиторий для работы с базой
type AuthService struct {
	repo repository.Authorization
}

// описываем метод инициализации сервиса авторизации
func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// имплементируем метод создания пользователя
// здесь данные будут передаваться на слой ниже, в repository
func (s *AuthService) CreateUser(user todo.User) (int, error) {
	// хэшируем пароль
	user.Password = generatePasswordHash(user.Password)
	// вызываем метод из модуля repository
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	// получаем юзера из бд
	user, err := s.repo.GetUser(username, generatePasswordHash(password))

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

// метод для хэширования пароля, использует библиотеку sha1 и "соль"
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

// распарсим токен и вытащим id
func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)

	if !ok {
		return 0, errors.New(("token claims are not of type"))
	}

	return claims.UserId, nil
}
