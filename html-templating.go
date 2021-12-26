package main

import (
	"fmt"
	"html/template"
	"os"
)

type entry struct {
	Name string
	Done bool
}

type ToDo struct {
	User string
	List []entry
}

func renderTodo() {
	// Files are provided as a slice of strings.
	paths := []string{
		"todo.html",
	}

	var todos = ToDo{
		User: "Miracle",
		List: []entry{
			{
				Name: "Enugu",
				Done: true,
			},
			{
				Name: "Gombe",
				Done: false,
			},
			{
				Name: "Rivers",
				Done: true,
			},
		},
	}

	templ := template.Must(template.ParseFiles(paths...))
	err := templ.Execute(os.Stdout, todos)
	fmt.Println("")

	if err != nil {
		panic(err)
	}
}
