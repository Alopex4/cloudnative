package main

import (
	"fmt"
	"time"
)

// //  创建error接口
// type error interface {
// 	Error() string
// }

// // 错误处理
// func main() {
// 	// 通过errors.new
// 	newError := errors.New("an error")
// 	fmt.Println(newError.Error())

// 	// 通过fmt.Errorf
// 	newError2 := fmt.Errorf("another error")
// 	fmt.Println(newError2.Error())

// 	// 判断错误
// 	_, err := divTen(0)
// 	if err != nil {
// 		println(err.Error())
// 	}
// }

// // 判断是否存在错误
// func divTen(dividend int) (int, error) {
// 	if dividend == 0 {
// 		return 0, errors.New(">> dividend could not be 0")
// 	} else {
// 		return 10 /dividend , nil
// 	}
// }

// // defer 延迟函数
// func main() {
// 	// main closed
// 	// defer no3
// 	// defer no2
// 	// defer no1
// 	println("main closed")
// 	defer println("defer no1")
// 	defer println("defer no2")
// 	defer println("defer no3")
// }

// // defer 延迟函数
// func main() {
// 	println(c()) // 2
// }

// func c() (i int) {
// 	defer func() { i++ }()
// 	return 1
// }

// //defer / panic / recover
// func main() {

// 	defer func() {
// 		fmt.Println("inside defer")
// 		if err := recover(); err != nil {
// 			fmt.Println("Wrong!", err)
// 		}
// 	}()
// 	panic("This is a panic")
// 	println("could not execute") // 无法执行
// }

//// defer / panic / recover 小例子
// func main() {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			fmt.Println(err)
// 		}
// 		fmt.Println("error is ok")
// 		fmt.Printf("keep going...")
// 	}()
// 	panic("error")
// }

// // 一个goroutine
// func main() {
// 	go func() {
// 		for i := 0; i <= 5; i++ {
// 			go println(i)
// 		}
// 	}()
// 	time.Sleep(time.Second)
// }

// // chan 的例子
// func main() {
// 	myChain := make(chan int, 3)
// 	go func() {
// 		// go println("inside goroutine") // 如果开启， 直接返回 get channel： 1 结束
// 		println("inside goroutine")
// 		// 将1写入myChain
// 		myChain <- 1
// 	}()
// 	value := <-myChain
// 	// time.Sleep(time.Second)
// 	println("get channel:", value)
// }

// // 遍历
// func main() {
// 	ch := make(chan int, 10)
// 	go func() {
// 		for i := 0; i <= 10; i++ {
// 			rand.Seed(time.Now().UnixNano())
// 			n := rand.Intn(10)
// 			fmt.Println("putting", n)
// 			ch <- i
// 		}
// 		// 关闭channel
// 		defer close(ch)
// 	}()

// 	fmt.Println("hello from main")
// 	for v := range ch {
// 		fmt.Println("receive", v)
// 	}
// }

// // 关闭通道
// func main() {
// 	ch := make(chan int)
// 	close(ch)

// 	if v, notClosed := <-ch; notClosed {
// 		print(v)
// 	}

// }

// // 多个通道
// func main() {
// 	ch1 := make(chan int)
// 	ch2 := make(chan int)
// 	ch1 <- 1 // fatal error: all goroutines are asleep - deadlock!

// 	select {
// 	case v := <-ch1:
// 		println("ch1", v)
// 	case v := <-ch2:
// 		println("ch2", v)
// 	default:
// 		println("noting")
// 	}
// }

// // 单线程处理
// var count int = 0

// func PrimeNum(n int) {
// 	for i := 2; i < n; i++ {
// 		if n%i == 0 {
// 			return
// 		}
// 	}
// 	count += 1
// 	fmt.Printf("%v\n\t", n)
// }

// func main() {
// 	for i := 2; i < 100001; i++ {
// 		PrimeNum(i)
// 	}

// 	fmt.Printf("Done %d\n", count)
// }

// // 多线程处理 (同步锁)
// var (
// 	count int
// 	lock  sync.Mutex
// )

// func PrimeNum(n int) {
// 	for i := 2; i < n; i++ {
// 		if n%i == 0 {
// 			return
// 		}
// 	}
// 	lock.Lock()
// 	count += 1
// 	lock.Unlock()
// 	fmt.Printf("%v\t", n)
// }

// func main() {
// 	for i := 2; i < 100001; i++ {
// 		go PrimeNum(i)
// 	}

// 	var key string
// 	fmt.Scanln(&key)
// 	fmt.Printf("Done %d\n", count)
// }

// // 多线程处理 (同步锁)
// // 关闭信道  close(chan)
// // for ... range 不断接收值， 直到关闭（缺少关闭机制会出错）
// // 取值方式 value， (ok) <- channel
// // select ... default 适用于无法确定汗是关闭信道的请卡
// var count int
// var ch chan int = make(chan int, 100)

// func PrimeNum(n int, ch chan int) {
// 	for i := 2; i < n; i++ {
// 		if n%i == 0 {
// 			return
// 		}
// 	}
// 	ch <- n
// }

// func main() {
// 	for i := 2; i < 100001; i++ {
// 		go PrimeNum(i, ch)
// 	}

// printing:
// 	for {
// 		select {
// 		case v := <-ch:
// 			fmt.Printf("%d\t", v)
// 			count += 1
// 		default:
// 			fmt.Printf("done %d\n", count)
// 			// 跳出printing的label
// 			break printing
// 		}
// 	}
// }

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
