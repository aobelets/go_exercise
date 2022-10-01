package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const N = 20

func main() {

	fn := func(x int) int {
		time.Sleep(time.Duration(rand.Int31n(N)) * time.Second)
		return x * 2
	}
	in1 := make(chan int, N)
	in2 := make(chan int, N)
	out := make(chan int, N)

	start := time.Now()
	Merge2Channels(fn, in1, in2, out, N+1)
	for i := 0; i < N+1; i++ {
		in1 <- i
		in2 <- i
	}

	orderFail := false
	EvenFail := false
	for i, prev := 0, 0; i < N; i++ {
		c := <-out
		if c%2 != 0 {
			EvenFail = true
		}
		if prev >= c && i != 0 {
			orderFail = true
		}
		prev = c
		fmt.Println(c)
	}
	if orderFail {
		fmt.Println("порядок нарушен")
	}
	if EvenFail {
		fmt.Println("Есть не четные")
	}
	duration := time.Since(start)
	if duration.Seconds() > N {
		fmt.Println("Время превышено")
	}
	fmt.Println("Время выполнения: ", duration)
}

func Merge2Channels(fn func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {

	wg := new(sync.WaitGroup)
	mu := new(sync.Mutex)
	wg.Add(n)
	Array1 := make([]int, n)
	Array2 := make([]int, n)

	go func(wg *sync.WaitGroup, mu *sync.Mutex) {
		for i := 0; i < n; i++ {
			mu.Lock()
			x1 := <-in1
			x2 := <-in2
			mu.Unlock()
			go func(i, x1 int) {
				Array1[i] = fn(x1)
				Array2[i] = fn(x2)
				wg.Done()
			}(i, x1)
		}
	}(wg, mu)

	go func(wg1 *sync.WaitGroup) {
		wg.Wait()
		for i := 0; i < n; i++ {
			out <- Array1[i] + Array2[i]
			//fmt.Println(Array1, Array2)
		}
		fmt.Println(Array1, Array2)
	}(wg)
}

/*
func Merge2Channels(fn func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {

	wg := new(sync.WaitGroup)
	mu := new(sync.Mutex)
	wg.Add(n * 2)
	Array1 := make([]int, n)
	Array2 := make([]int, n)

	go func(wg *sync.WaitGroup, mu *sync.Mutex) {
		for i := 0; i < n; i++ {
			mu.Lock()
			x1 := <-in1
			mu.Unlock()
			Array1[i] = fn(x1)
			wg.Done()
		}
	}(wg, mu)

	go func(wg *sync.WaitGroup, mu *sync.Mutex) {
		for i := 0; i < n; i++ {
			mu.Lock()
			x2 := <-in2
			mu.Unlock()
			Array2[i] = fn(x2)
			wg.Done()
		}
	}(wg, mu)

	go func(wg *sync.WaitGroup, mu *sync.Mutex) {
		wg.Wait()
		for i := 0; i < n; i++ {
			out <- Array1[i] + Array2[i]
			//fmt.Println(Array1, Array2)
		}
		//fmt.Println(Array1, Array2)
	}(wg, mu)
}


*/

/*
func Merge2Channels(fn func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	wg := new(sync.WaitGroup)
	mu := new(sync.Mutex)
	Array1 := make([]int, n)
	Array2 := make([]int, n)
	go func(wg *sync.WaitGroup, mu *sync.Mutex) {
		wg.Add(n)
		for i := 0; i < n; i++ {
			mu.Lock()
			//x1 := <-in1
			Array1[i] = <-in1
			Array2[i] = <-in2
			mu.Unlock()
			//Array1[i] = fn(x1)
			wg.Done()
		}
		//wg.Done()
	}(wg, mu)

	wg.Wait()
	go func(wg *sync.WaitGroup, mu *sync.Mutex) {
		wg.Wait()
		for i := 0; i < n; i++ {
			//wg.Add(1)
			out <- Array1[i] + Array2[i]
			fmt.Println(Array1, Array2)
			//wg.Done()
		}
		fmt.Println(Array1, Array2)
		//wg.Done()
	}(wg, mu)
}
}*/
