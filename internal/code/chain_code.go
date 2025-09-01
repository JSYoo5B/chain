package code

import (
	_ "embed"
	"io"
	"text/template"
)

//go:embed chain_code.tmpl
var chainCodeTmpl string
var chainCodeTemplate *template.Template

func init() {
	var err error
	if chainCodeTemplate, err = template.New("chainCode").Parse(chainCodeTmpl); err != nil {
		panic(err)
	}
}

type ChainCode struct {
	PackageName string
	Imports     []string
	Workflows   []WorkflowDef
}

func (c ChainCode) GenerateChainCode(wr io.Writer) error {
	return chainCodeTemplate.Execute(wr, c)
}
