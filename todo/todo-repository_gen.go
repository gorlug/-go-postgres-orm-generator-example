package todo

// GENERATED FILE
// DO NOT EDIT

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-generator-example/logger"
)

type TodoRepository struct {
	connPool *pgxpool.Pool
	dialect  goqu.DialectWrapper
}

func NewTodoRepository(connPool *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{
		connPool: connPool,
		dialect:  goqu.Dialect("postgres"),
	}
}

func (r *TodoRepository) Create(todo Todo) error {
	sql, args, err := r.dialect.Insert("Todo").
		Prepared(true).
		Rows(goqu.Record{

			"name":    todo.Name,
			"checked": todo.Checked,
			"state":   todo.State,
		}).
		ToSQL()
	if err != nil {
		logger.Error("error creating create Todo sql: %v", err)
		return err
	}

	_, err = r.connPool.Exec(context.Background(), sql, args...)
	if err != nil {
		logger.Error("error creating Todo: %v", err)
		return err
	}

	return nil
}
