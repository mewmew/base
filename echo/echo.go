package main

import "flag"
import "fmt"
import "os"

// suppressNewline specifies whether the trailing newline should be suppressed.
var suppressNewline bool

func init() {
	flag.BoolVar(&suppressNewline, "n", false, "suppress the trailing newline.")
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: echo [OPTION]... [STRING]...")
	fmt.Fprintln(os.Stderr, "Echo STRING(s) to standard output.")
	fmt.Fprintln(os.Stderr)
	flag.PrintDefaults()
}

func main() {
	echo(flag.Args())
}

// echo prints the provided arguments, separated by spaces and terminated by a
// newline, to the standard output.
func echo(args []string) {
	for i, arg := range args {
		fmt.Print(arg)
		if i < len(args)-1 {
			fmt.Print(" ")
		}
	}
	if !suppressNewline {
		fmt.Println()
	}
}
