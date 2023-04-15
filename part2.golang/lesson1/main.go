package main

import (
	"fmt"
	"time"
)

// // 多个lock
// func main() {
// 	go lock()
// 	go rlock()
// 	go wlock()
// 	time.Sleep(5 * time.Second)
// }

// func lock() {
// 	lock := sync.Mutex{}
// 	for i := 0; i < 3; i++ {
// 		lock.Lock()
// 		// defer是在函数运行完成后执行， 因此这里会卡住
// 		defer lock.Unlock()
// 		fmt.Println("I'm sync.Mutex{}:", i)
// 	}
// 	fmt.Println("I'm sync.Mutex{}>> Done")
// }

// func rlock() {
// 	lock := sync.RWMutex{}
// 	for i := 0; i < 3; i++ {
// 		lock.RLock()
// 		// defer是在函数运行完成后执行， 因此这里会卡住
// 		defer lock.RUnlock()
// 		fmt.Println("I'm sync.RWMutex{read}:", i)
// 	}
// 	fmt.Print("I'm sync.Mutex{read}>> Done")
// }

// func wlock() {
// 	lock := sync.RWMutex{}
// 	for i := 0; i < 3; i++ {
// 		lock.Lock()
// 		// defer是在函数运行完成后执行， 因此这里会卡住
// 		defer lock.Unlock()
// 		fmt.Println("I'm sync.RWMutex{write}:", i)
// 	}
// 	fmt.Println("I'm sync.RWMutex{write}>> Done")
// }

// waitgroup 示范
// 为了让主线程可以等待其他线程完成后结束调用
// 1. 通过延迟调用等待 sleep
// 2. 通过等待channel数据完成
// 3. 通过waitgroup 等待
func main() {
	go waitByChan()
}

func waitBySleep() {
	for i := 0; i < 100; i++ {
		fmt.Print(i, " ")
	}

	time.Sleep(time.Second * 5)
}

func waitByChan() {
	ch := make(chan int, 100)

	for i := 0; i < 100; i++ {
		go func(i int) {
			ch <- i
		}(i)
	}

	for j := 0; j < 100; j++ {
		x := <-ch
		fmt.Print(x)
	}
	// defer close(ch)
}
