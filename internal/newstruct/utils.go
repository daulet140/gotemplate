package newstruct

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func generateDirectory(dirName string) error {
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		return err
	}
	return nil
}

type StructDefinition struct {
	Name       string
	Definition string
}
type Data struct {
	ProjectPath       string
	FileName          string
	StructDefinitions []StructDefinition
	StructName        string
	LowerStructName   string
}

// Function to generate struct definition from JSON
func generateStructFromJSON(jsonData map[string]interface{}, structName string, structDefinitions *[]StructDefinition) (string, error) {
	structFields := make([]string, 0)
	for key, value := range jsonData {

		fieldName := strings.Title(key)
		fieldType, nestedStruct := getGoTypeAndNestedStruct(value, strings.Title(key), structDefinitions)
		if nestedStruct != "" {
			*structDefinitions = append(*structDefinitions, StructDefinition{
				Name:       structName + strings.Title(key),
				Definition: nestedStruct,
			})
		}
		structFields = append(structFields, fmt.Sprintf("\t%s %s `json:\"%s\"`", fieldName, fieldType, key))

	}

	structDef := fmt.Sprintf("type %s struct {\n%s\n}", structName, strings.Join(structFields, "\n"))
	return structDef, nil
}

// Function to determine Go type from JSON value and handle nested structs
func getGoTypeAndNestedStruct(value interface{}, structName string, structDefinitions *[]StructDefinition) (string, string) {
	switch v := value.(type) {
	case string:
		return "string", ""
	case float64:
		return "float64", ""
	case bool:
		return "bool", ""
	case map[string]interface{}:
		nestedStruct, _ := generateStructFromJSON(v, structName, structDefinitions)
		return structName, nestedStruct
	case []interface{}:
		if len(v) > 0 {
			elemType, _ := getGoTypeAndNestedStruct(v[0], structName+"Item", structDefinitions)
			return "[]" + elemType, ""
		}
		return "[]interface{}", ""
	default:
		return "interface{}", ""
	}
}

// Function to write struct definitions to file using text/template
func writeStructToFile(structDefs Data, filename string) error {
	const structTemplate = `package models

{{range .StructDefinitions}}
{{.Definition}}
{{end}}`

	tmpl, err := template.New("struct").Parse(structTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, structDefs)
	if err != nil {
		return err
	}

	return nil
}
func readJSONFromFile(filename string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func generateFile(filePath string, tmpl string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	t := template.Must(template.New("").Parse(tmpl))
	err = t.Execute(file, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return err
	}
	return nil
}
