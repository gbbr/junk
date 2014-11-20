package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"time"
)

var (
	thieves = [...]string{"Mordy", "Jerk", "Jameson", "Janice", "Picasso"}
	donors  = [...]string{"Jim", "Willy", "Babadook", "Hendrix", "Johnson"}
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Expose a box to interact with
	go exposeBox()
	for i := 0; i < len(donors); i++ {
		// Start all the donors
		go donor(i)
	}
	for i := 0; i < len(thieves); i++ {
		// Start an equal number of thieves
		go thief(i)
	}

	ntr := make(chan os.Signal, 1)
	signal.Notify(ntr, os.Interrupt)
	<-ntr
	close(kill)
}

type box struct {
	donations int
	openedBy  string
}

var act = make(chan func(b *box))

// exposeBox exposes a donation box for contestants to act on.
func exposeBox() {
	var b box
	for {
		select {
		case fn := <-act:
			fn(&b)
			clearScreen()
			fmt.Printf("Contents: %+v\n", b)
		case <-kill:
			return
		}
	}
}

var kill = make(chan bool)

// donor creates a new donor that is capable of adding
// num items at a time.
func donor(num int) {
	for {
		select {
		case act <- func(b *box) {
			b.openedBy = donors[rand.Intn(len(donors))]
			b.donations += num
		}:
		case <-kill:
			return
		default:
		}
		time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
	}
}

// thief creates a new thief that is capable of taking
// num items in one go.
func thief(num int) {
	for {
		select {
		case act <- steal(num):
		case <-kill:
			return
		default:
		}
		time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
	}
}

// steal removes N items from the charity box
// in the name of a random thief
func steal(N int) func(*box) {
	return func(b *box) {
		b.openedBy = thieves[rand.Intn(len(thieves))]
		b.donations -= N
		if b.donations < 0 {
			b.donations = 0
		}
	}
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
