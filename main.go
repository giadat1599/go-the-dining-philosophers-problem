package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Philosopher is a struct which stores information about a philosopher.
type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

// philosophers is list of all philosophers
var philosophers = []Philosopher{
	{name: "Plato", leftFork: 4, rightFork: 0},
	{name: "Socrates", leftFork: 0, rightFork: 1},
	{name: "Aristotle", leftFork: 1, rightFork: 2},
	{name: "Pascal", leftFork: 2, rightFork: 3},
	{name: "Locke", leftFork: 4, rightFork: 0},
}

var hunger = 3 // how many times does a person eat
var eatTime = 1 * time.Second
var thinkTime = 3 * time.Second
var sleepTime = 1 * time.Second

var orderMutex sync.Mutex
var orderFinished []string

func main() {

	// print out the starting message
	fmt.Println("The table is empty")
	time.Sleep(sleepTime)

	// start the meal
	dine()

	// print out finished message
	fmt.Println("The table is empty")

	time.Sleep(sleepTime)
	fmt.Printf("Order finished: %s.\n", strings.Join(orderFinished, ", "))
}

func dine() {
	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	seatedWg := &sync.WaitGroup{}
	seatedWg.Add(len(philosophers))

	// forks is map of all 5 forks.
	var forks = make(map[int]*sync.Mutex)
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	// start the meal.
	for i := 0; i < len(philosophers); i++ {
		// fire off the go routine for the current philosopher
		go func(philosopher Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seatedWg *sync.WaitGroup) {

			defer wg.Done()

			// seat the philosopher at the table
			fmt.Printf("%s is seated at the table.\n", philosopher.name)
			seatedWg.Done()

			seatedWg.Wait()

			// eat three times
			for i := hunger; i > 0; i-- {
				// get a lock on both forks
				forks[philosopher.leftFork].Lock()
				fmt.Printf("\t%s takes the left fork.\n", philosopher.name)
				forks[philosopher.rightFork].Lock()
				fmt.Printf("\t%s takes the right fork.\n", philosopher.name)

				fmt.Printf("\t%s has both forks and is eating.\n", philosopher.name)
				time.Sleep(eatTime)

				fmt.Printf("\t%s is thinking.\n", philosopher.name)
				time.Sleep(thinkTime)

				forks[philosopher.leftFork].Unlock()
				forks[philosopher.rightFork].Unlock()

				fmt.Printf("\t%s put down the forks.\n", philosopher.name)
			}

			fmt.Println(philosopher.name, "is satisified")
			fmt.Println(philosopher.name, "left the table")

			orderMutex.Lock()
			orderFinished = append(orderFinished, philosopher.name)
			orderMutex.Unlock()
		}(philosophers[i], wg, forks, seatedWg)
	}

	wg.Wait()
}
