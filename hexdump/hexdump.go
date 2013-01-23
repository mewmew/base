package main

import "encoding/hex"
import "flag"
import "fmt"
import "io"
import "log"
import "os"

// flagLength is the number of bytes which should be interpreted. 0 corresponds
// to the entire input.
var flagLength int64

// flagOffset is the number of bytes which should be skipped from the beginning
// of the input.
var flagOffset int64

func init() {
	flag.Int64Var(&flagLength, "n", 0, "Interpret only x bytes of input.")
	flag.Int64Var(&flagOffset, "s", 0, "Skip x bytes from the beginning of the input.")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: hexdump [OPTION]... [FILE]...")
	fmt.Fprintln(os.Stderr, "Display file contents in hex and ASCII.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "With no FILE, or when FILE is -, read standard input.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, "  hexdump -n 0x10 f  Output a hex dump of the first 16 bytes of f's contents.")
	fmt.Fprintln(os.Stderr, "  hexdump -s 4 f  Skip the first 4 bytes and output a hex dump of f's contents.")
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		// Read from stdin when no FILE has been provided.
		err := hexdump(StdinFileName)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	for _, filePath := range flag.Args() {
		err := hexdump(filePath)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// StdinFileName is a reserved file name used for standard input.
const StdinFileName = "-"

// hexdump writes a hex dump of the provided file or standard input (when the
// provided file path is "-") to standard output. The format of the dump matches
// the output of `hexdump -C` on the command line.
func hexdump(filePath string) (err error) {
	// Open input file.
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

	// Ignore directories.
	fi, err := fr.Stat()
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return nil
	}

	if flagOffset != 0 {
		// Skip x bytes if "-s" flag is used.
		_, err = fr.Seek(flagOffset, os.SEEK_SET)
		if err != nil {
			return err
		}
	}

	// Write hex dump to standard output.
	w := hex.Dumper(os.Stdout)
	if flag.NArg() > 1 {
		// Output path if more than one input file.
		fmt.Println("path:", filePath)
		// Output new line after output from w.Close().
		defer fmt.Println()
	}
	defer w.Close()
	if flagLength != 0 {
		// Interpret only x bytes if "-n" flag is used.
		_, err = io.CopyN(w, fr, flagLength)
		if err != nil && err != io.EOF {
			return err
		}
	} else {
		_, err = io.Copy(w, fr)
		if err != nil {
			return err
		}
	}

	return nil
}
