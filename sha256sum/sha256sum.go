package main

import "crypto/sha256"
import "flag"
import "fmt"
import "io"
import "log"
import "os"

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: sha256sum [FILE]...")
	fmt.Fprintln(os.Stderr, "Print SHA256 checksums.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "With no FILE, or when FILE is -, read standard input.")
}

// StdinFileName is a reserved file name used for standard input.
const StdinFileName = "-"

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		// Read from stdin when no FILE has been provided.
		err := sha256sum(StdinFileName)
		if err != nil {
			log.Println(err)
		}
		return
	}

	for _, filePath := range flag.Args() {
		err := sha256sum(filePath)
		if err != nil {
			log.Println(err)
		}
	}
}

// sha256sum outputs the SHA256 checksum of a provided file or standard input
// (when the provided file path is "-").
func sha256sum(filePath string) (err error) {
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

	// Output SHA256 checksum.
	h := sha256.New()
	_, err = io.Copy(h, fr)
	if err != nil {
		return err
	}
	fmt.Printf("%x  %s\n", h.Sum(nil), filePath)

	return nil
}
