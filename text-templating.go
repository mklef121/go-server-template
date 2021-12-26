package main

import (
	"os"
	"text/template"
)

type Todo struct {
	Name        string
	Description string
}

func testTemplating() {
	td := Todo{"Test templates", "Let's test a template to see the magic."}

	templ, err := template.New("todos").Parse("You have a task named \"{{ .Name}}\" with description: \"{{ .Description}}\"")

	if err != nil {
		panic(err)
	}

	err = templ.Execute(os.Stdout, td)

	if err != nil {
		panic(err)
	}
}
