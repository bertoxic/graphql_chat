package render

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"path/filepath"
)

var pathToTemplate = "templates"

func Template(w http.ResponseWriter, tmpl string) {
	templateMap, err := CreateTemplateCache()
	if err != nil {

	}
	tc := templateMap[tmpl]
	buf := &bytes.Buffer{}
	err = tc.Execute(buf, struct{}{})
	if err != nil {
	}
	_, err = buf.WriteTo(w)
	if err != nil {

	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	templateMap := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*page.gohtml", pathToTemplate))
	if err != nil {

	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {

		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*layout.gohtml", pathToTemplate))
		if err != nil {

		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*layout.gohtml"))
			if err != nil {
				return nil, err
			}
		}
		templateMap[name] = ts
	}
	return templateMap, nil
}
