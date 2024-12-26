package workerpool

import (
	"fmt"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	p := NewPool(5)

	for i := 0; i < 10; i++ {
		err := p.Schedule(func() {
			time.Sleep(time.Second * 3)
		})
		if err != nil {
			println("task: ", i, "err: ", err)
		}
	}

	p.Free()
}

func TestA(t *testing.T) {
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			fmt.Println("j: ", &j, " val: ", i)
		}()
	}
}
