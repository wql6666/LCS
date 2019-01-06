package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Result struct {
	r   *http.Response
	err error
}

func process() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	//后边的1啥意思?
	//var c chan Result//为什么不行
	c := make(chan Result, 1) //初始化一个管道用来存储http请求的结果
	//网址少了个s，报错误，空指针？
	req, err := http.NewRequest("Get", "https://www.baidu.com", nil)
	if err != nil {
		fmt.Println("http request failed ,err", err)
		return
	}
	go func() {
		resp, err := client.Do(req)
		pack := Result{r: resp, err: err}
		c <- pack
		//time.Sleep(3*time.Second)
	}()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		res := <-c
		//fmt.Println("Time out")
		fmt.Println("Time out", res.err)
	case res := <-c:
		defer res.r.Body.Close()
		out, _ := ioutil.ReadAll(res.r.Body)
		fmt.Printf("sever response:%s", out)
	}
	return
}

func main() {
	process()
}
