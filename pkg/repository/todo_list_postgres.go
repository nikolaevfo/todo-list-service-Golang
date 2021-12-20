package repository

import (
	"fmt"
	"strings"
	todo "to-do-list"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// создаем структуру репозитория
type TodoListPostgres struct {
	db *sqlx.DB
}

// создаем конструктор репозитория для работы со списками
func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

// для создания списка необходимо сделать вставки в две таблицы:
// usersListsTable(связывает пользователей с их списками) и todoListsTable
// эти операции проводятся в транзакции (последовательность действий, которые должны выполниться полностью)
func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	// создаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// создаем запись в todoListsTable
	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	// записваем в переменную id
	if err := row.Scan(&id); err != nil {
		// в случае ошибки останавливаем транзакцию и откатываем изменения
		tx.Rollback()
		return 0, err
	}

	// осуществляем вставку в usersListsTable
	// связываем id пользователя и id нового списка
	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2) RETURNING id", usersListsTable)
	// метод Exec не возварщает никакой информации
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		// в случае ошибки останавливаем транзакцию и откатываем изменения
		tx.Rollback()
		return 0, err
	}

	// применяем изменения к БД и заканчиваем транзакцию
	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList
	// в $1 будет поподать userId в r.db.Select
	// команда INNER JOIN позволяет выбрать только те элементы, которые есть в обеих таблицах
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1", todoListsTable, usersListsTable)

	// записываем в lists результат запроса с помощью метода Select
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList

	// команда INNER JOIN позволяет выбрать только те элементы, которые есть в обеих таблицах
	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`, todoListsTable, usersListsTable)

	// записываем в list результат запроса с помощью метода Select
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}

func (r *TodoListPostgres) DeleteList(userId, listId int) error {
	// записи удаляем сразу из 2 таблиц
	query := fmt.Sprintf(`DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2`, todoListsTable, usersListsTable)

	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *TodoListPostgres) UpdateList(userId, listId int, input todo.UpdateListInput) error {
	// иницилазируем переменные,
	// после используем их для формирования запроса к БД
	// setValues будет использоваться для подстановки в строку запроса к БД
	setValues := make([]string, 0)
	// args будет содержать слайс переменных для передачи методу Exec
	args := make([]interface{}, 0)
	// argId будет хранить индекс переменной для подстановки в строку запроса
	argId := 1

	// осущствляем проверку наличия полей
	// если они существуют, добавляем их в созданные переменные
	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	// переменная setValues используются для создания запроса такого вида:
	// title=$1
	// description=$1
	// title=$1, decription=$2
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d", todoListsTable, setQuery, usersListsTable, argId, argId+1)
	args = append(args, listId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)

	return err
}
