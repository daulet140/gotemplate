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

	internalDir := "internal"
	err = generateDirectory(internalDir)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("directory created:  %s", internalDir)
	}

	modelsDir := fmt.Sprintf("%s/models", internalDir)
	err = generateDirectory(modelsDir)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("directory created:  %s", modelsDir)
	}

	databaseDir := fmt.Sprintf("%s/database", internalDir)
	err = generateDirectory(databaseDir)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("directory created:  %s", databaseDir)
	}

	managerDir := fmt.Sprintf("%s/manager", internalDir)
	err = generateDirectory(managerDir)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("directory created:  %s", managerDir)
	}

	managerV1Dir := fmt.Sprintf("%s/v1", managerDir)
	err = generateDirectory(managerV1Dir)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("directory created:  %s", managerV1Dir)
	}

	httpDir := fmt.Sprintf("%s/http", internalDir)
	err = generateDirectory(httpDir)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("directory created:  %s", httpDir)
	}

	httpDirV1 := fmt.Sprintf("%s/v1", httpDir)
	err = generateDirectory(httpDirV1)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("directory created:  %s", httpDirV1)
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
	modelFile := fmt.Sprintf(modelsDir+"/%s_%s.go", filename, data.LowerStructName)
	err = writeStructToFile(*data, modelFile)
	if err != nil {
		log.Printf("Error writing struct to file: %v", err)
	} else {
		log.Printf("file created: %s", modelFile)
	}

	file, err := generateExampleRepo(data)
	if err != nil {
		log.Printf("%s error creating to file err: %v", file, err)
	} else {
		log.Printf("file created: %s", file)
	}

	file, err = generateManager(data)
	if err != nil {
		log.Printf("%s error creating to file err: %v", file, err)
	} else {
		log.Printf("file created: %s", file)
	}

	httpDirV1Struct := fmt.Sprintf("%s/%s", httpDirV1, data.LowerStructName)
	err = generateDirectory(httpDirV1Struct)
	if err != nil {
		log.Printf("Error creating directory: %s", httpDirV1Struct)
	} else {
		log.Printf("directory created: %s", httpDirV1Struct)
	}

	file, err = generateV1Struct(data)
	if err != nil {
		log.Printf("%s error creating to file err: %v", file, err)
	} else {
		log.Printf("file created: %s", file)
	}

	fmt.Println("Struct generated and written to file successfully.")

	return "", err
}
