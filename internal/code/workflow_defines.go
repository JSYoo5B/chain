package code

import (
	"bytes"
	"fmt"
	"github.com/JSYoo5B/chain"
	"text/template"
)

type WorkflowDef struct {
	ConstructName   string
	ConstructParams string

	NodeConstructors []string

	WorkflowName string
	WorkflowType string

	Nodes []WorkflowNode
	Edges []WorkflowEdge
}

type WorkflowNode struct {
	VarName       string
	ConstructExpr string
}

type WorkflowEdge struct {
	BaseNode string
	WorkType string
	Plan     map[string]string
}

func (e WorkflowEdge) ActionPlanCode() string {
	if len(e.Plan) == 0 {
		return fmt.Sprintf("chain.TerminationPlan[%s]()",
			e.WorkType)
	} else if len(e.Plan) == 1 &&
		e.Plan[chain.Success] != "" {
		return fmt.Sprintf("chain.SuccessOnlyPlan(%s)",
			e.Plan[chain.Success])
	} else if len(e.Plan) == 2 &&
		e.Plan[chain.Success] != "" && e.Plan[chain.Error] != "" {
		return fmt.Sprintf("chain.DefaultPlan(%s, %s)",
			e.Plan[chain.Success],
			e.Plan[chain.Error])
	} else if len(e.Plan) == 3 &&
		e.Plan[chain.Success] != "" && e.Plan[chain.Error] != "" && e.Plan[chain.Abort] != "" {
		return fmt.Sprintf("chain.DefaultPlanWithAbort(%s, %s, %s)",
			e.Plan[chain.Success],
			e.Plan[chain.Error],
			e.Plan[chain.Abort])
	}

	buf := bytes.Buffer{}
	if err := actionPlanTemplate.Execute(&buf, e); err != nil {
		return ""
	}
	return buf.String()
}

var (
	actionPlanTemplate *template.Template
)

func init() {
	tmpl := `chain.ActionPlan[{{.WorkType}}]{
        {{- range $key, $value := .Plan }}
        "{{ $key }}": {{ $value }},
        {{- end }}
    }`
	var err error
	if actionPlanTemplate, err = template.New("").Parse(tmpl); err != nil {
		panic(err)
	}
}
