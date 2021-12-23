package main

import "os"

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

}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: data}, nil
}
