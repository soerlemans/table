package main

// Regular Imports:
import (
	"fmt"
	"github.com/alexflint/go-arg"
	"os"

	util "github.com/soerlemans/table/util"
)

// Arguments Struct:
type Arguments struct {
	ProgramFile string `arg:"-f,--file" help:"Path to file containing queries."`

	// Positional.
	ProgramText string   `arg:"positional" help:"Query to execute."`
	InputFiles  []string `arg:"positional" help:"Files to source as input."`
}

// Arguments Methods:
func (Arguments) Version() string {
	return fmt.Sprintf("Version: %s", VERSION)
}

// Globals:
const (
	VERSION = "0.1"
)

// Cli arguments variable.
var args Arguments

// Functions:
func initArgs() error {
	// Parse and handle arguments.
	arg.MustParse(&args)

	// Logging:
	if len(args.ProgramFile) != 0 {
		// If a
		args.InputFiles = append([]string{args.ProgramText}, args.InputFiles...)
	}

	if len(args.ProgramText) != 0 {
	}

	util.Printf("args: %+v", args)
	return nil
}
