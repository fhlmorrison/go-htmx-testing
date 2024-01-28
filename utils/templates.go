package utils

import (
	"bytes"
	"html/template"
)

func ExecuteTemplateToString(template *template.Template, name string, data any) string {
	var buffer bytes.Buffer
	template.ExecuteTemplate(&buffer, name, data)
	return buffer.String()
}
