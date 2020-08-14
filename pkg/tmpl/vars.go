package tmpl

import "time"

// SetVar -
func SetVar(key, value string) {
	tmplVars[key] = value
}

// TODO add more template variables
func makeFileVars(path string) {
	SetVar("TemplateFileName", path)
}

// TODO add more template variables
func makeRunVars() {
	SetVar("Timestamp", time.Now().String())
}
