package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/beefsack/gar"
	"github.com/go-martini/martini"
)

func compileTemplates() *template.Template {
	t := template.New("base")
	// Compile all templates
	files, err := gar.Files()
	if err != nil {
		log.Fatalf("Unable to get files, %v", err)
	}
	for _, f := range files {
		if !strings.HasPrefix(f, "template/") {
			continue
		}
		file, _, err := gar.Open(f)
		if err != nil {
			log.Fatalf("Unable to open template file, %v", err)
		}
		b, err := ioutil.ReadAll(file.Content)
		if err != nil {
			log.Fatalf("Unable to read template file, %v", err)
		}
		if err := file.Content.Close(); err != nil {
			log.Fatalf("Unable to close template file, %v", err)
		}
		t = template.Must(t.New(f).Parse(string(b)))
	}
	return t
}

func main() {
	m := martini.Classic()
	m.Map(compileTemplates())
	m.Get("/", func(w http.ResponseWriter, tmpl *template.Template) {
		tmpl.ExecuteTemplate(w, "template/home.html", nil)
	})
	m.Get("/public/**", func(
		w http.ResponseWriter,
		r *http.Request,
		params martini.Params,
	) {
		path := fmt.Sprintf("public/%s", params["_1"])
		file, ok, err := gar.Open(path)
		if err != nil {
			log.Fatalf("Error opening %s, %v", path, err)
		}
		defer file.Content.Close()
		if !ok {
			http.NotFound(w, r)
			return
		}
		http.ServeContent(w, r, file.FileInfo.Name(), file.FileInfo.ModTime(), file.Content)
	})
	m.Run()
}
