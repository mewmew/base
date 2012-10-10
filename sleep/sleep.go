package main

import "flag"
import "fmt"
import "log"
import "os"
import "strconv"
import "time"

func init() {
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: sleep TIME...")
	fmt.Fprintln(os.Stderr, "Pause execution for the duration TIME.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, "  sleep 3       Sleep for 3 seconds.")
	fmt.Fprintln(os.Stderr, "  sleep 1m3.1s  Sleep for 1 minute and 3.1 seconds.")
}

func main() {
	for _, duration := range flag.Args() {
		err := Sleep(duration)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// Sleep sleeps for the specified duration.
func Sleep(duration string) (err error) {
	d, err := time.ParseDuration(duration)
	if err != nil {
		sec, err := strconv.Atoi(duration)
		if err != nil {
			return err
		}
		d = time.Duration(sec) * time.Second
	}
	time.Sleep(d)
	return nil
}
