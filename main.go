package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Body  []byte
}

func (page *Page) save() error {

	filename := page.Title + ".txt"

	return os.WriteFile(filename, page.Body, 0600)

}
func main() {
	// os.
	p1 := &Page{Title: "TestPage", Body: []byte("I am come to the place")}
	p1.save()

	p2, _ := loadPage("TestPage")

	fmt.Println(p2.Body)

	fmt.Println("About Starting server")

	setupServer()

}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression.
}

func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Handler collect them", *request.URL, request.URL.Scheme)
	fmt.Fprintf(writer, "Hi there, I love %s!", request.URL.Path[1:])
}

// Let's create a handler, viewHandler that will allow users to view a wiki page. It will handle URLs prefixed with "/view/".
func viewHandler(writer http.ResponseWriter, request *http.Request, title string) {
	page, err := loadPage(title)

	// fmt.Println(string(title), page)
	// fmt.Println("We loaded page")
	if err != nil {
		// writer.Write([]byte("We cannot find the page you are looking for"))
		http.Redirect(writer, request, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(writer, "view", page)
	// fmt.Fprintf(writer, "<h1>%s</h1><div>%s</div>", page.Title, page.Body)
}

func editHandler(writer http.ResponseWriter, request *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	// fmt.Fprintf(writer, "<h1>Editing %s</h1>"+
	// 	"<form action=\"/save/%s\" method=\"POST\">"+
	// 	"<textarea name=\"body\">%s</textarea><br>"+
	// 	"<input type=\"submit\" value=\"Save\">"+
	// 	"</form>",
	// 	p.Title, p.Title, p.Body)

	renderTemplate(writer, "edit", p)
}

func saveHandler(writer http.ResponseWriter, request *http.Request, title string) {

	body := request.FormValue("body")

	pa := &Page{Body: []byte(body), Title: title}

	pa.save()

	fmt.Println("We just finished saving")
	http.Redirect(writer, request, "/view/"+title, http.StatusFound)
}

func makeHandler(callback func(w http.ResponseWriter, r *http.Request, title string)) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		m := validPath.FindStringSubmatch(request.URL.Path)
		if m == nil {
			http.NotFound(rw, request)
			return
		}

		callback(rw, request, m[2])
	}
}
func renderTemplate(writer http.ResponseWriter, templ string, page *Page) {
	// te, err := template.ParseFiles(templ)

	// if err != nil {
	// 	http.Error(writer, err.Error(), http.StatusInternalServerError)

	// 	return
	// }

	if err := templates.ExecuteTemplate(writer, templ+".html", page); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func setupServer() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8083", nil))
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: data}, nil
}
