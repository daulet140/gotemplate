package files

import (
	"fmt"
	"os"
	"text/template"
)

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
