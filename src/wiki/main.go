package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

// data structure of wiki
type Page struct {
	Title string // title
	Body  []byte // inner title
}

// set path address and get string length
const lenPath = len("/view/")

// make array of template file
 var templates = make(map[string]*template.Template)

// initialize function which is executed before main function
func init(){
	for _, tmpl := range []string{"edit","view"}{
		// if error, Must function generates panic so, there's no err description.
		t := template.Must(template.ParseFiles(tmpl + ".html"))
		templates[tmpl] = t

	}
}



// response handler
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[lenPath:]
	p, err := loadPage(title)

	// err handle (when you accessed the notexist page)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)

	// fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[lenPath:]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit", p)

	// fmt.Fprintf(w, "<h1>Editing %s</h1>"+
	// 	"<form action=\"/save/%s\" method=\"POST\">"+
	// 	"<textarea name=\"body\">%s</textarea><br>"+
	// 	"<input type=\"submit\" value=\"Save\">"+
	// 	"</form>", p.Title, p.Title, p.Body)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[lenPath:]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {

	// Embedd　Title, Body　to edit.html
	err := templates[tmpl].Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// save method of txt file
func (p *Page) save() error {
	// make text file by Title name and save the file.
	filename := p.Title + ".txt"
	// 0600 is permission settings, 0600 is permission which your own is permitted.
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// load file name from texttitle and return new Page pointer
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	// if get err, make body value nil and return as error
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler) // make edit page
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
