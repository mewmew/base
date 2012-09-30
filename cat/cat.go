package main

import "flag"
import "fmt"
import "io"
import "log"
import "os"

func init() {
	flag.Usage = usage
	flag.Parse()
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
	if flag.NArg() == 0 {
		// Read from stdin when no FILE has been provided.
		err := Cat(StdinFileName)
		if err != nil {
			log.Fatalln(err)
		}
	}
	for _, filePath := range flag.Args() {
		err := Cat(filePath)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// Cat outputs the content of a provided file or standard input (when the
// provided file path is "-").
func Cat(filePath string) (err error) {
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
	_, err = io.Copy(os.Stdout, fr)
	if err != nil {
		return err
	}
	return nil
}
