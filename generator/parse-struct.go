package generator

import (
	"reflect"
	"strings"
)

type ParsedStructField struct {
	Name                string
	DbName              string
	Type                string
	ParentStructName    string
	OriginalStructField reflect.StructField
	EnumValues          []string
	IsEnum              bool
}

type ParsedStruct struct {
	Name   string
	Fields []ParsedStructField
}

func GenerateParsedStruct(structType any) ParsedStruct {
	t := reflect.TypeOf(structType)
	parsed := ParsedStruct{
		Name:   t.Name(),
		Fields: []ParsedStructField{},
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		isStructName := isStructName(field)
		if isStructName {
			parsed.Name = field.Name
			continue
		}
		dbValue := field.Tag.Get("db")
		typeName := field.Type.String()
		parsed.Fields = append(parsed.Fields, ParsedStructField{
			Name:                field.Name,
			DbName:              dbValue,
			Type:                typeName,
			ParentStructName:    t.Name(),
			OriginalStructField: field,
			EnumValues:          getEnumValues(field),
			IsEnum:              isEnum(field),
		})
	}
	return parsed
}

func isStructName(field reflect.StructField) bool {
	return field.Tag.Get("isStructName") == "true"
}

func getEnumValues(field reflect.StructField) []string {
	enumValues := field.Tag.Get("enum")
	valuesSplit := strings.Split(enumValues, ",")
	return valuesSplit
}

func isEnum(field reflect.StructField) bool {
	enumValues := field.Tag.Get("enum")
	return enumValues != ""
}
