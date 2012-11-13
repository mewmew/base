package main

import "flag"
import "fmt"
import "io"
import "log"
import "os"

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: cat [OPTION]... [FILE]...")
	fmt.Fprintln(os.Stderr, "Concatenate FILE(s), or standard input, to standard output.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "With no FILE, or when FILE is -, read standard input.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, "  cat f - g  Output f's contents, then standard input, then g's contents.")
	fmt.Fprintln(os.Stderr, "  cat        Copy standard input to standard output.")
}

// StdinFileName is a reserved file name used for standard input.
const StdinFileName = "-"

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		// Read from stdin when no FILE has been provided.
		err := cat(StdinFileName)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	for _, filePath := range flag.Args() {
		err := cat(filePath)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// cat outputs the content of a provided file or standard input (when the
// provided file path is "-").
func cat(filePath string) (err error) {
	// Open file.
	var fr *os.File
	if filePath == StdinFileName {
		fr = os.Stdin
	} else {
		fr, err = os.Open(filePath)
		if err != nil {
			return err
		}
		defer fr.Close()
	}

	// Write file contents to standard output.
	_, err = io.Copy(os.Stdout, fr)
	if err != nil {
		return err
	}
	return nil
}
