package operators

type Operator interface {
	SubmitTask(string, []string) error
	Name() string
	Run() error
}
