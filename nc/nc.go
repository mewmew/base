package main

import "flag"
import "fmt"
import "io"
import "log"
import "net"
import "os"

// flagProto specifies the protocol to be used for connections.
var flagProto string

// When flagListen is true, listen for incoming connections.
var flagListen bool

func init() {
	flag.StringVar(&flagProto, "proto", "tcp", "Transfer protocol.")
	flag.BoolVar(&flagListen, "l", false, "Listen for incoming connections.")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: nc [OPTION]... ADDR")
	fmt.Fprintln(os.Stderr, "Read and write data across networks.")
	fmt.Fprintln(os.Stderr)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, "  Connect to example.org on TCP port 8080.")
	fmt.Fprintln(os.Stderr, "    nc example.org:8080")
	fmt.Fprintln(os.Stderr, "  Listen for connections on TCP port 8080.")
	fmt.Fprintln(os.Stderr, "    nc -l :8080")
	fmt.Fprintln(os.Stderr, "  Listen for connections on localhost at TCP port 8080.")
	fmt.Fprintln(os.Stderr, "    nc -l 127.0.0.1:8080")
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	addr := flag.Arg(0)
	if flagListen {
		// listen
		err := listen(addr)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		// connect
		err := connect(addr)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// listen listens for incoming connections and handles their input and output.
func listen(addr string) (err error) {
	l, err := net.Listen(flagProto, addr)
	if err != nil {
		return err
	}
	cl := NewConnList()
	go listenInput(cl)
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		cl.Add(conn.RemoteAddr().String(), conn)
		go listenOutput(conn, cl)
	}
	return nil
}

// listenInput writes to all connected clients from standard input.
func listenInput(cl *connList) {
	buf := make([]byte, 32*1024)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			if err == io.EOF {
				// clean exit on EOF.
				os.Exit(0)
			}
			log.Fatalln(err)
		}
		for _, conn := range cl.Conns() {
			// write input to all connected clients.
			_, err := conn.Write(buf[:n])
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// listenOutput writes to standard output from conn. Client connections are
// removed from the list once they disconnect.
func listenOutput(conn net.Conn, cl *connList) {
	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		log.Println(err)
	}
	// client has disconnected.
	cl.Del(conn.RemoteAddr().String())
}

// connect connects to host:port and handles the connection's input and output.
func connect(addr string) (err error) {
	conn, err := net.Dial(flagProto, addr)
	if err != nil {
		return err
	}
	done := make(chan bool)
	go connectInput(conn, done)
	go connectOutput(conn, done)
	<-done
	return nil
}

// connectInput writes to conn from standard input. Once complete, it sends a
// notification on the done channel.
func connectInput(conn net.Conn, done chan bool) {
	_, err := io.Copy(conn, os.Stdin)
	if err != nil {
		log.Println(err)
	}
	done <- true
}

// connectOutput writes to standard output from conn. Once complete, it sends a
// notification on the done channel.
func connectOutput(conn net.Conn, done chan bool) {
	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		log.Println(err)
	}
	done <- true
}
