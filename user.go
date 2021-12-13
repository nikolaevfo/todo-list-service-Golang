package todo

// Описываем структуру пользователя, соответствующую базе данных.
// Применяется при чтении запросов клиента.
// Добавлены json-теги для корректного чтения и
// вывода данных при запросах, а также теги для базы данных.
type User struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
