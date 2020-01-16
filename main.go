package main

import (
	"fmt"
	"github.com/pruxa/tets_trotler/throttler"
)

func main() {
	fmt.Println("Example of throttler usage")
	fmt.Println("Task will print self number with '.'")
	t := throttler.NewThrottler()
	fmt.Println("Run empty")
	toRun, err := t.Run(0)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println("toRun = ", toRun)
	fmt.Println("Add 10 tickets and run 5")
	for i := 0; i < 10; i++ {
		t.AddTask(newTask())
	}
	toRun, err = t.Run(5)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println()
	fmt.Println("toRun =", toRun)
	fmt.Println("Add 100 tickets in concurrency and run by 10 until empty the queue")
	for i := 0; i < 100; i++ {
		go t.AddTask(newTask())
	}
	for t.QueueLen() != 0 {
		toRun, err = t.Run(10)
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Println()
		fmt.Println("ToRun = ", toRun)
	}
}

type task struct {
	id int
}

func (t task) Run() {
	fmt.Print(t.id, ".")
}

func newTask() task {
	c++
	return task{c - 1}
}

var c = 1
