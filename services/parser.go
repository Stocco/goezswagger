package services

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var (
	yamlOutput *YamlOutput
	golangNumeric map[string]string
)

func init() {
	golangNumeric = make(map[string]string)
	golangNumeric["int"] = "integer"
	golangNumeric["int8"] = "integer"
	golangNumeric["int16"] = "integer"
	golangNumeric["int32"] = "integer"
	golangNumeric["int64"] = "integer"
	golangNumeric["uint8"] = "integer"
	golangNumeric["uint16"] = "integer"
	golangNumeric["uint32"] = "integer"
	golangNumeric["uint64"] = "integer"
	golangNumeric["float64"] = "number"
	golangNumeric["float32"] = "number"

}

func GenerateFile(filePath string) {

	yamlOutput = &YamlOutput{Openapi: "3.0.0", Info: &Info{}, Paths: make(map[string] map[string] *Method), Components: &Components{Schema: make(map[string] *ModelSchema)}}

	log.Println(fmt.Sprintf("Parsing main dir ..."))
	parseHeaders(filePath)
	getAnnotatedModels(filePath)
	parseRoutesFromFile(filePath)

	dirs, err := ioutil.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dirs {
		if file.IsDir() && file.Name()[:1] != "." && file.Name() != "vendor" {
			log.Println(fmt.Sprintf("Parsing dir: %s ...", file.Name()))

			parseHeaders(filePath + "/" + file.Name())
			getAnnotatedModels(filePath + "/" + file.Name())
			parseRoutesFromFile(filePath + "/" + file.Name())

		}
	}

	bytes, err := yaml.Marshal(yamlOutput)
	if err != nil {
		log.Fatal(err.Error())
	}


	err = ioutil.WriteFile("./generated_swagger.yaml", bytes, 0400)
	if err != nil {
		log.Fatal(err.Error())
	}


	log.Println(fmt.Sprintf("Parse completed. File generate at project dir: generated_swagger.yaml"))

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
			reqName := strings.TrimSpace(comments[pos].Text[len("//@request")+1:])
			if strings.Contains(reqName, "[]") {
				requestWithoutArray := strings.Replace(reqName, "[]", "", 1)
				method.RequestBody = &RequestBody{Content: &Content{ApplicationType: &ApplicationType{Schema: &Schema{Type: "array", SchemaItems: &SchemaItems{Ref: `#/components/schemas/` + requestWithoutArray}}}}}
			} else {
				method.RequestBody = &RequestBody{Content: &Content{ApplicationType: &ApplicationType{Schema: &Schema{Ref: `#/components/schemas/` + reqName}}}}
			}

		}


		if pos := findKeyInComments("//@response", comments) ; pos > 0 {

			method.Responses = make(map[string] *Response)
			for _, response := range strings.Split(comments[pos].Text[len("//@response") + 1:], " ") {
				resp := strings.Split(response, ":")

				if strings.Contains(resp[1], "[]") {
					responseWithoutArray := strings.Replace(resp[1], "[]", "", 1)
					method.Responses[resp[0]] = &Response{Description: ``, Content: &Content{ApplicationType: &ApplicationType{Schema: &Schema{Type: "array", SchemaItems:&SchemaItems{Ref: `#/components/schemas/` + responseWithoutArray}}}}}
				} else {
					method.Responses[resp[0]] = &Response{Description: ``, Content: &Content{ApplicationType: &ApplicationType{Schema: &Schema{Ref: `#/components/schemas/` + resp[1]}}}}
				}



			}
		}

		if method.Responses == nil {
			method.Responses = make(map[string] *Response)
			method.Responses[`200`]  = &Response{Description: ``}
		}

		if yamlOutput.Paths[path] == nil {
			yamlOutput.Paths[path] = make(map[string] *Method)
			yamlOutput.Paths[path][methodName] = method
		} else {
			yamlOutput.Paths[path][methodName] = method
		}

	}


}



func getAnnotatedModels(fileName string) {
	fset := token.NewFileSet()
	astr, err := parser.ParseDir(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, astt := range astr {
		for _, f := range astt.Files  {
			for _, modelName := range retrieveAnnotatedNames("//@model", f) {
				log.Println(fmt.Sprintf("Found model %s", modelName))
				parseAnottatedModel(modelName, fileName)
			}
		}
	}
}

func retrieveAnnotatedNames(key string, ast *ast.File) []string {

	var modelNames []string
	for _, cg := range ast.Comments {

		for _, v := range cg.List {

			if strings.Contains(strings.ToLower(v.Text), key) {
				modelNames = append(modelNames, strings.TrimSpace(v.Text[len(key):]))
			}

		}

	}

	return modelNames
}

func parseAnottatedModel(modelName string, fileName string) {

	fset := token.NewFileSet()
	astr, err := parser.ParseDir(fset, fileName, nil, 0)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, astt := range astr {
		for _, f := range astt.Files  {
			for _, model := range f.Scope.Objects {
				if model.Name == modelName {

					decl := model.Decl.(*ast.TypeSpec)
					structDecl, ok := decl.Type.(*ast.StructType)
					if !ok {
						return
					}
					fields := structDecl.Fields.List

					for _, field := range fields {
						if field.Tag != nil && strings.Contains(field.Tag.Value, "json:") {

							safeCast, okIdent := field.Type.(*ast.Ident)
							safePointerCast, okStarExpr := field.Type.(*ast.StarExpr)
							safeMapCast, okMap := field.Type.(*ast.MapType)
							safeArrayCast, okArray := field.Type.(*ast.ArrayType)
							safeSelectorCast, okSelector := field.Type.(*ast.SelectorExpr)
							safeInterfaceCast, okInterface := field.Type.(*ast.InterfaceType)

							if !okIdent && !okStarExpr && !okMap && !okArray && !okSelector && !okInterface {
								return
							}

							var fieldType string
							if safeCast != nil {
								fieldType = safeCast.Name
							} else if safePointerCast  != nil {

								switch v := safePointerCast.X.(type) {
								case *ast.ArrayType:
									fieldType = "array"
									safeArrayCast = v
								case *ast.Ident:
									fieldType = v.Name
								case *ast.SelectorExpr:
									fieldType = "string"
								}

							} else if safeMapCast != nil {
								fieldType = "object"
							} else if safeArrayCast != nil {
								fieldType = "array"
							} else if safeSelectorCast != nil {
								fieldType = "string"
							} else if safeInterfaceCast != nil {
								fieldType = "object"
							}


							builtin := isBasicType(fieldType)
							fieldType = correctFileTypeName(fieldType)

							if yamlOutput.Components.Schema[model.Name] == nil {
								yamlOutput.Components.Schema[model.Name] = &ModelSchema{Properties: make(map[string] *Properties), Example: make(map[string] interface{})}
							}

							if builtin {
								yamlOutput.Components.Schema[model.Name].Properties[extractKey(field.Tag.Value, "json")] = &Properties{Type: fieldType, Description: extractKey(field.Tag.Value, "description")}
							} else if fieldType == "array" {

								switch v := safeArrayCast.Elt.(type) {
								case *ast.StarExpr:
									yamlOutput.Components.Schema[model.Name].Properties[extractKey(field.Tag.Value, "json")] = &Properties{Type: fieldType, Items: &Schema{Ref: `#/components/schemas/` + v.X.(*ast.Ident).Name}, Description: extractKey(field.Tag.Value, "description")}
								case *ast.Ident:
									if isBasicType(v.Name) {
										yamlOutput.Components.Schema[model.Name].Properties[extractKey(field.Tag.Value, "json")] = &Properties{Type: fieldType, Items: &Schema{Type: correctFileTypeName(v.Name)}, Description: extractKey(field.Tag.Value, "description")}
									} else {
										yamlOutput.Components.Schema[model.Name].Properties[extractKey(field.Tag.Value, "json")] = &Properties{Type: fieldType, Items: &Schema{Ref: `#/components/schemas/` + v.Name}, Description: extractKey(field.Tag.Value, "description")}
									}
								case *ast.InterfaceType:
									yamlOutput.Components.Schema[model.Name].Properties[extractKey(field.Tag.Value, "json")] = &Properties{Type: fieldType, Items: &Schema{Type: "object"}, Description: extractKey(field.Tag.Value, "description")}
								case *ast.MapType:
									yamlOutput.Components.Schema[model.Name].Properties[extractKey(field.Tag.Value, "json")] = &Properties{Type: fieldType, Items: &Schema{Type: "object"}, Description: extractKey(field.Tag.Value, "description")}
								case *ast.SelectorExpr:
									yamlOutput.Components.Schema[model.Name].Properties[extractKey(field.Tag.Value, "json")] = &Properties{Type: fieldType, Items: &Schema{Type: "string"}, Description: extractKey(field.Tag.Value, "description")}
								}


							} else if fieldType == "object" {
								yamlOutput.Components.Schema[model.Name].Properties[extractKey(field.Tag.Value, "json")] = &Properties{Type: fieldType, Description: extractKey(field.Tag.Value, "description")}
							} else {
								yamlOutput.Components.Schema[model.Name].Properties[extractKey(field.Tag.Value, "json")] = &Properties{Ref: `#/components/schemas/` + fieldType, Description: extractKey(field.Tag.Value, "description")}
							}

							if strings.Contains(field.Tag.Value, "example:") {
								exampleValue := extractKey(field.Tag.Value, "example")

								if exampleValue == "true" || exampleValue == "false" {
									yamlOutput.Components.Schema[model.Name].Example[extractKey(field.Tag.Value, "json")] = exampleValue == "true"
									continue
								}

								if num, err := strconv.Atoi(exampleValue) ; err == nil {
									yamlOutput.Components.Schema[model.Name].Example[extractKey(field.Tag.Value, "json")] = num
									continue
								}

								if fl, err := strconv.ParseFloat(exampleValue, 10); err == nil {
									yamlOutput.Components.Schema[model.Name].Example[extractKey(field.Tag.Value, "json")] = fl
									continue
								}

								yamlOutput.Components.Schema[model.Name].Example[extractKey(field.Tag.Value, "json")] = exampleValue
							}

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
		if fieldTag[pos:pos + 1] == `"` || fieldTag[pos:pos + 1] == `,` {
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

func isBasicType(typeString string) bool {

	for _, v := range types.Typ {
		if typeString == v.Name() {
			return true
		}
	}

	return false
}

func correctFileTypeName(fieldType string) string {

	correctFieldName := fieldType

	numberSwagger, isNumber := golangNumeric[fieldType]
	if isNumber {
		correctFieldName = numberSwagger
	}

	if correctFieldName == "bool" {
		correctFieldName = "boolean"
	}

	return correctFieldName

}