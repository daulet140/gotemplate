package folder

import "os"

func generateDirectory(dirName string) error {
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		return err
	}
	return nil
}
