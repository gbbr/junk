package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net"
)

var netAddr = flag.String("addr", ":1234", "Address to listen on (ie. \"localhost:8080\" or \":1234\")")

// Incoming message format
type message struct {
	Message string `xml:"body"`
	From    string `xml:"nickname"`
	origin  net.Conn
}

// Server structure, holds broadcast channel,
// leaver channel, TCP listener and connection pool
type Server struct {
	conns map[net.Conn]bool // Connected clients
	ln    net.Listener      // TCP listener
	bcast chan message      // Message broadcasting channel
	leave chan net.Conn     // Notify disconnect
	join  chan net.Conn     // Notify join
}

// Returns a new server instance that is listening
// on the given address (<host:port>)
func NewServer(addr string) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return &Server{}, err
	}

	log.Printf("Server running on %s\n", addr)

	return &Server{
		make(map[net.Conn]bool),
		ln,
		make(chan message),
		make(chan net.Conn),
		make(chan net.Conn),
	}, nil
}

// Accepts incoming connections and reads input
func (s *Server) Start() {
	log.Println("Waiting for connections...")
	go s.manageTraffic()

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Printf("Error accepting incoming connection (%s)\n", err)
		}

		s.join <- conn
		go s.Handle(conn)
	}
}

// Handles connection and listens for message from
// connected client
func (s *Server) Handle(c net.Conn) {
	scn := bufio.NewScanner(c)
	msg := &message{origin: c}

	for scn.Scan() {
		xml.Unmarshal(scn.Bytes(), msg)
		s.bcast <- *msg
	}

	s.leave <- c
}

// Manages the server connections and broadcasts messages
func (s *Server) manageTraffic() {
	// Private to protect from potential race conditions
	announce := func(msg message) {
		mm, err := xml.Marshal(&msg)
		if err != nil {
			log.Printf("Error marshalling message (%s)\n", err)
		}

		for c := range s.conns {
			if msg.origin != c {
				fmt.Fprintf(c, "%s\n", mm)
			}
		}
	}

	for {
		select {
		// Broadcast
		case m := <-s.bcast:
			announce(m)

		// Someone leaves
		case l := <-s.leave:
			delete(s.conns, l)
			announce(message{"You hear the sound of a door closing.", "Server", l})

		// Somone joins
		case j := <-s.join:
			s.conns[j] = true
			announce(message{"You hear the sound of a door opening.", "Server", j})
		}
	}
}

// Closes the server listener
func (s *Server) Close() error {
	return s.ln.Close()
}

func main() {
	flag.Parse()

	server, err := NewServer(*netAddr)
	if err != nil {
		log.Fatalf("Error opening socket (%s)\n", err)
	}

	server.Start()
	server.Close()
}
