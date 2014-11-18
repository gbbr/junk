package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"time"
)

var (
	thieves      = [...]string{"Mordy", "Jerk", "Jameson", "Janice", "Picasso"}
	donors       = [...]string{"Jim", "Willy", "Babadook", "Hendrix", "Johnson"}
	participants = len(thieves) + len(donors)
)

type box struct {
	donations int
	openedBy  string
	waiting   bool
}

var (
	act = make(chan func(b *box))
	die = make(chan bool)
)

// exposeBox exposes a donation box for contestants
// to act on.
func exposeBox() {
	var beat box
	for {
		select {
		case fn := <-act:
			fn(&beat)
			clearScreen()
			fmt.Printf("Contents: %+v\n", beat)
			beat.waiting = false
		default:
			if !beat.waiting {
				log.Println("Waiting...")
				beat.waiting = true
			}
		}
	}
}

// donor creates a new donor that is capable of adding
// num items at a time.
func donor(num int) {
	for {
		select {
		case act <- func(b *box) {
			b.openedBy = donors[rand.Intn(5)]
			b.donations += num
		}:
		case <-die:
			return
		}

		<-time.After(time.Duration(rand.Intn(2)) * time.Second)
	}
}

// thief creates a new thief that is capable of taking
// num items in one go.
func thief(num int) {
	for {
		select {
		case act <- steal(num):
		case <-die:
			return
		}

		<-time.After(time.Duration(rand.Intn(2)) * time.Second)
	}
}

// steal removes N items from the charity box
// in the name of a random thief
func steal(N int) func(*box) {
	return func(b *box) {
		b.openedBy = thieves[rand.Intn(5)]
		b.donations -= N
		if b.donations < 0 {
			b.donations = 0
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	go exposeBox()
	for i := 0; i < len(donors); i++ {
		go donor(i)
	}
	for i := 0; i < len(thieves); i++ {
		go thief(i)
	}

	ntr := make(chan os.Signal, 1)
	signal.Notify(ntr, os.Interrupt)
	<-ntr
	// Alert participants
	for i := 0; i < participants; i++ {
		die <- true
	}
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
