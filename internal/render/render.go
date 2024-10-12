package render

import (
	"bytes"
	"fmt"
	"github.com/bertoxic/graphqlChat/internal/models"
	"github.com/bertoxic/graphqlChat/internal/utils"
	"github.com/bertoxic/graphqlChat/pkg/config"
	"html/template"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	app = a
}
func Template(w http.ResponseWriter, tmpl string, templateData *models.TemplateData) {
	tc, ok := app.TemplateCache[tmpl]
	if !ok {
		http.Error(w, fmt.Sprintf("Template %s not found", tmpl), http.StatusInternalServerError)
		return
	}

	buf := &bytes.Buffer{}
	err := tc.Execute(buf, templateData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %s", err), http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing response: %s", err), http.StatusInternalServerError)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {

	templateMap := make(map[string]*template.Template)
	pathToTemplate, err := utils.FindDirectory("templates")
	if err != nil {
		return nil, fmt.Errorf("error finding templates directory: %w", err)
	}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*page.gohtml", pathToTemplate))
	if err != nil {
		return nil, fmt.Errorf("error finding page templates: %w", err)
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil, fmt.Errorf("error parsing template %s: %w", name, err)
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*layout.gohtml", pathToTemplate))
		if err != nil {
			return nil, fmt.Errorf("error finding layout templates: %w", err)
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*layout.gohtml", pathToTemplate))
			if err != nil {
				return nil, fmt.Errorf("error parsing layout templates: %w", err)
			}
		}

		templateMap[name] = ts
	}

	return templateMap, nil
}
