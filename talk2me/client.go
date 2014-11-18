package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var netAddr = flag.String("addr", ":1234", "Server address & port to connect to.")

// Message structure
type message struct {
	Body   string `xml:"body"`
	Sender string `xml:"nickname"`
}

// Reads a message from given I/O
func readMessage(buf *bufio.Reader) string {
	input, err := buf.ReadString('\n')
	if err != nil {
		log.Println("Error scanning string.")
	}

	return strings.Trim(input, "\r\n")
}

// Listens for incoming message on network connection
// and outputs them to standard I/O
func readIncoming(conn net.Conn) {
	inc := bufio.NewScanner(conn)

	for inc.Scan() {
		msg := new(message)
		xml.Unmarshal(inc.Bytes(), msg)
		fmt.Printf("%s says: %s\n", msg.Sender, msg.Body)
	}
}

func main() {
	flag.Parse()

	conn, err := net.Dial("tcp", *netAddr)
	if err != nil {
		log.Fatalf("Error dialing in (%s)", err)
	}

	userInput := bufio.NewReader(os.Stdin)
	enc := xml.NewEncoder(conn)

	fmt.Print("Enter your name: ")
	nick := readMessage(userInput)
	go readIncoming(conn)

	for {
		input := readMessage(userInput)
		if err := enc.Encode(&message{input, nick}); err != nil {
			log.Printf("Error marshalling (%s)", err)
		}

		fmt.Fprint(conn, "\n")
	}
}
