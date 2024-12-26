package workerpool

import (
	"errors"
	"fmt"
	"sync"
)

type Task func()

type Pool struct {
	capacity int
	active   chan struct{} // active channel: 作为 worker 的计数信号量
	tasks    chan Task     // tasks channel: 任务队列

	wg   sync.WaitGroup // 等待所有 worker 完成
	quit chan struct{}  // 退出信号, 用一个空结构体（节省内存）
}

const (
	defaultCapacity = 100
	maxCapacity     = 10000
)

func NewPool(capacity int) *Pool {
	if capacity <= 0 {
		capacity = defaultCapacity
	}
	if capacity > maxCapacity {
		capacity = maxCapacity
	}

	p := &Pool{
		capacity: capacity,
		active:   make(chan struct{}, capacity),
		tasks:    make(chan Task),
		quit:     make(chan struct{}),
	}

	fmt.Println("workerpool start")
	go p.run()

	return p
}

func (p *Pool) run() {
	idx := 0

	for {
		select {
		case <-p.quit:
			return
		/*
			struct{} 和 struct{}{} 的区别：
			1. struct{} 是类型
			2. struct{}{} 是该类型的实例
		*/
		case p.active <- struct{}{}:
			idx++
			// worker 的编号从 1 开始
			p.newWorker(idx)
		}
	}
}

func (p *Pool) newWorker(idx int) {
	p.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("worker[%03d]: recover from panic: %v and exit\n", idx, err)
				<-p.active
			}
			p.wg.Done()
		}()

		fmt.Printf("worker[%03d]: start\n", idx)

		for {
			select {
			case <-p.quit:
				fmt.Printf("worker[%03d]: exit\n", idx)
				<-p.active
				return
			case task := <-p.tasks:
				fmt.Printf("worker[%03d]: received a task\n", idx)
				task()
			}
		}
	}()
}

var ErrWorkerPoolFreed = errors.New("workerpool freed")

func (p *Pool) Schedule(task Task) error {
	select {
	case <-p.quit:
		return ErrWorkerPoolFreed
	case p.tasks <- task:
		return nil
	}
}

func (p *Pool) Free() {
	close(p.quit)
	p.wg.Wait()
	fmt.Println("workerpool freed")
}
