package main

import (
	"log"
	"sync"

	f "github.com/soerlemans/table/filter"
	u "github.com/soerlemans/table/util"
)

const (
	// Lets prevent an unlimited amount of goroutines being kicked off.
	MAX_TASKS = 10
)

// Orchestrate processing of different files.
func run(t_args Arguments) error {
	tables, err := readInput(t_args)
	if err != nil {
		return err
	}

	filter, err := f.InitFilter(t_args.ProgramText)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	// TODO: Use taskCounter to decide on a max taskcount.
	// var taskCounter uint64 = 0

	for _, table := range tables {
		// Use a lambda to capture localized variables.
		task := func() {
			// Create the context for the task to run.
			ctx := initProcessContext(&filter, table)

			// Start processing.
			err := Process(ctx)
			if err != nil {
				log.Fatalln(err)
			}

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
	u.FailIf(err)

	err = run(args)
	u.FailIf(err)
}
