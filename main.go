package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go-postgres-generator-example/logger"
	"go-postgres-generator-example/todo"
	"go-postgres-generator-example/user"
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

	userRepository := user.NewUserRepository(connpool)
	todoRepository := todo.NewTodoRepository(connpool)

	userId, userFromDb := createUser(userRepository)
	todoId, todoFromDb := createTodo(userId, todoRepository)

	updateUser(userFromDb, userRepository, userId)
	updateTodo(todoFromDb, todoRepository, todoId)
	createSecondTodo(userId, todoRepository)

	getAllCheckedTodos(todoRepository, userId)

	deleteTodo(todoRepository, todoId)
	deleteUser(userRepository, userId)
}

func getAllCheckedTodos(repository *todo.TodoRepository, id int) {
	checkedTodos, err := repository.GetCheckedTodos(id)
	if err != nil {
		panic(err)
	}
	logger.Debug("checked todos", checkedTodos)
}

func createSecondTodo(id int, repository *todo.TodoRepository) {
	todoToCreate := todo.Todo{
		Name:    "Test2",
		Checked: false,
		State:   todo.TodoStateCreated,
		UserId:  id,
	}
	todoId, err := repository.Create(todoToCreate)
	if err != nil {
		panic(err)
	}

	logger.Debug("create todo with todoId", todoId)
	todoFromDb, err := repository.GetById(todoId)
	if err != nil {
		panic(err)
	}
	logger.Debug("todo from db", todoFromDb)
}

func deleteUser(userRepository *user.UserRepository, userId int) {
	err := userRepository.Delete(userId)
	if err != nil {
		panic(err)
	}
}

func deleteTodo(todoRepository *todo.TodoRepository, todoId int) {
	err := todoRepository.Delete(todoId)
	if err != nil {
		panic(err)
	}
}

func updateTodo(todoFromDb todo.Todo, todoRepository *todo.TodoRepository, todoId int) {
	todoFromDb.Checked = true
	err := todoRepository.Update(todoFromDb)
	if err != nil {
		panic(err)
	}
	updatedTodoFromDb, err := todoRepository.GetById(todoId)
	if err != nil {
		panic(err)
	}
	logger.Debug("updated todo from db", updatedTodoFromDb)
}

func updateUser(userFromDb user.User, userRepository *user.UserRepository, userId int) {
	userFromDb.State.SomeValue = "updated value"
	err := userRepository.Update(userFromDb)
	if err != nil {
		panic(err)
	}
	updatedUserFromDb, err := userRepository.GetById(userId)
	if err != nil {
		panic(err)
	}
	logger.Debug("updated user from db", updatedUserFromDb)
}

func createTodo(userId int, todoRepository *todo.TodoRepository) (int, todo.Todo) {
	todoToCreate := todo.Todo{
		Name:    "Test",
		Checked: false,
		State:   todo.TodoStateCreated,
		UserId:  userId,
	}
	todoId, err := todoRepository.Create(todoToCreate)
	if err != nil {
		panic(err)
	}

	logger.Debug("create todo with todoId", todoId)
	todoFromDb, err := todoRepository.GetById(todoId)
	if err != nil {
		panic(err)
	}
	logger.Debug("todo from db", todoFromDb)
	return todoId, todoFromDb
}

func createUser(userRepository *user.UserRepository) (int, user.User) {
	userInstance := user.User{
		Email: "doesnot@shouldnotexist.com",
		State: user.UserState{
			SomeValue: "",
		},
	}
	userId, err := userRepository.Create(userInstance)
	if err != nil {
		logger.LogError("Failed to create user: %v", err)
		panic(err)
	}
	logger.Debug("create user with userId", userId)
	userFromDb, err := userRepository.GetById(userId)
	if err != nil {
		panic(err)
	}
	logger.Debug("user from db", userFromDb)
	return userId, userFromDb
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
