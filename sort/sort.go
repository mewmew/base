package main

import "flag"
import "fmt"
import "log"
import "os"
import "sort"

import "github.com/mewkiz/pkg/bufioutil"

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: sort [OPTION]... [FILE]...")
	fmt.Fprintln(os.Stderr, "Write sorted concatenation of all FILE(s) to standard output.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "With no FILE, or when FILE is -, read standard input.")
	fmt.Fprintln(os.Stderr)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	err := sortFiles(flag.Args())
	if err != nil {
		log.Fatalln(err)
	}
}

// sortFiles writes the sorted concatenation of all provided files or standard
// input (when no file has been provided).
func sortFiles(filePaths []string) (err error) {
	lines := make([]string, 0)
	for _, filePath := range filePaths {
		l, err := bufioutil.ReadLines(filePath)
		if err != nil {
			return err
		}
		lines = append(lines, l...)
	}
	if len(filePaths) == 0 {
		br := bufioutil.NewReader(os.Stdin)
		lines, err = br.ReadLines()
		if err != nil {
			return err
		}
	}
	sort.Strings(lines)
	for _, line := range lines {
		fmt.Println(line)
	}
	return nil
}
