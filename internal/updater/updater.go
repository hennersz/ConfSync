package updater

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

const ConfigFileName = "config.json"

type Config map[string]string

type Updater struct {
	sourceDir string
	config    Config
}

func NewUpdater(sourceDir string) (*Updater, error) {
	confFileData, err := ioutil.ReadFile(path.Join(sourceDir, ConfigFileName))
	if err != nil {
		return nil, err
	}

	config := make(Config)
	err = json.Unmarshal(confFileData, &config)
	if err != nil {
		return nil, err
	}

	return &Updater{sourceDir, config}, nil
}

func (u *Updater) Update() error {
	for fileName, dest := range u.config {
		sourcePath := path.Join(u.sourceDir, fileName)
		destPath := path.Join(dest, path.Base(fileName))

		fileData, err := ioutil.ReadFile(sourcePath)
		if err != nil {
			return err
		}

		err = os.MkdirAll(dest, 0755)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(destPath, fileData, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
