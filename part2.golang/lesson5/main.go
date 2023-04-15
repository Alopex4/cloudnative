package main

// 1.接收客户端 request，并将 request 中带的 header 写入 response header
// 2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
// 3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
// 4.当访问 localhost/healthz 时，应返回200

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/version", version)
	http.HandleFunc("/client", client)
	http.HandleFunc("/healthz", healthz)
	err := http.ListenAndServe(":8080", nil)
	// 80端口打开需要root权限
	// err := http.ListenAndServe(":80", nil)

	if err != nil {
		println(err)
	}
}

// 1.接收客户端 request，并将 request 中带的 header 写入 response header
func root(w http.ResponseWriter, r *http.Request) {
	for key, valSlice := range r.Header {
		val, err := json.Marshal(valSlice)
		if err != nil {
			println(err)
		}
		w.Header().Add(key, string(val))
	}
}

// 2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
func version(w http.ResponseWriter, r *http.Request) {
	var ver = os.Getenv("version")
	if ver != "" {
		w.Header().Add("version", ver)
	} else {
		w.Header().Add("version", "Not found")
	}
}

// 3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
func client(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Printf("user-agent: %s; IP: %s; statusCode: %d\n", r.UserAgent(), r.Host, http.StatusOK)
}

// 4.当访问 localhost/healthz 时，应返回200
func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
