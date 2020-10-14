package updater

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

const ConfigFileName = "config.json"

type Config map[string]string

type Updater struct {
	sourceDir string
	Config    Config
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
