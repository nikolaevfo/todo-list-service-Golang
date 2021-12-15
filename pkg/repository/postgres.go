// в данном файле реалзиуется логика подключения БД и хранятся имена таблиц

package repository

// sqlx - пакет для работы с бд
import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// описвания названия таблиц из БД для использования в методах модуля
const (
	usersTable      = "users"
	todoListsTable  = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

// параметры для БД
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	// запуск БД
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	// проверка подключения к БД
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
