package render

import (
	"bytes"
	"html/template"
	"log"
	"myapp/pkg/config"
	"myapp/pkg/models"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}
var app *config.AppConfig

// NewTemplate sets the config for the template packege
func NewTemplate(a *config.AppConfig) {
	app = a
}
func AddDefaultData(templateData *models.TemplateData) *models.TemplateData {

	return templateData
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, templateData *models.TemplateData) {
	var templateCache map[string]*template.Template
	if app.UseCache {
		//get the template cache from the app config
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}
	//get requested template from cache
	template, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}
	buf := new(bytes.Buffer)
	templateData = AddDefaultData(templateData)
	_ = template.Execute(buf, templateData)
	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	Cache := map[string]*template.Template{}
	//get all the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return Cache, nil
	}
	//range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return Cache, nil
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return Cache, nil
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return Cache, nil
			}
		}
		Cache[name] = ts
	}
	return Cache, nil
}
