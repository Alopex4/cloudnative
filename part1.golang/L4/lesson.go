package main

import (
	"encoding/json"
	"fmt"
)

// // 传入参数
// func main() {
// 	// 非命名参数 ./main.go hello world
// 	args := os.Args
// 	fmt.Println(args)
// 	// 命名参数 ./main.go --hello world
// 	name := flag.String("hello", "<name>", "say hello to")
// 	flag.Parse()
// 	fmt.Printf("hello %s\n", *name)
// }

// // 多个返回值
// func main() {
// 	ret, str := myFun(10, 20)
// 	fmt.Println(ret, str)
// }

// func myFun(a, b int) (int, string) {
// 	return a + b, "ok"
// }

// // 命名函数返回
// func main() {
// 	ret, str := myFun(10, 20)
// 	fmt.Println(ret, str)
// }

// // var retVal = int 0
// // var status = string ""
// func myFun(a, b int) (retVal int, status string) {
// 	retVal = a + b
// 	status = "ok"
// 	return
// }

// // 返回值选择
// func main() {
// 	// 错误的写法
// 	// all = myFun(10, 20)

// 	// 两个变量都接收
// 	ret, str := myFun(10, 20)

// 	// 忽略 返回的str
// 	ret, _ = myFun(10, 20)
// 	fmt.Println(ret, str)
// }

// // var retVal = int 0
// // var status = string ""
// func myFun(a, b int) (retVal int, status string) {
// 	retVal = a + b
// 	status = "ok"
// 	return
// }

// // 传递多个参数
// func main() {
// 	fmt.Print(calSum(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10))
// }

// func calSum(init int, item ...int) int {
// 	for _, i := range item {
// 		init += i
// 	}
// 	return init
// }

// 内置函数 new vs make
// func main() {
// 	aa := make(map[string]int, 2)
// 	aa["a"] = 10
// 	aa["b"] = 20

// 	bb := new(map[string]int)
// 	// 全初始化
// 	*bb = map[string]int{"a": 10, "b": 20}
// 	// 局部修改
// 	(*bb)["a"] = 100

// 	// var aa = map[string]int{
// 	// 	"a": 10,
// 	// 	"b": 20,
// 	// }

// 	// bb := new(map[string]int){
// 	// 	"a": 100,
// 	// 	"b": 200
// 	// }

// }

// // 回调函数
// func main() {
// 	cal(10, 20, add)
// 	cal(10, 20, minus)
// }

// func cal(a, b int, f func(int, int) int) int {
// 	return f(a, b)
// }

// func add(a, b int) int {
// 	return a + b
// }

// func minus(a, b int) int {
// 	return a - b
// }

// // 匿名函数
// func main() {
// 	who := "world"

// 	// 调用
// 	func() {
// 		println("hello ", who)
// 	}()

// 	// 赋值给变量
// 	myFun := func() {
// 		print("hello ", who)
// 	}
// 	println(myFun)

// 	// 直接调用
// 	func(x, y int) {
// 		println(x + y)
// 	}(10, 20)

// 	// 函数作为返回值调用
// 	println(square(10)())

// }

// // 作函数返回值
// func square(x int) func() int {
// 	return func() int {
// 		return x * x
// 	}
// }

// // 方法
// type circuit struct {
// 	cid string
// }

// func main() {
// 	c1 := circuit{cid: "LDN/CT"}
// 	println(c1.getCid())
// }

// func (c circuit) getCid() string {
// 	return c.cid
// }

// // 值和指针传递
// func main() {
// 	var a int = 0
// 	// 复制一个a的副本值
// 	addOneCant(a)
// 	println(a)

// 	// 创建一个a'指针指向a的地址的内容 (和a的值相同)
// 	addOne(&a)
// 	println(a)

// }

// // 值传递
// func addOneCant(a int) {
// 	a += 1
// }

// // 指针传递
// func addOne(a *int) {
// 	*a += 1
// }

// // 接口
// type Phone interface {
// 	call() string
// }

// type iPhone struct {
// 	model string
// }

// type Nokia struct {
// 	model string
// }

// func (iphone iPhone) call() string {
// 	return iphone.model
// }

// func (nokia Nokia) call() string {
// 	return nokia.model
// }

// func main() {
// 	phones := []Phone{}
// 	myIphone := iPhone{model: "14"}
// 	myNokia := Nokia{model: "a413"}

// 	phones = append(phones, myIphone)
// 	phones = append(phones, myNokia)

// 	for _, p := range phones {
// 		println(p.call())
// 	}
// }

// // reflect 反射器
// type cat struct {
// 	name string
// 	age  int
// }

// func main() {
// 	myCat := cat{name: "leslie", age: 20}
// 	// println(reflect.ValueOf(myCat.age))
// 	fmt.Printf("%v\n", reflect.ValueOf(myCat)) // {leslie 20}
// 	fmt.Printf("%v\n", reflect.TypeOf(myCat))  // main.cat

// 	myMap := make(map[string]int, 1)
// 	myMap["a"] = 100
// 	fmt.Printf("%v\n", reflect.ValueOf(myMap)) // map[a:100]
// 	fmt.Printf("%v\n", reflect.TypeOf(myMap))  // map[string]int
// }

// json编码
// 定义Huamn, 使用大写使其可以公开访问, 否则marshal2JsonString为空
// https://stackoverflow.com/a/8271160
type Human struct {
	Name string
	Age  int
}

func main() {
	// 实例化一个human
	singer := Human{Name: "leslie", Age: 30}
	singerStr := marshal2JsonString(singer)
	fmt.Printf("%v\n", singerStr)

	// 定义一个名称为obj的map对象
	var obj interface{}
	err := json.Unmarshal([]byte(singerStr), &obj)
	if err != nil {
		print(err)
	}

	objMap, _ := obj.(map[string]interface{})
	// fmt.Printf("%v\n", objMap) // map[Age:30 Name:leslie]
	for k, v := range objMap {
		switch value := v.(type) { // 这里的v.(type)为什么可以这样调用
		case string:
			fmt.Printf("type of %s is string, value is %v\n", k, value)
		case interface{}:
			fmt.Printf("type of %s is interface{}, value is %v\n", k, value)
		default:
			fmt.Printf("type of %s is wrong, value is %v\n", k, value)
		}
	}
}

// 将string 变成 struct
func unmarshal2Struct(humanStr string) Human {
	h := Human{}
	err := json.Unmarshal([]byte(humanStr), &h)
	if err != nil {
		print(err)
	}
	return h
}

// 将struct转换为string
func marshal2JsonString(h Human) string {
	myByte, err := json.Marshal(&h)
	if err != nil {
		print(err)
	}
	return string(myByte)
}
