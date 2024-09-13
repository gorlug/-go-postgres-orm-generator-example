package generator

import "reflect"

type ParsedStructField struct {
	Name                string
	DbName              string
	Type                string
	ParentStructName    string
	OriginalStructField reflect.StructField
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
		dbValue := field.Tag.Get("db")
		parsed.Fields = append(parsed.Fields, ParsedStructField{
			Name:                field.Name,
			DbName:              dbValue,
			Type:                field.Type.Name(),
			ParentStructName:    t.Name(),
			OriginalStructField: field,
		})
	}
	return parsed
}
