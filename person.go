package main

import (
	"fmt"
	html_template "html/template"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"
)

type Person struct {
	Name   string
	Emails []string
}

const tmpl = `The name is {{.Name}}.
{{$name := .Name}}
{{range .Emails}}
His email id is {{.}} and his Name is {{$name}}
{{end}}
`

func parsePerson() {
	person := Person{
		Name:   "Satish",
		Emails: []string{"satish@rubylearning.org", "satishtalim@gmail.com"},
	}

	templ := template.New("person-parse")

	templ = template.Must(templ.Parse(tmpl))

	err := templ.Execute(os.Stdout, person)

	if err != nil {
		fmt.Println("An error occured", err)
		return
	}

	fmt.Println()

	fsH := http.FileServer(http.Dir("./public"))

	http.Handle("/public/", http.StripPrefix("/public", fsH))

	http.HandleFunc("/", ServeTemplate)

	fmt.Println("Listening...")
	err = http.ListenAndServe(GetPort(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return
	}
}

// Get the Port from the environment so we can run on Heroku (more of this later)
func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func ServeTemplate(w http.ResponseWriter, r *http.Request) {
	lp := path.Join("templates", "layout.html")
	fp := path.Join("templates", r.URL.Path)

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)

	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	templ, err := html_template.ParseFiles(lp, fp)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "500 Internal Server Error", 500)
		return
	}

	// hi := 'j'
	templ.ExecuteTemplate(w, "layout", nil)

}
