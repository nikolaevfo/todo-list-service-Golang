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

// Методы сервиса вызывают соответствующие методы из модуля repository,
// передаем данные на уровень ниже

// объявляем соль для хэширования пароля,время жизни токена
// и набор случайных байтов для подписи токена (ключ подписи)
const (
	salt       = "adfj;f894ojfeoaf93"
	tokenTTL   = 12 * time.Hour
	signingKey = "dfa3fef4breg43f43"
)

// объявим структуру для настройки генерации токена с UserId
type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

// описываем структуру сервиса авторизации, котороя принимает в контструктор
// репозиторий для работы с базой
type AuthService struct {
	repo repository.Authorization
}

// описываем конструктор инициализации сервиса авторизации
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

// публичные метод для генерация токена
func (s *AuthService) GenerateToken(username, password string) (string, error) {
	// получаем юзера из БД, использую метод из repository
	// так как пароль храниться в хэшированном виде, передаем его с помощью generatePasswordHash
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	// если user существует, генерируем токен с помощью библиотеки jwt:
	// превый аргумент - метод подписи,
	// второй аргумент - json с различными полями,
	// стандартные настройки: время жизни токена (ExpiresAt),
	// время, когда токен был сгенерирован (IssuedAt),
	// а также user.Id (токен будет содержать его внутри себя)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	// возвращаем подписанный с помощью ключа токен
	return token.SignedString([]byte(signingKey))
}

// метод для хэширования пароля, использует библиотеку sha1 и "соль"
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

// метод для парсинга токена и получения id, используется в хендлере middleware.go
func (s *AuthService) ParseToken(accessToken string) (int, error) {
	// используем ParseWithClaims пакета jwt
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// проверяем  метод подписи токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		// если ок, возвращаем ключ подписи
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	// приведем поле Claims объекта token к структуре tokenClaims
	claims, ok := token.Claims.(*tokenClaims)

	if !ok {
		return 0, errors.New(("token claims are not of type"))
	}

	// возвращаем UserId
	return claims.UserId, nil
}
