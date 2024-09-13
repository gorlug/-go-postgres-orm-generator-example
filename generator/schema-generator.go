package generator

import (
	"bytes"
	"fmt"
	"go-postgres-generator-example/logger"
	"strings"
	"text/template"
)

type SchemaModelField struct {
	Name             string
	Type             string
	SchemaAdditional string
}

type SchemaModel struct {
	Name   string
	Fields []SchemaModelField
}

type SchemaEnum struct {
	Name   string
	Values []string
}

type Schema struct {
	Models []SchemaModel
	Enums  []SchemaEnum
}

func GenerateSchema(structTypes []any) error {
	schema := Schema{
		Models: []SchemaModel{},
		Enums:  []SchemaEnum{},
	}
	for _, structType := range structTypes {
		parsed := GenerateParsedStruct(structType)
		generateSchemaFromStruct(parsed, &schema)
	}
	tmpl := template.Must(template.ParseGlob("generator/template/*.tmpl"))
	buffer := bytes.NewBuffer([]byte{})
	err := tmpl.ExecuteTemplate(buffer, "schema", schema)
	if err != nil {
		logger.Error("Failed to generate schema template", err.Error())
		return err
	}
	err = WriteBufferToFile(buffer, "schema.prisma")
	if err != nil {
		logger.Error("failed to write schema file", err)
	}
	return nil
}

func generateSchemaFromStruct(parsedStruct ParsedStruct, schema *Schema) SchemaModel {
	model := SchemaModel{
		Name:   parsedStruct.Name,
		Fields: make([]SchemaModelField, len(parsedStruct.Fields)),
	}
	for index, field := range parsedStruct.Fields {
		enumValues := field.OriginalStructField.Tag.Get("enum")
		valuesSplit := strings.Split(enumValues, ",")
		schemaAdditional := field.OriginalStructField.Tag.Get("prisma")
		model.Fields[index] = SchemaModelField{
			Name:             field.DbName,
			Type:             getDbType(field.Type),
			SchemaAdditional: fmt.Sprint(schemaAdditional, " ", getDefaultValue(field, valuesSplit, schemaAdditional)),
		}
		if enumValues != "" {
			schemaEnum := SchemaEnum{
				Name:   field.Type,
				Values: valuesSplit,
			}
			schema.Enums = append(schema.Enums, schemaEnum)
		}
	}
	schema.Models = append(schema.Models, model)
	return model
}

func getDefaultValue(field ParsedStructField, values []string, additional string) string {
	if strings.Contains(additional, "default") {
		return ""
	}
	if len(values) > 1 {
		return fmt.Sprint("@default(", values[0], ")")
	}
	if field.DbName == "id" {
		return "@id @default(autoincrement())"
	}
	switch field.Type {
	case "int":
		return "@default(0)"
	case "string":
		return "@default(\"\")"
	case "bool":
		return "@default(false)"
	default:
		return ""
	}
}

func getDbType(t string) string {
	switch t {
	case "int":
		return "Int"
	case "string":
		return "String"
	case "bool":
		return "Boolean"
	default:
		return t
	}
}
