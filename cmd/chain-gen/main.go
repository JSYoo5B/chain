// Command chain-gen translates one .chain file into Go workflow constructors.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/JSYoo5B/chain/internal/dsl/compiler"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("chain-gen", flag.ContinueOnError)
	flags.SetOutput(stderr)
	output := flags.String("o", "", "write generated Go to this file (default: stdout)")
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("usage: chain-gen [-o output.go] input.chain")
	}
	input := flags.Arg(0)
	if *output != "" {
		inputPath, err := filepath.Abs(input)
		if err != nil {
			return err
		}
		outputPath, err := filepath.Abs(*output)
		if err != nil {
			return err
		}
		if inputPath == outputPath {
			return fmt.Errorf("output must differ from the DSL input file")
		}
	}
	text, err := os.ReadFile(input)
	if err != nil {
		return err
	}
	generated, err := compiler.Compile(input, text)
	if err != nil {
		return err
	}
	if *output == "" {
		_, err = stdout.Write(generated)
		return err
	}
	// Compile before touching the output, and replace it only after a complete
	// write so an error cannot leave a truncated generated file behind.
	file, err := os.CreateTemp(filepath.Dir(*output), ".chain-gen-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(generated); err != nil {
		file.Close()
		return err
	}
	if err := file.Chmod(0644); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), *output)
}
