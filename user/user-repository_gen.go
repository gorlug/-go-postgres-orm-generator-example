package user

// GENERATED FILE
// DO NOT EDIT

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-generator-example/logger"
	"time"
)

type UserRepository struct {
	connPool *pgxpool.Pool
	dialect  goqu.DialectWrapper
}

func NewUserRepository(connPool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		connPool: connPool,
		dialect:  goqu.Dialect("postgres"),
	}
}

func (r *UserRepository) Create(user User) (int, error) {
	sql, args, err := r.dialect.Insert("user").
		Prepared(true).
		Rows(goqu.Record{

			"updated_at": time.Now(),
			"email":      user.Email,
			"state":      jsonToString(user.State),
		}).
		Returning("id").
		ToSQL()
	if err != nil {
		logger.Error("error creating create User sql: %v", err)
		return -1, err
	}

	rows, err := r.connPool.Query(context.Background(), sql, args...)
	if err != nil {
		logger.Error("error creating User: %v", err)
		return -1, err
	}
	defer rows.Close()
	var id int
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			logger.Error("error scanning User: %v", err)
			return -1, err
		}
	} else {
		logger.Error("User already exists")
		return -1, UserAlreadyExistsError{User: user}
	}

	return id, nil

}

type UserAlreadyExistsError struct {
	User User
}

func (e UserAlreadyExistsError) Error() string {
	return fmt.Sprint("User ", e.User, " already exists")
}

func (r *UserRepository) GetById(id int) (User, error) {
	logger.Debug("Getting User by id ", id)
	sql, args, _ := r.dialect.From("user").
		Prepared(true).
		Select(
			"id",
			"created_at",
			"updated_at",
			"email",
			"state",
		).
		Where(goqu.Ex{"id": id}).
		ToSQL()

	rows, err := r.connPool.Query(context.Background(), sql, args...)
	if err != nil {
		logger.Error("Failed to get User: ", err)
	}
	defer rows.Close()
	item := User{}
	for rows.Next() {
		err = rows.Scan(
			&item.Id,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.Email,
			&item.State,
		)
		if err != nil {
			logger.Error("Failed to scan User: ", err)
			return item, err
		}
	}
	return item, nil
}

func (r *UserRepository) Update(user User) error {
	sql, args, err := r.dialect.Update("user").
		Prepared(true).
		Set(goqu.Record{

			"updated_at": time.Now(),
			"email":      user.Email,
			"state":      jsonToString(user.State),
		}).
		Where(goqu.Ex{"id": user.Id}).
		ToSQL()
	if err != nil {
		logger.Error("error creating update User sql: %v", err)
		return err
	}

	_, err = r.connPool.Exec(context.Background(), sql, args...)
	if err != nil {
		logger.Error("error updating User: %v", err)
		return err
	}

	return nil
}

func (r *UserRepository) Delete(id int) error {
	sql, args, err := r.dialect.Delete("user").
		Prepared(true).
		Where(goqu.Ex{"id": id}).
		ToSQL()
	if err != nil {
		logger.Error("error creating delete User sql: %v", err)
		return err
	}

	_, err = r.connPool.Exec(context.Background(), sql, args...)
	if err != nil {
		logger.Error("error deleting User: %v", err)
		return err
	}

	return nil
}

func jsonToString(jsonData any) string {
	bytes, err := json.Marshal(jsonData)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}
