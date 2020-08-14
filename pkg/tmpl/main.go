package tmpl

import (
	"strings"
	"text/template"
)

// Package global for storing template variables and their user-assigned values
var tmplVars map[string]string

// TODO add more functions
var funcMap = template.FuncMap{
	"Title": strings.Title,
}

func init() {
	tmplVars = make(map[string]string)
	makeRunVars()
}
