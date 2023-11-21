package utils

import "strings"

func TagsFromString(stringTag string) []string {
	tags := strings.Split(stringTag, ",")
	for i, tag := range tags {
		tags[i] = strings.ToLower(strings.TrimSpace(tag))
	}
	return tags
}
