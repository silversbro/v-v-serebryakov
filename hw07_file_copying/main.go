package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

//nolint:unused
func main() {
	flag.Parse()

	if err := validateArgs(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := Copy(from, to, offset, limit); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

//nolint:unused
func validateArgs() error {
	if from == "" || to == "" {
		return fmt.Errorf("both 'from' and 'to' paths must be specified")
	}

	if offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}

	if limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}

	return nil
}
