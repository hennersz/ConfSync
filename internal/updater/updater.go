package updater

import (
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/hennersz/ConfSync/internal/operators"
)

const ConfigFileName = "config.json"

type TaskConfig struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

type Config map[string][]TaskConfig

type Updater struct {
	sourceDir    string
	config       Config
	operators    map[string]operators.Operator
	operatorKeys []string
}

type UpdaterBuilder struct {
	sourceDir string
	operators []operators.Operator
}

func (b *UpdaterBuilder) SrcDir(sourceDir string) *UpdaterBuilder {
	b.sourceDir = sourceDir
	return b
}

func (b *UpdaterBuilder) WithOperator(operator operators.Operator) *UpdaterBuilder {
	b.operators = append(b.operators, operator)
	return b
}

func (b *UpdaterBuilder) Build() (*Updater, error) {
	confFileData, err := ioutil.ReadFile(path.Join(b.sourceDir, ConfigFileName))
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = json.Unmarshal(confFileData, &config)
	if err != nil {
		return nil, err
	}

	operators := make(map[string]operators.Operator)
	operatorKeys := make([]string, 0, len(b.operators))
	for _, operator := range b.operators {
		key := operator.Name()
		operators[key] = operator
		operatorKeys = append(operatorKeys, key)
	}

	return &Updater{b.sourceDir, config, operators, operatorKeys}, nil
}

func New() *UpdaterBuilder {
	return &UpdaterBuilder{}
}

func (u *Updater) Update() error {
	for filePath, taskList := range u.config {
		for _, task := range taskList {
			err := u.submitTask(path.Join(u.sourceDir, filePath), task)
			if err != nil {
				return err
			}
		}
	}

	for _, operator := range u.operators {
		err := operator.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

type OperatorNotFoundError struct {
	OperatorName string
}

func (e OperatorNotFoundError) Error() string {
	return "Operator: " + e.OperatorName + " was not found"
}

func (u *Updater) submitTask(filePath string, config TaskConfig) error {
	operator, ok := u.operators[config.Name]

	if !ok {
		return OperatorNotFoundError{config.Name}
	}

	fileName := path.Base(filePath)

	return operator.SubmitTask(fileName, config.Args)
}
