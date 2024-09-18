package main

import (
	"bytes"
	"text/template"
)

type Example struct {
	Name string
}

func main() {
	tmpl := template.Must(template.ParseGlob("generator/template/*.tmpl"))
	buffer := bytes.NewBuffer([]byte{})
	example := Example{Name: "World"}
	err := tmpl.ExecuteTemplate(buffer, "example", example)
	if err != nil {
		panic(err)
	}
	println(buffer.String())
}
