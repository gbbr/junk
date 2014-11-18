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

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	type box struct {
		items       int
		lastDeposit string
		waiting     bool
	}

	cmds := make(chan func(b *box))
	go func() {
		var beat box
		for {
			select {
			case fn := <-cmds:
				clear()
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
	}()

	// Charity workers
	die := make(chan bool, 9)
	for i := 0; i < 5; i++ {
		name := [5]string{"Jim", "Willy", "Babadook", "Hendrix", "Johnson"}
		go func(num int) {
		DEPOSITS:
			for {
				select {
				case cmds <- func(b *box) {
					b.items += num
					b.lastDeposit = name[rand.Intn(5)]
				}:
				case <-die:
					break DEPOSITS
				}
				<-time.After(time.Duration(rand.Intn(2)) * time.Second)
			}
			fmt.Println("I'm DEAD")
		}(i)
	}

	// Thieves
	for i := 0; i < 5; i++ {
		name := [5]string{"Mordy", "Jerk", "Jameson", "Janice", "Picasso"}
		go func(num int) {
		DEPOSITS:
			for {
				select {
				case cmds <- func(b *box) {
					if b.items-num < 0 {
						num = b.items
					}
					b.items -= num
					b.lastDeposit = name[rand.Intn(5)]
				}:
				case <-die:
					break DEPOSITS
				}
				<-time.After(time.Duration(rand.Intn(2)) * time.Second)
			}
			fmt.Println("I'm DEAD")
		}(i)
	}

	ntr := make(chan os.Signal, 1)
	signal.Notify(ntr, os.Interrupt)
	<-ntr

	for i := 0; i < 10; i++ {
		die <- true
	}
}
