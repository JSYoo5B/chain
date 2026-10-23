package compiler

import (
	"go/ast"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/JSYoo5B/chain/internal/dsl/frontend"
)

// Compile parses and analyzes a DSL file, then emits formatted Go source.
// Imports and Go type correctness are validated when the generated package is
// built; compilation never executes prerequisite code or action constructors.
func Compile(filename string, text []byte) ([]byte, error) {
	return CompileWithOptions(filename, text, Options{})
}

// Options controls the source path embedded in generated Go diagnostics.
type Options struct {
	// SourceName is relative to the generated file's directory. The default
	// assumes output beside the input and uses the input file's base name.
	SourceName string
}

// CompileWithOptions is Compile with an explicit generated-source path mapping.
func CompileWithOptions(filename string, text []byte, options Options) ([]byte, error) {
	file, err := frontend.Parse(filename, text)
	if err != nil {
		return nil, err
	}
	program, err := Analyze(file)
	if err != nil {
		return nil, err
	}
	if options.SourceName == "" {
		options.SourceName = filepath.Base(filename)
	}
	return generate(program, filepath.ToSlash(options.SourceName))
}

func generate(program *Program, sourceName string) ([]byte, error) {
	output := mappedOutput{source: program.File.Source, filename: sourceName}
	output.copy(program.File.HeaderSpan)
	if len(program.Workflows) == 0 {
		return output.formatted()
	}
	used := usedIdentifiers(program)
	runtimeName := freshName("_chain", used)
	output.printf(0, "\nimport %s %q\n\n", runtimeName, "github.com/JSYoo5B/chain")
	for _, workflow := range program.Workflows {
		w := workflow.Syntax
		typeSpan := frontend.Span{Start: w.Go.Position(w.Go.StateType.Pos()).Offset, End: w.Go.Position(w.Go.StateType.End()).Offset}
		// A parameter or prerequisite variable may shadow a predeclared type or
		// an imported package. Resolve the state type outside the function scope.
		state := freshName("_chainState_"+w.Constructor.Name, used)
		output.printf(typeSpan.Start, "type %s = ", state)
		output.copy(typeSpan)
		output.insert(typeSpan.End, "\n\n")
		local := freshName("_workflow", used)
		output.printf(w.Constructor.Span.Start, "func %s", w.Constructor.Name)
		output.copy(w.Parameters.Span)
		output.printf(typeSpan.Start, " *%s.Workflow[", runtimeName)
		output.copy(typeSpan)
		output.insert(typeSpan.End, "] {")
		if w.Prerequisite != nil {
			output.copy(w.Prerequisite.Content)
			// Preserve a statementless label at the end of the original block.
			// The newline also prevents a trailing // comment hiding generated code.
			output.insert(w.Prerequisite.Content.End, "\n;\n")
		}
		output.printf(w.Constructor.Span.Start, "%s := %s.NewWorkflow[%s](%s", local, runtimeName, state, strconv.Quote(w.Name.Name))
		for _, node := range workflow.Nodes {
			output.insert(node.Name.Span.Start, ", ")
			output.copy(node.Name.Span)
		}
		output.insert(w.Constructor.Span.Start, ")\n")
		for _, node := range workflow.Nodes {
			if len(node.Directions) == 0 {
				output.printf(node.Name.Span.Start, "%s.SetRunPlan(%s, %s.TerminationPlan[%s]())\n", local, node.Name.Name, runtimeName, state)
				continue
			}
			output.printf(node.Name.Span.Start, "%s.SetRunPlan(%s, %s.RunPlan[%s]{\n", local, node.Name.Name, runtimeName, state)
			directions := make([]string, 0, len(node.Directions))
			for direction := range node.Directions {
				directions = append(directions, direction)
			}
			sort.Strings(directions)
			for _, direction := range directions {
				target := node.Directions[direction]
				output.printf(target.Span.Start, "%s: ", strconv.Quote(direction))
				if target.Name == "end" {
					output.printf(target.Span.Start, "%s.Terminate[%s]()", runtimeName, state)
				} else {
					output.copy(target.Span)
				}
				output.insert(target.Span.End, ",\n")
			}
			output.insert(node.Name.Span.End, "})\n")
		}
		output.printf(w.Constructor.Span.Start, "return %s\n}\n\n", local)
	}
	return output.formatted()
}

// Avoid introducing an import or local variable that shadows user identifiers
// in parameters, state types, prerequisite code, or workflow declarations.
func usedIdentifiers(program *Program) map[string]bool {
	used := map[string]bool{}
	collect := func(root ast.Node) {
		ast.Inspect(root, func(node ast.Node) bool {
			if name, ok := node.(*ast.Ident); ok {
				used[name.Name] = true
			}
			return true
		})
	}
	collect(program.File.Header)
	for _, workflow := range program.Workflows {
		w := workflow.Syntax
		used[w.Constructor.Name], used[w.Name.Name] = true, true
		collect(w.Go.Parameters)
		collect(w.Go.StateType)
		if w.Go.Prerequisite != nil {
			collect(w.Go.Prerequisite)
		}
		for _, node := range workflow.Nodes {
			used[node.Name.Name] = true
		}
	}
	return used
}

func freshName(base string, used map[string]bool) string {
	name := base
	for i := 1; used[name]; i++ {
		name = base + strconv.Itoa(i)
	}
	used[name] = true
	return name
}
