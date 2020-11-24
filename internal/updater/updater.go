package updater

import (
	"encoding/json"
	"fmt"
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

type Builder struct {
	sourceDir string
	operators []operators.Operator
}

func (b *Builder) SrcDir(sourceDir string) *Builder {
	b.sourceDir = sourceDir

	return b
}

func (b *Builder) WithOperator(operator operators.Operator) *Builder {
	b.operators = append(b.operators, operator)

	return b
}

func (b *Builder) Build() (*Updater, error) {
	confFileData, err := ioutil.ReadFile(path.Join(b.sourceDir, ConfigFileName))
	if err != nil {
		return nil, fmt.Errorf("error occurred reading config file: %w", err)
	}

	config := Config{}

	err = json.Unmarshal(confFileData, &config)
	if err != nil {
		return nil, fmt.Errorf("error occurred parsing config: %w", err)
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

func New() *Builder {
	return &Builder{}
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
			return fmt.Errorf("error occurred running operator: %s\nerror: %w", operator.Name(), err)
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
		return OperatorNotFoundError{OperatorName: config.Name}
	}

	fileName := path.Base(filePath)

	return operator.SubmitTask(fileName, config.Args)
}
