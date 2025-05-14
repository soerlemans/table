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
	ProgramFile    string `arg:"-f,--file" help:"Path to file containing filters."`
	FromStdin      bool   `arg:"--stdin" help:"Specifies if the program should read from stdin." default:"false"`
	FieldSeparator rune   `arg:"-F,--field-separator" help:"Define the field separator."`
	Csv            bool   `help:"Define that input is CSV."`
	Json           bool   `help:"Define that input is JSON."`
	Excel          bool   `help:"Define that input is Excel. "`

	// Positional.
	ProgramText string   `arg:"positional" help:"Filter to execute."`
	InputFiles  []string `arg:"positional" help:"Files to source as input."`
}

// Arguments Methods:
func (Arguments) Version() string {
	return fmt.Sprintf("Version: %s", VERSION)
}

// Globals:
const (
	VERSION = "0.1"

	DEFAULT_FILTER = "."
)

// Functions:
func handleProgramFile(t_args *Arguments) error {
	if len(t_args.ProgramFile) != 0 {
		// If the program file was supplied then the program text positional arg.
		// Should be moved to other input files to parse.
		slice := []string{t_args.ProgramText}
		t_args.InputFiles = append(slice, t_args.InputFiles...)

		_, err := os.Stat(t_args.ProgramFile)
		if err != nil {
			return err
		}

		// TODO: Implement the rest.
		t_args.ProgramText = "TODO: Overwrite with args.ProgramFile content."
	}

	return nil
}

func initArgs() (Arguments, error) {
	var args Arguments

	// Parse and handle arguments.
	arg.MustParse(&args)
	defer util.Logf("args: %+v", args)

	// Logging:
	err := handleProgramFile(&args)
	if err != nil {
		return args, err
	}

	if len(args.ProgramText) == 0 {
		defaultFilter := util.Quote(DEFAULT_FILTER)
		util.Logf("No program text given, apply default: %s", defaultFilter)
		args.ProgramText = DEFAULT_FILTER
	}

	// If no input files are supplied check stdin.
	if len(args.InputFiles) == 0 {
		args.FromStdin = true
	}

	return args, nil
}
