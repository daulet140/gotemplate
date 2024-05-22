package folder

func createManagerDir(internalDir string) (string, error) {
	managerDir := internalDir + "/manager"
	err := generateDirectory(managerDir)
	if err != nil {
		return managerDir, err
	}

	v1Dir := managerDir + "/v1"
	err = generateDirectory(v1Dir)
	if err != nil {
		return v1Dir, err
	}
	return "", err
}
