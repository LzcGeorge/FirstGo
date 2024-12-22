package main

import (
	"fmt"
	"time"
)

type signal struct{}

func worker() {
	println("worker is working...")
	time.Sleep(1 * time.Second)
}

func spawn(f func()) <-chan signal {
	c := make(chan signal)
	go func() {
		println("worker start to work...")
		f()
		c <- signal{}
	}()
	/*
		var s chan signal
		return s
		// 返回一个未初始化的通道，会导致死锁
		fatal error: all goroutines are asleep - deadlock!

		goroutine 1 [chan receive (nil chan)]:
		main.main()
	*/
	return c
}

func main() {
	println("start a worker...")
	c := spawn(worker)
	// 等待并接受来自通道c的信号
	comma, ok := <-c
	fmt.Println(comma, ok)
	fmt.Println("worker work done!")
}
