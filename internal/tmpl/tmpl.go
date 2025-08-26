package tmpl

import (
	_ "embed"
	"html/template"
	"sync"
)

//go:embed main.html
var mainHTML string

var (
	parsedTmpl *template.Template
	once       sync.Once
	parseErr   error
)

func GetTemplate() (*template.Template, error) {
	once.Do(func() {
		parsedTmpl, parseErr = template.New("main").Parse(mainHTML)
	})
	if parseErr != nil {
		return nil, parseErr
	}
	return parsedTmpl, nil
}
