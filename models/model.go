package models

type FuncCase struct {
	Name  string
	Cases []*TestCase
}

type TestCase struct {
	Name  string
	Steps []*CaseStep
}

type CaseStep struct {
	Name        string
	Level       string
	Method      string
	URL         string
	Data        string
	ResType     string
	Verfication string

	Result *CaseOutput
}

type CaseOutput struct {
	Status  int
	Body    []byte
	Success bool
	Error   error
}
