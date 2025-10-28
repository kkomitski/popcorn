package main

import (
	"fmt"
	"os"
	BE "pop/backend"
)

func main() {
	// If a file argument is provided, run the file
	if len(os.Args) > 1 {
		filePath := os.Args[1]
		err := BE.RunFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running file: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Otherwise, start the REPL
		BE.Repl()
	}
}
