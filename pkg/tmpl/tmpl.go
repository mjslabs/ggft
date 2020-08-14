package tmpl

import (
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/template"
)

// CreateFileFromTemplate - parse a template file and write out results
func CreateFileFromTemplate(input, output string) error {
	var (
		f            *os.File
		err          error
		tmplContents string
	)

	if f, err = os.Create(output); err != nil {
		return err
	}
	defer f.Close()

	if tmplContents, err = readTemplateFile(input); err != nil {
		return err
	}

	// Set file-specific variables
	makeFileVars(input)

	// Write out the template
	return writeTemplateToFile(f, tmplContents)
}

// ScanTemplateForVars - scan the template for variables that need to be defined
func ScanTemplateForVars(templateFile string) ([]string, error) {
	buf := &bytes.Buffer{}
	t, err := getTemplateObject(templateFile)
	r := regexp.MustCompile(`map has no entry for key "([^ ]+?)"`)
	varList := []string{}

	// Seed vars specific to this file
	makeFileVars(templateFile)
	// Our place to keep fake values for undefined vars, so we can iterate through the whole template
	tmpVars := tmplVars

	err = t.Execute(buf, tmpVars)
	for err != nil {
		if !strings.Contains(err.Error(), "map has no entry for key") {
			// Encountered an error that isn't an undefined variable
			return nil, err
		}

		// Grab name of undefined template variable
		matches := r.FindStringSubmatch(err.Error())
		// Save the undefined var twice, once to allow our next run through to succeed...
		tmpVars[matches[1]] = ""
		// ... and the second, to return to the caller
		varList = append(varList, matches[1])
		// Try again
		err = t.Execute(buf, tmpVars)
	}

	return varList, err
}

func readTemplateFile(templateFile string) (string, error) {
	var tmplContents []byte
	var err error
	if tmplContents, err = ioutil.ReadFile(templateFile); err != nil {
		return "", err
	}
	return string(tmplContents), nil
}

func getTemplateObject(templateFile string) (*template.Template, error) {
	var tmplContents string
	var err error
	if tmplContents, err = readTemplateFile(templateFile); err != nil {
		return nil, err
	}

	t := template.New("").Option("missingkey=error")
	return t.Funcs(funcMap).Parse(string(tmplContents))
}

func writeTemplateToFile(f *os.File, tmplContents string) error {
	t := template.New("").Option("missingkey=error")
	packageTemplate := template.Must(t.Funcs(funcMap).Parse(string(tmplContents)))
	return packageTemplate.Execute(f, tmplVars)
}
