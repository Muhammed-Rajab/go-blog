package utils

import (
	"html/template"
	"time"
)

func formatAsDate(t time.Time) string {
	return t.Format("Monday, 2 May 2006")
}

func GetTemplateFuncsMap() template.FuncMap {
	return template.FuncMap{
		"formatAsDate": formatAsDate,
	}
}
