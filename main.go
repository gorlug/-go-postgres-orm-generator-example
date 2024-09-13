package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go-postgres-generator-example/logger"
	"go-postgres-generator-example/todo"
	"log"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.LogError("Error loading .env file: %v", err)
		panic(err)
	}

	connpool, err := CreatePostgesConnpool(os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.LogError("Failed to create connection pool: %v", err)
		panic(err)
	}

	todoRepository := todo.NewTodoRepository(connpool)
	todoToCreate := todo.Todo{
		Name:    "Test",
		Checked: false,
		State:   todo.TodoStateCreated,
	}
	err = todoRepository.Create(todoToCreate)
	if err != nil {
		panic(err)
	}
}

func CreatePostgesConnpool(dbUrl string) (*pgxpool.Pool, error) {
	connPool, err := pgxpool.NewWithConfig(context.Background(), config(dbUrl))
	if err != nil {
		logger.Error("Error while creating connection to the database!!", err)
		return nil, err
	}

	connection, err := connPool.Acquire(context.Background())
	if err != nil {
		logger.LogError("Error while acquiring connection from the database pool!! %v", err)
		return nil, err
	}
	defer connection.Release()

	err = connection.Ping(context.Background())
	if err != nil {
		logger.LogError("Could not ping database")
		return nil, err
	}

	logger.LogDebug("Connected to the database!!")
	return connPool, nil
}

func config(dbUrl string) *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	return dbConfig
}
