package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

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

func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Handler collect them", *request.URL, request.URL.Scheme)
	fmt.Fprintf(writer, "Hi there, I love %s!", request.URL.Path[1:])
}

// Let's create a handler, viewHandler that will allow users to view a wiki page. It will handle URLs prefixed with "/view/".
func viewHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len("/view/"):]

	page, err := loadPage(title)

	// fmt.Println(string(title), page)
	// fmt.Println("We loaded page")
	if err != nil {
		// writer.Write([]byte("We cannot find the page you are looking for"))
		http.Redirect(writer, request, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(writer, "view.html", page)
	// fmt.Fprintf(writer, "<h1>%s</h1><div>%s</div>", page.Title, page.Body)
}

func editHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len("/edit/"):]
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

	renderTemplate(writer, "edit.html", p)
}

func saveHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len("/save/"):]

	body := request.FormValue("body")

	pa := &Page{Body: []byte(body), Title: title}

	pa.save()

	fmt.Println("We just finished saving")
	http.Redirect(writer, request, "/view/"+title, http.StatusFound)
}

func renderTemplate(writer http.ResponseWriter, templ string, page *Page) {
	te, _ := template.ParseFiles(templ)

	te.Execute(writer, page)
}

func setupServer() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)

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
