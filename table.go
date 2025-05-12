package main

import (
	"io"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/soerlemans/table/util"
)

const (
	// Lets prevent an unlimited amount of goroutines being kicked off.
	MAX_TASKS = 10

	NEWLINE = "\n"
)

// Convert a byte blob into a text array.
func bytes2Text(t_input []byte) []string {
	text := strings.Split(string(t_input), NEWLINE)

	// We should only strip the last line if we have at least one line.
	if len(text) > 1 {
		// Strip any empty lines from the end.
		for index, line := range slices.Backward(text) {
			if len(line) != 0 {
				text = text[:index + 1]
				break
			}
		}
	}

	return text
}

func readInputStdin(t_args Arguments) ([]Context, error) {
	var ctxList []Context

	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return ctxList, err
	}

	input := bytes2Text(content)
	ctx := initContext(input)
	ctxList = append(ctxList, ctx)

	return ctxList, nil
}

func readInputFiles(t_args Arguments) ([]Context, error) {
	var ctxList []Context

	for _, filePath := range t_args.InputFiles {
		content, err := os.ReadFile(filePath)

		if err != nil {
			return ctxList, err
		}

		input := bytes2Text(content)
		ctx := initContext(input)
		ctxList = append(ctxList, ctx)
	}

	return ctxList, nil
}

func readInput(t_args Arguments) ([]Context, error) {
	var (
		ctxList []Context
		err     error
	)

	if t_args.FromStdin {
		ctxList, err = readInputStdin(t_args)
	} else {
		ctxList, err = readInputFiles(t_args)
	}

	return ctxList, err
}

// Orchestrate processing of different files.
func run(t_args Arguments) error {
	ctxList, err := readInput(t_args)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	// TODO: Use taskCounter to decide on a max taskcount.
	// var taskCounter uint64 = 0

	for _, ctx := range ctxList {
		// Use a lambda to capture localized variables.
		task := func() {
			Process(ctx)

			wg.Done()
		}

		// Increment wait count.
		wg.Add(1)

		// Launch async task.
		go task()

	}

	// Wait for all goroutines to finish.
	wg.Wait()

	return nil
}

func main() {
	args, err := initArgs()
	util.FailIf(err)

	err = run(args)
	util.FailIf(err)
}
