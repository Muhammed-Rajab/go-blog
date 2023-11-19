package utils

import (
	"html/template"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func formatAsDate(t time.Time) string {
	return t.Format("Monday, 2 May 2006")
}

func objectID(id primitive.ObjectID) string {
	return id.Hex()
}

func GetTemplateFuncsMap() template.FuncMap {
	return template.FuncMap{
		"formatAsDate": formatAsDate,
		"objectid":     objectID,
	}
}
