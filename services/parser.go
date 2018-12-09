package services

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

var (
	yamlOutput *YamlOutput
	modelsUsed map[string] bool
)


func GenerateFile(filePath string) {

	yamlOutput = &YamlOutput{Openapi: "3.0.0", Info: &Info{}, Paths: make(map[string] map[string] *Method), Components: &Components{Schema: make(map[string] *ModelSchema)}}
	modelsUsed = make(map[string]bool)

	log.Println(fmt.Sprintf("Parsing main dir ..."))
	parseHeaders(filePath)
	parseModelsFromFile(filePath)
	parseRoutesFromFile(filePath)

	dirs, err := ioutil.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dirs {
		if file.IsDir() && file.Name()[:1] != "." && file.Name() != "vendor" {
			log.Println(fmt.Sprintf("Parsing dir: %s ...", file.Name()))

			parseHeaders(filePath + "/" + file.Name())
			parseModelsFromFile(filePath + "/" + file.Name())
			parseRoutesFromFile(filePath + "/" + file.Name())

		}
	}

	for k, _ := range yamlOutput.Components.Schema {
		_, used := modelsUsed[k]
		if !used {
			delete(yamlOutput.Components.Schema, k)
		}
	}

	bytes, err := yaml.Marshal(yamlOutput)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = ioutil.WriteFile("./generated_swagger.yaml", bytes, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}


	log.Println(fmt.Sprintf("Parse completed. Models used [%v] File generate at project dir: generated_swagger.yaml", modelsUsed))

}

func parseRoutesFromFile(fileName string) {

	fset := token.NewFileSet()
	astr, err := parser.ParseDir(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, astt := range astr {
		for _, f := range astt.Files {

			for _, cg := range f.Comments {
				checkIfRouteComentGroup(cg.List)
			}


		}
	}

}
func checkIfRouteComentGroup(comments []*ast.Comment) {

	if len(comments) < 3 {
		return
	}


	if strings.Contains(comments[0].Text, "//@path") {

		method := &Method{}
		path := strings.TrimSpace(comments[0].Text[len("//@path"):])
		methodName := strings.TrimSpace(comments[1].Text[len("//@method"):])
		method.Summary = strings.TrimSpace(comments[2].Text[len("//@summary"):])


		if pos := findKeyInComments("//@tags", comments) ; pos > 0 {
			method.Tags = strings.Split(comments[pos].Text[len("//@tags")+1:], " ")
		}

		if pos := findKeyInComments("//@request", comments) ; pos > 0 {
			method.RequestBody = &RequestBody{Content: &Content{ApplicationType: &ApplicationType{Schema: &Schema{Ref: `#/components/schemas/` + strings.TrimSpace(comments[pos].Text[len("//@request")+1:])}}}}
			modelsUsed[strings.TrimSpace(comments[pos].Text[len("//@request")+1:])] = true
		}


		if pos := findKeyInComments("//@response", comments) ; pos > 0 {

			method.Responses = make(map[string] *Response)
			for _, response := range strings.Split(comments[pos].Text[len("//@response") + 1:], " ") {
				resp := strings.Split(response, ":")
				respStruct := &Response{Description: `ok`, Content: &Content{ApplicationType: &ApplicationType{Schema: &Schema{Ref: `#/components/schemas/` + resp[1]}}}}
				method.Responses[resp[0]] = respStruct
				modelsUsed[resp[1]] = true

			}
		}

		if method.Responses == nil {
			method.Responses = make(map[string] *Response)
			method.Responses[`200`]  = &Response{Description: `ok`}
		}

		if yamlOutput.Paths[path] == nil {
			yamlOutput.Paths[path] = make(map[string] *Method)
			yamlOutput.Paths[path][methodName] = method
		} else {
			yamlOutput.Paths[path][methodName] = method
		}

	}


}



func parseModelsFromFile(fileName string) {
	fset := token.NewFileSet()
	astr, err := parser.ParseDir(fset, fileName, nil, 0)
	if err != nil {
		log.Fatal(err.Error())
	}


	for _, astt := range astr {
		for _, f := range astt.Files  {

			for _, model := range f.Scope.Objects {

				if model.Kind.String() == "type" {

					decl := model.Decl.(*ast.TypeSpec)
					structDecl, ok := decl.Type.(*ast.StructType)
					if !ok {
						return
					}
					fields := structDecl.Fields.List

					for _, field := range fields {

						if field.Tag != nil && strings.Contains(field.Tag.Value, "json:") {

							safeCast, ok := field.Type.(*ast.Ident)
							if !ok {
								return
							}

							fieldType := safeCast.Name
							if fieldType == "int" || fieldType == "float64" || fieldType == "float32"  {
								fieldType = "number"
							}

							if yamlOutput.Components.Schema[model.Name] == nil {
								yamlOutput.Components.Schema[model.Name] = &ModelSchema{Properties: make(map[string] *Properties)}
							}

							yamlOutput.Components.Schema[model.Name].Properties[extractKey(field.Tag.Value, "json")] = &Properties{Type: fieldType, Description: extractKey(field.Tag.Value, "description")}
						}

					}


				}
			}
		}
	}
}
func extractKey(fieldTag string, key string) string {

	if strings.Index(fieldTag, key) < 0 {
		return ""
	}

	value := ""
	for pos := strings.Index(fieldTag, key) + len(key) + 2; pos < len(fieldTag) ; pos++ {
		if fieldTag[pos:pos + 1] == `"` {
			return value
		}

		value = value + fieldTag[pos:pos + 1]
	}

	return  value
}


func parseHeaders(fileName string) {
	fset := token.NewFileSet()
	astr, err := parser.ParseDir(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, v := range astr {
		for _, v := range v.Files {

			if yamlOutput.Info.Title == "" {
				yamlOutput.Info.Title = strings.TrimSpace(retrieveKeyFromAstFile("//@title", v))
			}

			if yamlOutput.Info.Version == "" {
				yamlOutput.Info.Version = strings.TrimSpace(retrieveKeyFromAstFile("//@version", v))
			}

			if yamlOutput.Info.Description == "" {
				yamlOutput.Info.Description = strings.TrimSpace(retrieveKeyFromAstFile("//@description", v))
			}

		}
	}



}

func findKeyInComments(key string, comments []*ast.Comment) int {

	for c, v := range comments {
		if strings.Contains(v.Text, key) {
			return c
		}
	}

	return -1

}

func retrieveKeyFromAstFile(key string, ast *ast.File) string {

	for _, cg := range ast.Comments {

		for _, v := range cg.List {

			if strings.Contains(v.Text, key) {
				return v.Text[len(key):]
			}

		}

	}

	return ""
}