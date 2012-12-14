package main

import "net"
import "sync"

// connList is a list of connections that can safely be used concurrently.
type connList struct {
	sync.Mutex
	m map[string]net.Conn
}

// NewConnList returns a new connection list.
func NewConnList() (cl *connList) {
	cl = &connList{
		m: make(map[string]net.Conn),
	}
	return cl
}

// Add adds the specified connection.
func (cl *connList) Add(addr string, conn net.Conn) {
	cl.Lock()
	defer cl.Unlock()
	cl.m[addr] = conn
}

// Del removes the specified connection.
func (cl *connList) Del(addr string) {
	cl.Lock()
	defer cl.Unlock()
	delete(cl.m, addr)
}

// Conns returns all connections.
func (cl *connList) Conns() (conns []net.Conn) {
	cl.Lock()
	defer cl.Unlock()
	for _, conn := range cl.m {
		conns = append(conns, conn)
	}
	return conns
}
