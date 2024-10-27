package extractor

import (
	"github.com/go-rod/rod"
)

type Schema struct {
	Name         string  `json:"name"`
	BaseSelector string  `json:"baseSelector"`
	Fields       []Field `json:"fields"`
}

type Field struct {
	Name      string  `json:"name"`
	Selector  string  `json:"selector"`
	Type      string  `json:"type"`
	Attribute string  `json:"attribute,omitempty"`
	Fields    []Field `json:"fields,omitempty"`
}

type Extractor struct {
	Schema  Schema
	Browser *rod.Browser
}

type ExtractedItem map[string]interface{}
