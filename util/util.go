// Utility functions for logging and formatting.
package util

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Specify an enumeration type for the input to the Etc function.
type EtcWidth int

const (
	// Terminal output should usually be max 80 characters.
	ETC80  EtcWidth = 80
	ETC120 EtcWidth = 120
	ETC160 EtcWidth = 160
)

// Determines if we should log (Not able to be a constant).
var DEBUG = initDebug()

const (
	// Exit code to give on failure.
	EXIT_CODE_ERR = 1
)

// Figure out if we should enable logging or not.
func initDebug() bool {
	var debug = false

	// Load environment from .env optionally.
	// Ignore any errors as a .env is not required.
	godotenv.Load()

	// Environment variable name to lookup.
	const DEBUG_ENV_NAME = "DEBUG"
	const DEBUG_FALSE = "false"

	value, exists := os.LookupEnv(DEBUG_ENV_NAME)
	if exists {

		// If the environment variable exists.
		// And is not false then we enable debug mode.
		debug = (value != DEBUG_FALSE)
	}

	if debug {
		log.Println("Debug mode on.")
	}

	return debug
}

// Convert any variable to a string and encircle with quotes.
func Quote[T any](t_var T) string {
	return fmt.Sprint("\"", t_var, "\"")

}

// Limit a string by a max character count, and append "..." (Etc is short for et cetera).
func Etc[T EtcWidth | int](t_str string, t_count T) string {
	const (
		etc    = "..."
		etcLen = len(etc)
	)

	maxStrlen := int(t_count) - etcLen
	sliced := t_str[:maxStrlen]

	return fmt.Sprintf("%s%s", sliced, etc)
}

// Fail unconditionally.
func Fail(t_err error) {
	log.Fatalln(t_err)

	os.Exit(EXIT_CODE_ERR)
}

// Fail if an error exists.
func FailIf(t_err error) {
	if t_err != nil {
		Fail(t_err)
	}
}

// Conditionally log only if DEBUG is set to true.
func Logf(t_fmt string, t_args ...interface{}) {
	if DEBUG {
		// Add a newline.
		fmtLn := fmt.Sprintln(t_fmt)

		log.Printf(fmtLn, t_args...)
	}
}

// Conditionally write either to logs or stdout depending on DEBUG.
func Printf(t_fmt string, t_args ...interface{}) {
	fmtLn := fmt.Sprintln(t_fmt)

	// A little more clear than !DEBUG.
	if DEBUG == false {
		// On non debug builds just print regularly.
		fmt.Printf(fmtLn, t_args...)
	} else {
		// Use log instead.
		log.Printf(fmtLn, t_args...)
	}
}

// Conditionally write either to logs or stdout depending on DEBUG.
func Println(t_args ...interface{}) {
	if !DEBUG {
		// On non debug builds just print regularly.
		fmt.Println(t_args...)
	} else {
		// Use log instead.
		log.Println(t_args...)
	}
}
