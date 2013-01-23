package main

import "crypto/md5"
import "flag"
import "fmt"
import "io"
import "log"
import "os"

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: md5sum [FILE]...")
	fmt.Fprintln(os.Stderr, "Print MD5 checksums.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "With no FILE, or when FILE is -, read standard input.")
}

// StdinFileName is a reserved file name used for standard input.
const StdinFileName = "-"

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		// Read from stdin when no FILE has been provided.
		err := md5sum(StdinFileName)
		if err != nil {
			log.Println(err)
		}
		return
	}

	for _, filePath := range flag.Args() {
		err := md5sum(filePath)
		if err != nil {
			log.Println(err)
		}
	}
}

// md5sum outputs the MD5 checksum of a provided file or standard input (when
// the provided file path is "-").
func md5sum(filePath string) (err error) {
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

	// Output MD5 checksum.
	h := md5.New()
	_, err = io.Copy(h, fr)
	if err != nil {
		return err
	}
	if filePath == StdinFileName {
		// don't output file path for standard input.
		fmt.Printf("%x\n", h.Sum(nil))
	} else {
		fmt.Printf("%x  %s\n", h.Sum(nil), filePath)
	}

	return nil
}
