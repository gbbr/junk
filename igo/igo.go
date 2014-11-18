package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// dir holds the name of the temporary directory used
var tmpDir string

const (
	welcomeMessage = `
Welcome to Interactive Go playground (alpha)!

The whole standard library is at your service. The rules are
simple: you are inside the main function and multiline commands
are not supported. Use ";" instead.

Type 'q' to quit and 'r' to reset the buffer.
`
	header = `
package main
func main() {
`
	footer = `}`
)

type state struct {
	// current line and command indentation
	line, indent int
	// buf holds the collected instructions of a succesfully compiling program
	buf bytes.Buffer
}

// checkError attempts to compile the current buffer with the provided line and returns
// the message returned by the compiler.
func (s *state) parseLine(line string) (parsedLine, outMsg, errMsg string) {
TRY_COMPILE:
	newLine := s.buf.String() + line + "\r\n"
	tmpFile := filepath.Join(tmpDir, "check.go")
	err := ioutil.WriteFile(tmpFile, []byte(header+newLine+footer), 0644)
	if err != nil {
		log.Fatalf("error writing tempfile: %s", err)
	}

	exec.Command("goimports", "-w", tmpFile).Run()
	cmd := exec.Command("go", "run", tmpFile)
	var stdout, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err = cmd.Run()
	if _, ok := err.(*exec.ExitError); err != nil && !ok {
		log.Fatalf("error executing file: %s", err)
	}

	scn := bufio.NewScanner(&stderr)
	for scn.Scan() {
		ln := scn.Text()
		if strings.HasPrefix(ln, "#") {
			continue
		}
		parts := strings.SplitN(ln, ": ", 2)
		if len(parts) != 2 {
			return line, ln + "\r\n", ""
		}
		if strings.HasSuffix(parts[1], " declared and not used") {
			var miss string
			fmt.Sscanf(parts[1], "%s declared and not used", &miss)
			line += fmt.Sprintf("; _ = %s", miss)
			goto TRY_COMPILE
		}
		return "", stdout.String(), parts[1]
	}
	return line, stdout.String(), ""
}

// accepted processes the passed command at the bottom of the current buffer and returns
// true if it has been succesfully compiled and accepted.
func (s *state) accepted(line string) bool {
	nline, stdout, errmsg := s.parseLine(line)
	if errmsg != "" {
		fmt.Printf("=> syntax error: %s\r\n", errmsg)
		return false
	}
	if stdout != "" {
		fmt.Printf("=> %s\r\n", strings.TrimRight(stdout, "\r\n"))
		return true
	}

	_, err := s.buf.WriteString(nline + "\r\n")
	if err != nil {
		log.Fatalf("could not write to buffer: %+v\r\n", err)
	}

	return true
}

func main() {
	var err error
	tmpDir, err = ioutil.TempDir("", "igo-")
	if err != nil {
		log.Fatalf("could not create tempdir: %s", err)
	}
	st := state{line: 1}
	scn := bufio.NewScanner(os.Stdin)
	fmt.Println(welcomeMessage)
LOOP:
	for {
		fmt.Printf("igo(main):%03d:%d> ", st.line, st.indent)
		if !scn.Scan() {
			break
		}
		switch scn.Text() {
		case "q":
			break LOOP
		case "help":
			fmt.Println(welcomeMessage)
			continue
		case "r":
			st.buf.Reset()
			st.line = 1
			log.Println("Flushed buffer.")
			continue
		case "buf":
			fmt.Println(st.buf.String())
			continue
		default:
			if st.accepted(scn.Text()) {
				st.line++
			}
		}
	}

	os.RemoveAll(tmpDir)
}
