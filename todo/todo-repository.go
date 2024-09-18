package todo

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"go-postgres-generator-example/logger"
)

func (r *TodoRepository) GetCheckedTodos(userId int) ([]Todo, error) {
	sql, args, err := r.dialect.From("todo").
		Prepared(true).
		Select("id", "created_at", "updated_at", "name", "checked", "state", "user_id").
		Where(goqu.Ex{"user_id": userId, "checked": true}).
		ToSQL()
	if err != nil {
		logger.Error("error creating get checked Todos sql: %v", err)
		return nil, err
	}

	rows, err := r.connPool.Query(context.Background(), sql, args...)
	if err != nil {
		logger.Error("error getting checked Todos: %v", err)
		return nil, err
	}
	defer rows.Close()

	var items []Todo
	for rows.Next() {
		item := Todo{}
		err = rows.Scan(
			&item.Id,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.Name,
			&item.Checked,
			&item.State,
			&item.UserId,
		)
		if err != nil {
			logger.Error("Failed to scan Todo: ", err)
			return items, err
		}
		items = append(items, item)
	}
	return items, nil
}
