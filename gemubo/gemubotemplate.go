package gemubo

import "strings"

type Template struct {
	Name          string
	Content       string
	PlaeceHolders map[string]interface{}
}

func NewTemplate(name, content string) *Template {

	tokens := strings.Fields(content)
	placeHolders := make(map[string]interface{})

	for _, token := range tokens {
		if IsPlaceHolder(token) {
			placeHolders[token] = interface{}(nil)
		}
	}

	return &Template{
		Name:          name,
		Content:       content,
		PlaeceHolders: placeHolders,
	}
}

func IsPlaceHolder(s string) bool {
	if s == "" {
		return false
	}

	ch := s[0]

	return string(ch) == "$"
}
