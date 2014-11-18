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
	cmds = make(chan func(b *box))
	die  = make(chan bool, participants)
)

func exposeBox() {
	var beat box
	for {
		select {
		case fn := <-cmds:
			clearScreen()
			fn(&beat)
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

func donor(num int) {
	for {
		select {
		case cmds <- func(b *box) {
			b.donations += num
			b.openedBy = donors[rand.Intn(5)]
		}:
		case <-die:
			return
		}
		<-time.After(time.Duration(rand.Intn(2)) * time.Second)
	}
}

func thief(num int) {
	for {
		select {
		case cmds <- func(b *box) {
			if b.donations-num < 0 {
				num = b.donations
			}
			b.donations -= num
			b.openedBy = thieves[rand.Intn(5)]
		}:
		case <-die:
			return
		}
		<-time.After(time.Duration(rand.Intn(2)) * time.Second)
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
