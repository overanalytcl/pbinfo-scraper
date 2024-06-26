package pbgen

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

type Example struct {
	Name        string
	Input       string
	Output      string
	Explanation string
}

type StatementDetails struct {
	InputData    string
	OutputData   string
	Restrictions string
	Examples     []Example
}

func NewStatementDetails() *StatementDetails {
	return &StatementDetails{}
}

func (d *StatementDetails) Parse(pd *ProblemDetails) *StatementDetails {
	// A no-op, for now
	return d
}

type Problem struct {
	Metadata ProblemMetadata
	Details  StatementDetails
}

func NewProblem(d *ProblemDetails) *Problem {
	details := NewStatementDetails().Parse(d)
	return &Problem{
		Metadata: *NewProblemMetadata(d),
		Details:  *details,
	}
}

func NewProblemFromId(id int) (*Problem, error) {
	details, err := NewProblemDetails(id)
	if err != nil {
		return nil, err
	}

	return NewProblem(details), nil
}

//go:embed problem.tmpl
var problemTemplate string

func (p *Problem) ToMarkdown() (string, error) {
	tmpl := template.New("problem")
	tmpl = tmpl.Funcs(template.FuncMap{
		"MetadataToMarkdown": func(metadata ProblemMetadata) (string, error) {
			md, err := metadata.ToMarkdown()
			if err != nil {
				return "", err
			}
			return md, nil
		},
	})

	tmpl, err := tmpl.Parse(problemTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, p)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return "", err
	}

	return buf.String(), nil
}
