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
	IsReference  bool
	IsJson       bool
}

type GenerateRepositoryModel struct {
	ProjectName string
	Name        string
	NameLower   string
	Package     string
	Fields      []GenerateRepositoryModelField
}

type GenerateStructField struct {
	Name   string
	Type   string
	DbName string
}

type GenerateStructEnumValue struct {
	Name  string
	Value string
}

type GenerateStructEnum struct {
	Name   string
	Values []GenerateStructEnumValue
}

type GenerateStructModel struct {
	Name    string
	Package string
	Fields  []GenerateStructField
	Enums   []GenerateStructEnum
}

func GenerateStruct(params GenerateRepositoryParams) error {
	parsedStruct := GenerateParsedStruct(params.StructType)
	logger.Debug("parsedStruct", parsedStruct)
	fields := make([]GenerateStructField, len(parsedStruct.Fields)+3)
	fields[0] = GenerateStructField{
		Name:   "Id",
		Type:   "int",
		DbName: "id",
	}

	enums := make([]GenerateStructEnum, 0)

	for index, field := range parsedStruct.Fields {
		typeName := field.Type
		if strings.HasPrefix(typeName, fmt.Sprint(params.Package, ".")) {
			typeName = strings.TrimPrefix(typeName, fmt.Sprint(params.Package, "."))
		}
		fields[index+1] = GenerateStructField{
			Name:   field.Name,
			Type:   typeName,
			DbName: field.DbName,
		}
		logger.Debug("field", field.Name, "and type", field.Type)
		if field.IsEnum {
			stateName := createEnumName(parsedStruct, field)
			fields[index+1].Type = stateName
			values := make([]GenerateStructEnumValue, len(field.EnumValues))
			for i, value := range field.EnumValues {
				values[i] = GenerateStructEnumValue{
					Name:  fmt.Sprint(stateName, firstLetterToUpper(value)),
					Value: value,
				}
			}
			enums = append(enums, GenerateStructEnum{
				Name:   stateName,
				Values: values,
			})
		}
	}

	fields[len(fields)-2] = GenerateStructField{
		Name:   "CreatedAt",
		Type:   "time.Time",
		DbName: "created_at",
	}
	fields[len(fields)-1] = GenerateStructField{
		Name:   "UpdatedAt",
		Type:   "time.Time",
		DbName: "updated_at",
	}

	generateStructModel := GenerateStructModel{
		Package: params.Package,
		Name:    parsedStruct.Name,
		Fields:  fields,
		Enums:   enums,
	}
	return generateAndWriteFile(params, generateStructModel, parsedStruct, "struct", fmt.Sprint(strings.ToLower(parsedStruct.Name), "_gen.go"))
}

func createEnumName(parsedStruct ParsedStruct, field ParsedStructField) string {
	return fmt.Sprint(parsedStruct.Name, field.Name)
}

func GenerateRepository(params GenerateRepositoryParams) error {
	err := GenerateStruct(params)
	if err != nil {
		return err
	}
	parsedStruct := GenerateParsedStruct(params.StructType)
	model := GenerateRepositoryModel{
		ProjectName: params.ProjectName,
		Name:        parsedStruct.Name,
		NameLower:   firstLetterToLower(parsedStruct.Name),
		Package:     params.Package,
		Fields:      createFields(parsedStruct.Fields),
	}

	fileName := fmt.Sprint(strings.ToLower(parsedStruct.Name), "-repository_gen.go")
	return generateAndWriteFile(params, model, parsedStruct, "repository", fileName)
}

func generateAndWriteFile(params GenerateRepositoryParams, model any, parsedStruct ParsedStruct, templateName string, fileName string) error {
	tmpl := template.Must(template.ParseGlob("generator/template/*.tmpl"))
	buffer := bytes.NewBuffer([]byte{})
	err := tmpl.ExecuteTemplate(buffer, templateName, model)
	if err != nil {
		logger.Error("Failed to generate repository template", err.Error())
		return err
	}
	formatted, err := FormatResult(buffer)
	if err != nil {
		logger.Error("Failed to format repository template", err.Error())
		return err
	}
	err = WriteBytesToFile(formatted, fmt.Sprint(params.Directory, "/", fileName))
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
			IsReference:       hasPrismaReference(field),
			IsJson:            isJson(field),
		}
	}
	return modelFields
}

func isJson(field ParsedStructField) bool {
	return field.Type != "string" && field.Type != "int" && field.Type != "bool"
}
