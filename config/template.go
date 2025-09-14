package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const (
	tmplPath = "templates/"
)

var templates map[string]*template.Template

func init() {
	templates = make(map[string]*template.Template)
	err := filepath.Walk(tmplPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := info.Name()
		tmpl, err := template.New(name).ParseFiles(path)
		if err != nil {
			return errors.Join(fmt.Errorf("parse template %q: %w", name, err), err)
		}
		templates[name] = tmpl
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func GetTemplate(name string) *template.Template {
	tmpl, ok := templates[name]
	if !ok {
		log.Fatalf("template %q not found", name)
	}
	return tmpl
}
