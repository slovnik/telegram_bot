package main

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/slovnik/slovnik"
)

type Template struct {
	tmpl *template.Template
}

func CreateTemplate() *Template {
	templateFiles := []string{
		"./tmpl/full-word.gotmpl",
		"./tmpl/short-word.gotmpl",
		"./tmpl/translation.gotmpl",
		"./tmpl/phrases.gotmpl",
	}

	funcs := template.FuncMap{
		"join": strings.Join,
	}

	tmpl, err := template.New("").Funcs(funcs).ParseFiles(templateFiles...)

	if err != nil {
		panic(err)
	}

	return &Template{tmpl}
}

func (t *Template) Execute(words []slovnik.Word) string {
	var buf bytes.Buffer
	t.tmpl.ExecuteTemplate(&buf, "translation", words)
	return buf.String()
}

func (t *Template) Phrases(words []slovnik.Word) string {
	var buf bytes.Buffer
	t.tmpl.ExecuteTemplate(&buf, "phrases", words)
	return buf.String()
}
