package main

import (
	"go-postgres-generator-example/generator"
	"go-postgres-generator-example/logger"
	"go-postgres-generator-example/todo"
	"go-postgres-generator-example/user"
)

//go:generate go run generate.go
//go:generate go run github.com/steebchen/prisma-client-go format

func main() {
	logger.Debug("Generating....")
	logger.Debug("Generating schema")
	todoEntity := todo.TodoEntity{}
	err := generator.GenerateSchema([]any{user.UserEntity{}, todoEntity})
	if err != nil {
		panic(err)
	}
	logger.Debug("Generating Todo repository")
	projectName := "go-postgres-generator-example"
	err = generator.GenerateRepository(generator.GenerateRepositoryParams{
		StructType:  todoEntity,
		Directory:   "todo",
		Package:     "todo",
		ProjectName: projectName,
	})
	if err != nil {
		panic(err)
	}

	logger.Debug("Generating User repository")
	err = generator.GenerateRepository(generator.GenerateRepositoryParams{
		StructType:  user.UserEntity{},
		Directory:   "user",
		Package:     "user",
		ProjectName: projectName,
	})
	if err != nil {
		panic(err)
	}
}
