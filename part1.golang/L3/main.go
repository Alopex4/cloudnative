package main

import "fmt"

// 给定一个字符串数组
// [“I”,“am”,“stupid”,“and”,“weak”]
// 用 for 循环遍历该数组并修改为
// [“I”,“am”,“smart”,“and”,“strong”]

func main() {
	var myArray = [5]string{"I", "am", "stupid", "and", "weak"}
	for index, value := range myArray {
		if value == "stupid" {
			myArray[index] = "smart"
		} else if value == "weak" {
			myArray[index] = "smart"
		} else {
			continue
		}
	}
	fmt.Printf("%v", myArray)
}
