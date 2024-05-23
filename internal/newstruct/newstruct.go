package newstruct

import (
	"fmt"
	"github.com/daulet140/gotemplate/internal/cfg"
	"log"
	"os"
	"strings"
)

type JSONStruct interface {
	GenerateStructFromJSON(jsonFile string) (string, error)
}

type jsonStruct struct {
	config cfg.Config
}

func NewJSONStruct(config cfg.Config) *jsonStruct {
	return &jsonStruct{config: config}
}

func (j *jsonStruct) GenerateStructFromJSON(jsonFile string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pwdArr := strings.Split(pwd, "\\")

	filename := strings.Split(jsonFile, ".")[0]

	jsonStructInternal := fmt.Sprintf("generated_%s_internal", filename)
	err = generateDirectory(jsonStructInternal)
	if err != nil {
		return jsonStructInternal, err
	}

	modelsDir := fmt.Sprintf("%s/models", jsonStructInternal)
	err = generateDirectory(modelsDir)
	if err != nil {
		return modelsDir, err
	}

	databaseDir := fmt.Sprintf("%s/database", jsonStructInternal)
	err = generateDirectory(databaseDir)
	if err != nil {
		return databaseDir, err
	}

	managerDir := fmt.Sprintf("%s/manager", jsonStructInternal)
	err = generateDirectory(managerDir)
	if err != nil {
		return managerDir, err
	}

	managerV1Dir := fmt.Sprintf("%s/v1", managerDir)
	err = generateDirectory(managerV1Dir)
	if err != nil {
		return managerV1Dir, err
	}

	httpDir := fmt.Sprintf("%s/http", jsonStructInternal)
	err = generateDirectory(httpDir)
	if err != nil {
		return httpDir, err
	}

	httpDirV1 := fmt.Sprintf("%s/v1", httpDir)
	err = generateDirectory(httpDirV1)
	if err != nil {
		return httpDirV1, err
	}

	jsonData, err := readJSONFromFile(jsonFile)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	structName := "" // Starting struct name
	structDefinitions := []StructDefinition{}
	_, err = generateStructFromJSON(jsonData, structName, &structDefinitions)
	if err != nil {
		log.Fatalf("Error generating struct: %v", err)
	}

	mainStruct := structDefinitions[len(structDefinitions)-1].Name
	data := &Data{
		ProjectPath:       pwdArr[len(pwdArr)-1],
		FileName:          filename,
		StructName:        mainStruct,
		StructDefinitions: structDefinitions,
		LowerStructName:   strings.ToLower(mainStruct),
	}
	err = writeStructToFile(*data, fmt.Sprintf(modelsDir+"/%s_%s.go", filename, strings.ToLower(mainStruct)))
	if err != nil {

		log.Fatalf("Error writing struct to file: %v", err)
	}

	file, err := generateExampleRepo(data)
	if err != nil {
		log.Fatalf("%s Error writing struct to file: %v", file, err)
	}

	file, err = generateManager(data)
	if err != nil {
		log.Fatalf("%s Error writing struct to file: %v", file, err)
	}

	httpDirV1Struct := fmt.Sprintf("%s/%s", httpDirV1, data.LowerStructName)
	err = generateDirectory(httpDirV1Struct)
	if err != nil {
		return httpDirV1Struct, err
	}

	file, err = generateV1Struct(data)
	if err != nil {
		log.Fatalf("%s Error writing struct to file: %v", file, err)
	}

	fmt.Println("Struct generated and written to file successfully.")

	return "", err
}
