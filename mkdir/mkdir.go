package main

import "flag"
import "fmt"
import "log"
import "os"
import "strconv"

// When flagParent is true, mkdir creates any necessary parent directories and
// does not complain if the target directory already exists.
var flagParent bool

// flagMode is the permissions to be used when creating a directory.
var flagMode = perm(0777)

type perm os.FileMode

func (p *perm) String() string {
	return fmt.Sprintf("%04o", *p)
}

func (p *perm) Set(v string) (err error) {
	m, err := strconv.ParseInt(v, 8, 32)
	if err != nil {
		return err
	}
	*p = perm(m)
	return nil
}

func init() {
	flag.BoolVar(&flagParent, "p", false, "Make parent directories as needed. No error if existing.")
	flag.Var(&flagMode, "m", "Set permissions of the directory.")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: mkdir [OPTION]... DIRECTORY...")
	fmt.Fprintln(os.Stderr, "Create DIRECTORY(ies).")
	fmt.Fprintln(os.Stderr)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	for _, dir := range flag.Args() {
		err := mkdir(dir)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// mkdir creates the provided directory, including any necessary parent
// directories if the "-p" flag is set. The "-m" flag sets the permissions to be
// used when creating the directory.
func mkdir(dir string) (err error) {
	if flagParent {
		err = os.MkdirAll(dir, os.FileMode(flagMode))
		if err != nil {
			return err
		}
		return nil
	}
	err = os.Mkdir(dir, os.FileMode(flagMode))
	if err != nil {
		return err
	}
	return nil
}
