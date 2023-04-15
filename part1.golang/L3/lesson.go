package main

import (
	"fmt"
)

func main() {

}

// type Books struct {
// 	title  string
// 	author string
// 	price  float32
// }

// func main() {
// 	book1 := Books{
// 		title:  "me",
// 		author: "author",
// 		price:  12.8}
// 	valuePassChangeAuthor("new author 123", book1)
// 	fmt.Println("main func: ", book1.author)
// 	refPassChangeAuthor("new author 123", &book1)
// 	fmt.Println("main func: ", book1.author)

// }

// func valuePassChangeAuthor(newAuthor string, myBook Books) {
// 	myBook.author = newAuthor
// 	fmt.Printf("Inside value pass func: %s\n", myBook.author)
// }

// func refPassChangeAuthor(newAuthor string, myBook *Books) {
// 	myBook.author = newAuthor
// 	fmt.Printf("Inside ref pass func: %s\n", myBook.author)
// }

func sumNumber(num int) int {
	var counter = 0
	for i := 0; i <= num; i++ {
		counter += i
	}
	return counter
}

func looping() {
	for {
		fmt.Println("abc")
	}
}

func forRange() {
	var myStr string = "hello"
	for index, value := range myStr {
		fmt.Printf("%d, %c\n", index, value)
	}
}

func deleteIndex(index int, tSlice []int) []int {
	return append(tSlice[:index], tSlice[index+1:]...)
}

// 安装 IDE 并安装 Go 语言插件
// 编写一个小程序：
// 给定一个字符串数组
// [“I”,“am”,“stupid”,“and”,“weak”]
// 用 for 循环遍历该数组并修改为
// [“I”,“am”,“smart”,“and”,“strong”]

func task1(strArray []string) []string {
	for index, value := range strArray {
		if value == "stupid" {
			strArray[index] = "smart"
		} else if value == "weak" {
			strArray[index] = "strong"
		} else {
			continue
		}
	}
	return strArray
}
