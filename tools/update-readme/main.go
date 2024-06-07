package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

type ReadmeData struct {
	ExampleYaml string
}

func main() {
	tpl := template.Must(template.ParseFiles("./README.tpl.md"))

	example, err := os.ReadFile("./examples/example.yaml")
	if err != nil {
		panic(fmt.Errorf("load example: %v", err))
	}

	rd := ReadmeData{
		ExampleYaml: strings.TrimSpace(string(example)),
	}

	readme, err := os.Create("./README.md")
	if err != nil {
		panic(fmt.Errorf("create readme: %v", err))
	}
	defer readme.Close()

	if err := tpl.Execute(readme, rd); err != nil {
		panic(fmt.Errorf("write readme: %v", err))
	}
}
