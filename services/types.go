package services

type YamlOutput struct {
	Info *Info
	Openapi string
	Paths map[string] map[string] *Method
	Components *Components
}

type Components struct{
	Schema map[string] *ModelSchema `yaml:"schemas"`
}

type ModelSchema struct {
	Properties map[string] *Properties
	Example   map[string] interface{} `yaml:"example,omitempty"`
}

type Properties struct {
	Type string `yaml:"type,omitempty"`
	Description string `yaml:"description,omitempty"`
	Ref  string `yaml:"$ref,omitempty"`
	Items *Schema `yaml:"items,omitempty"`
}


type Method struct {
	Summary string
	Tags []string
	RequestBody *RequestBody `yaml:"requestBody,omitempty"`
	Responses map[string] *Response
}

type RequestBody struct {
	Content *Content `yaml:"content"`
}

type Response struct {
	Description string
	Content *Content `yaml:"content,omitempty"`
}

type Content struct {
	ApplicationType *ApplicationType `yaml:"application/json"`
}

type ApplicationType struct {
	Schema *Schema `yaml:"schema"`
}

type Schema struct {
	Ref string  `yaml:"$ref,omitempty"`
	Type string `yaml:"type,omitempty"`
}


type Info struct {
	Title string
	Version string
	Description string
}
