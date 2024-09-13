package generator

import (
	"bytes"
	"fmt"
	"go-postgres-generator-example/logger"
	"strings"
	"text/template"
)

type GenerateRepositoryParams struct {
	StructType  any
	Directory   string
	Package     string
	ProjectName string
}

type GenerateRepositoryModelField struct {
	ParsedStructField
	IsNotIdField bool
}

type GenerateRepositoryModel struct {
	ProjectName string
	Name        string
	NameLower   string
	Package     string
	Fields      []GenerateRepositoryModelField
}

func GenerateRepository(params GenerateRepositoryParams) error {
	parsedStruct := GenerateParsedStruct(params.StructType)
	model := GenerateRepositoryModel{
		ProjectName: params.ProjectName,
		Name:        parsedStruct.Name,
		NameLower:   strings.ToLower(parsedStruct.Name[:1]) + parsedStruct.Name[1:],
		Package:     params.Package,
		Fields:      createFields(parsedStruct.Fields),
	}

	tmpl := template.Must(template.ParseGlob("generator/template/*.tmpl"))
	buffer := bytes.NewBuffer([]byte{})
	err := tmpl.ExecuteTemplate(buffer, "repository", model)
	if err != nil {
		logger.Error("Failed to generate repository template", err.Error())
		return err
	}
	formatted, err := FormatResult(buffer)
	if err != nil {
		logger.Error("Failed to format repository template", err.Error())
		return err
	}
	err = WriteBytesToFile(formatted, fmt.Sprint(params.Directory, "/", strings.ToLower(parsedStruct.Name), "-repository_gen.go"))
	if err != nil {
		logger.Error("failed to write repository file", err)
	}
	return nil
}

func createFields(fields []ParsedStructField) []GenerateRepositoryModelField {
	modelFields := make([]GenerateRepositoryModelField, len(fields))
	for index, field := range fields {
		modelFields[index] = GenerateRepositoryModelField{
			ParsedStructField: field,
			IsNotIdField:      field.Name != "Id",
		}
	}
	return modelFields
}
