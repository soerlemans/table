package sink

import (
	"fmt"
	"os"
)

// Very straightforward implementation.
type FileSink struct {
	File *os.File
}

func (this *FileSink) Writef(t_fmt string, t_args ...interface{}) {
	str := fmt.Sprintf(t_fmt, t_args...)

	_, err := this.File.WriteString(str)
	if err != nil {
		// TODO: Handle err.
	}
}

func (this *FileSink) Writeln(t_args ...interface{}) {
	str := fmt.Sprintln(t_args...)

	_, err := this.File.WriteString(str)
	if err != nil {
		// TODO: Handle err.
	}
}

func InitFileSink(t_path string) (FileSink, error) {
	var sink FileSink

	file, err := os.Open(t_path)
	if err != nil {
		return sink, err
	}

	// Set the file sink.l
	sink.File = file

	return sink, nil
}
