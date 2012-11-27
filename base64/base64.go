package main

import "encoding/base64"
import "flag"
import "fmt"
import "io"
import "log"
import "os"

// When flagDecode is true, decode the provided input.
var flagDecode bool

func init() {
	flag.BoolVar(&flagDecode, "d", false, "Decode data.")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: base64 [OPTION]... [FILE]")
	fmt.Fprintln(os.Stderr, "Base64 encode or decode FILE, or standard input, to standard output.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "With no FILE, or when FILE is -, read standard input.")
	fmt.Fprintln(os.Stderr)
	flag.PrintDefaults()
}

// StdinFileName is a reserved file name used for standard input.
const StdinFileName = "-"

func main() {
	flag.Parse()
	switch flag.NArg() {
	case 0:
		// Read from stdin when no FILE has been provided.
		err := b64(StdinFileName)
		if err != nil {
			log.Fatalln(err)
		}
	case 1:
		err := b64(flag.Arg(0))
		if err != nil {
			log.Fatalln(err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}

// b64 encodes or decodes the content of a provided file or standard input (when
// the provided file path is "-").
func b64(filePath string) (err error) {
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

	if flagDecode {
		// Decode data.
		err = decode(os.Stdout, fr)
		if err != nil {
			return err
		}
	} else {
		// Encode data.
		err = encode(os.Stdout, fr)
		if err != nil {
			return err
		}
	}

	return nil
}

// decode decodes the provided base64-encoded data from the io.Reader to the
// io.Writer.
func decode(w io.Writer, r io.Reader) (err error) {
	// decode.
	dec := base64.NewDecoder(base64.StdEncoding, r)
	_, err = io.Copy(w, dec)
	if err != nil {
		return err
	}
	return nil
}

// encode encodes the provided data to base64-encoding from the io.Reader to the
// io.Writer.
func encode(w io.Writer, r io.Reader) (err error) {
	// encode.
	enc := base64.NewEncoder(base64.StdEncoding, w)
	// output a newline after the output from enc.Close()
	defer fmt.Fprintln(w)
	defer enc.Close()
	_, err = io.Copy(enc, r)
	if err != nil {
		return err
	}
	return nil
}
