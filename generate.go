package main

import (
	"go-postgres-generator-example/generator"
	"go-postgres-generator-example/logger"
	"go-postgres-generator-example/todo"
)

//go:generate go run generate.go
//go:generate go run github.com/steebchen/prisma-client-go format

func main() {
	logger.Debug("Generating....")
	logger.Debug("Generating schema")
	err := generator.GenerateSchema([]any{todo.Todo{}})
	if err != nil {
		panic(err)
	}
	logger.Debug("Generating Todo repository")
	err = generator.GenerateRepository(generator.GenerateRepositoryParams{
		StructType:  todo.Todo{},
		Directory:   "todo",
		Package:     "todo",
		ProjectName: "go-postgres-generator-example",
	})
	if err != nil {
		panic(err)
	}
}
