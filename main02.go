package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{Timeout: 5 * time.Second}

	// 创建http.Request对象
	req, err := http.NewRequest("GET", "http://rnproxy.zhangjiashu1.tech", nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println(1)
		return
	}

	// 发送请求并获取响应结果
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println(2)
		return
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 打印响应结果
	fmt.Println(string(body))
}
