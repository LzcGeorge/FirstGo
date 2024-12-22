package main

import "net/http"

func main() {
	// 创建一个 http server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World"))
	})
	// 启动 http server
	http.ListenAndServe(":8080", nil)
}
