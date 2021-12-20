package repository

import (
	"fmt"
	"strings"
	todo "to-do-list"

	"github.com/jmoiron/sqlx"
)

// создаем структуру репозитория
type TodoItemPostgres struct {
	db *sqlx.DB
}

// создаем конструктор репозитория для работы с задачами
func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) CreateItem(listId int, item todo.TodoItem) (int, error) {
	// создаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// создаем запись в todoItemsTable
	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)
	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	if err := row.Scan(&itemId); err != nil {
		// в случае ошибки останавливаем транзакцию и откатываем изменения
		tx.Rollback()
		return 0, err
	}

	// создаем запись в listsItemsTable
	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemsTable)
	_, err = tx.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		// в случае ошибки останавливаем транзакцию и откатываем изменения
		tx.Rollback()
		return 0, err
	}

	// применяем изменения к БД и закрываем транзакцию
	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAllItems(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem

	// команда INNER JOIN позволяет выбрать только те элементы, которые есть в обеих таблицах
	// делаем выборку из todoItemsTable, при этом "джойним" listsItemsTable и usersListsTable
	query := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id INNER JOIN %s ul on ul.list_id = li.list_id WHERE li.list_id = $1 AND ul.user_id = $2", todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemPostgres) GetItemById(userId, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem

	query := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id INNER JOIN %s ul on ul.list_id = li.list_id WHERE ti.id = $1 AND ul.user_id = $2", todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Get(&item, query, itemId, userId); err != nil {
		return item, err
	}

	return item, nil
}

func (s *TodoItemPostgres) UpdateItem(userId, itemId int, input todo.UpdateItemInput) error {
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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	// переменная setValues используются для создания запроса такого вида:
	// title=$1
	// description=$1
	// decription=$1, done=$2
	// title=$1, decription=$2, done=$3
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s ti SET %s FROM %s li, %s ul WHERE ti.id = li.item_id AND li.list_id=ul.list_id AND ul.user_id = $%d AND ti.id=$%d", todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)
	args = append(args, userId, itemId)

	_, err := s.db.Exec(query, args...)

	return err
}

func (r *TodoItemPostgres) DeleteItem(userId, itemId int) error {
	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2`, todoItemsTable, listsItemsTable, usersListsTable)

	_, err := r.db.Exec(query, userId, itemId)

	return err
}
