package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

func routine(timeout int, ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Second * time.Duration(timeout))
	ch <- fmt.Sprintf("routine - %v", timeout)
}

func actionRoutine() {
	var wg sync.WaitGroup

	var r chan string = make(chan string)
	var r2 chan string = make(chan string)
	var r3 chan string = make(chan string)
	var r1Val string
	var r2Val string
	var r3Val string

	wg.Add(1)
	go routine(5, r, &wg)
	wg.Add(1)
	go routine(3, r2, &wg)
	wg.Add(1)
	go routine(5, r3, &wg)

	go func() {
		wg.Wait()
		close(r)
		close(r2)
		close(r3)
	}()

	r1Val = <-r
	r2Val = <-r2
	r3Val = <-r3

	fmt.Println(r1Val)
	fmt.Println(r2Val)
	fmt.Println(r3Val)
}

func dualChannel(timeout int, ch1 chan string, ch2 chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0

	time.Sleep(time.Second * time.Duration(timeout))

	for i := 0; i < 10; i++ {
		count++
	}

	for i := 0; i < 10; i++ {
		count++
	}

	output := fmt.Sprintf("%v time - dualroutine - %v", count, timeout)

	func() {
		ch1 <- output
		ch2 <- output
	}()
}

func routineDualChannels() {
	var wg sync.WaitGroup

	var ch1 chan string = make(chan string)
	var ch2 chan string = make(chan string)

	var r1Val string
	var r2Val string

	wg.Add(1)
	go dualChannel(5, ch1, ch2, &wg)

	go func() {
		wg.Wait()
		close(ch1)
		close(ch2)
	}()

	r2Val = <-ch2
	r1Val = <-ch1

	fmt.Println(r1Val)
	fmt.Println(r2Val)
}

func insideRoutine(timeout int, data <-chan string, ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	lastRoutineData := <-data

	time.Sleep(time.Second * time.Duration(timeout))
	output := fmt.Sprintf("insideroutine - %v last(%s)", timeout, lastRoutineData)
	ch <- output
	ch <- output
}

func routineInsideRoutine() {
	var wg sync.WaitGroup

	var r chan string = make(chan string)
	var r2 chan string = make(chan string)
	var r3 chan string = make(chan string)
	var r1Val string
	var r2Val string
	var r3Val string

	wg.Add(1)
	go routine(5, r, &wg)
	wg.Add(1)
	go insideRoutine(3, r, r2, &wg)
	wg.Add(1)
	go insideRoutine(3, r2, r3, &wg)

	go func() {
		wg.Wait()
		close(r)
		close(r2)
		close(r3)
	}()

	r2Val = <-r2
	r3Val = <-r3

	fmt.Println(r1Val)
	fmt.Println(r2Val)
	fmt.Println(r3Val)
}

func main() {
	vals := strings.Split("06_01_2020=23:25", "=")
	dates := strings.Split(vals[0], "_")

	fmt.Println(dates[0] + " month - " + dates[1] + " day - " + dates[2] + " year")
	fmt.Println(vals[1])

	// initialActionRoutine := time.Now()
	// actionRoutine()
	// fmt.Println(time.Since(initialActionRoutine))

	// initialDualRoutine := time.Now()
	// routineDualChannels()
	// fmt.Println(time.Since(initialDualRoutine))

	// initialInsideRoutine := time.Now()
	// routineInsideRoutine()
	// fmt.Println(time.Since(initialInsideRoutine))

}
