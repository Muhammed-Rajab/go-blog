package utils

import (
	"html/template"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func formatAsDate(t time.Time) string {
	return t.Format("Monday, 2 May 2006")
}

func tagsAsString(tags []string) string {
	return strings.Join(tags, ",")
}

func objectID(id primitive.ObjectID) string {
	return id.Hex()
}

func GetTemplateFuncsMap() template.FuncMap {
	return template.FuncMap{
		"formatAsDate": formatAsDate,
		"objectid":     objectID,
		"tagsAsString": tagsAsString,
	}
}
