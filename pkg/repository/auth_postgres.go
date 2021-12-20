package repository

import (
	"fmt"
	todo "to-do-list"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

// описываем структуру репозитория для работы с авторизацией
type AuthPostgres struct {
	db *sqlx.DB
}

// создаем конструктор репозитория авторизации
func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// описываем метод, работающий с базой данных
func (r *AuthPostgres) CreateUser(user todo.User) (int, error) {
	var id int

	// осуществим запрос к базе данных, используя функцию для форматирования строк из fmt
	// используем метод INSERT для добавления в usersTable, возвращаем id
	// в плейсхолдеры $1, $2, $3 подставятся значения аргументов метода QueryRow, начиная со 2-го
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)

	// с помощью метода Scan записываем значение id в переменную
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// метод для получения user при проверке токена
func (r *AuthPostgres) GetUser(username, password string) (todo.User, error) {
	// объявим структуру, в которую будем записывать результат
	var user todo.User

	// описываем запрос к БД
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)

	// делаем запрос к БД, записываем в user
	err := r.db.Get(&user, query, username, password)

	return user, err
}
