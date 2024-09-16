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

type SchemaReference struct {
	Name  string
	Model string
	Field string
}

type SchemaModel struct {
	Name                     string
	Fields                   []SchemaModelField
	References               []SchemaReference
	ModelsReferencingThisOne []string
}

type SchemaEnum struct {
	Name   string
	Values []string
}

type Schema struct {
	Models []*SchemaModel
	Enums  []SchemaEnum
}

func GenerateSchema(structTypes []any) error {
	schema := Schema{
		Models: []*SchemaModel{},
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
	logger.Debug("parsedStruct", parsedStruct)
	model := SchemaModel{
		Name:                     strings.ToLower(parsedStruct.Name),
		Fields:                   make([]SchemaModelField, len(parsedStruct.Fields)+3),
		ModelsReferencingThisOne: []string{},
	}
	addIdFieldToModel(&model)
	for index, field := range parsedStruct.Fields {
		if hasPrismaReference(field) {
			addPrismaReference(&model, field, schema)
		}
		annotations := getPrismaAnnotations(field)
		model.Fields[index+1] = SchemaModelField{
			Name:             field.DbName,
			Type:             getDbType(parsedStruct, field),
			SchemaAdditional: fmt.Sprint(annotations, " ", getDefaultValue(field, field.EnumValues, annotations)),
		}
		if field.IsEnum {
			schemaEnum := SchemaEnum{
				Name:   getEnumName(parsedStruct, field),
				Values: field.EnumValues,
			}
			schema.Enums = append(schema.Enums, schemaEnum)
		}
	}
	addCreatedAndUpdatedFields(&model)
	logger.Debug("model fields", model.Fields)
	schema.Models = append(schema.Models, &model)
	return model
}

func addPrismaReference(s *SchemaModel, field ParsedStructField, schema *Schema) {
	reference := getPrismaReference(field)
	s.References = append(s.References, SchemaReference{
		Name:  field.DbName,
		Model: reference,
		Field: "id",
	})
	for _, model := range schema.Models {
		if reference == model.Name {
			model.ModelsReferencingThisOne = append(model.ModelsReferencingThisOne, s.Name)
			break
		}
	}
}

func getPrismaAnnotations(field ParsedStructField) string {
	return field.OriginalStructField.Tag.Get("prisma")
}

func getPrismaReference(field ParsedStructField) string {
	return field.OriginalStructField.Tag.Get("prismaReference")
}

func addCreatedAndUpdatedFields(s *SchemaModel) {
	s.Fields[len(s.Fields)-2] = SchemaModelField{
		Name:             "created_at",
		Type:             "DateTime",
		SchemaAdditional: "@default(now())",
	}
	s.Fields[len(s.Fields)-1] = SchemaModelField{
		Name:             "updated_at",
		Type:             "DateTime",
		SchemaAdditional: "@default(now()) @updatedAt",
	}
}

func addIdFieldToModel(s *SchemaModel) {
	s.Fields[0] = SchemaModelField{
		Name:             "id",
		Type:             "Int",
		SchemaAdditional: "@id @default(autoincrement())",
	}
}

func getEnumName(parsedStruct ParsedStruct, field ParsedStructField) string {
	return fmt.Sprint(parsedStruct.Name, field.Name)
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
		return "@default(\"{}\")"
	}
}

func getDbType(parsedStruct ParsedStruct, field ParsedStructField) string {
	if field.IsEnum {
		return getEnumName(parsedStruct, field)
	}
	t := field.Type
	switch t {
	case "int":
		return "Int"
	case "string":
		return "String"
	case "bool":
		return "Boolean"
	default:
		return "Json"
	}
}
