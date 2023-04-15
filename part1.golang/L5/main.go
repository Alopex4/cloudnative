package main

import (
	"fmt"
	"time"
)

// 队列长度10，队列元素类型为 int
// • 生产者：
// 每1秒往队列中放入一个类型为 int 的元素，队列满时生产者可以阻塞
// • 消费者：
// 每一秒从队列中获取一个元素并打印，队列为空时消费者阻塞

// customer vs producer
// 只写数据
func producer(ch chan<- int) {
	for {
		if cap(ch) == len(ch) {
			fmt.Println("wait customer")
			time.Sleep(time.Second)
			continue
		} else {
			// fmt.Println("chan length:", len(ch))
			time.Sleep(time.Second)
			ch <- 1
		}
	}
}

// 只读数据
func customer(ch <-chan int, done <-chan bool) {
	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		select {
		case v := <-ch:
			fmt.Println(v)
		case done := <-done:
			fmt.Print("It is over!", done, "\n")
		default:
			time.Sleep(time.Second)
			fmt.Println("wait producer")
		}
	}
}

func main() {
	var ch chan int = make(chan int, 10)
	var done chan bool = make(chan bool)
	go producer(ch)
	go customer(ch, done)

	var key string
	fmt.Scanln(&key)
	done <- true

	fmt.Print("main done\n")
	defer close(ch)
	defer close(done)
}
