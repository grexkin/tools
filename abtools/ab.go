package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	cocurrent  int
	totalRequest int32
	url string
	wg sync.WaitGroup

)

var (
	totalFaild int32
	totalSucc  int32
	totalNot200 int32
	totalFinished int32
)

func run()  {
	defer wg.Done()
	partNum := totalRequest/10
	for {
		totalFinishedRequest := atomic.LoadInt32(&totalFinished)
		if totalFinishedRequest > totalRequest {
			break
		}
		//进度打印
		if totalFinishedRequest > 0 && totalFinishedRequest % partNum == 0 {
			fmt.Printf("total finished:%d requests\n", totalFinishedRequest)
		}

		resp,err := http.Get(url)
		if err != nil {
			atomic.AddInt32(&totalFaild,1)
			atomic.AddInt32(&totalFinished,1)
			continue
		}
		atomic.AddInt32(&totalFinished, 1)
		if resp.StatusCode != http.StatusOK {
			atomic.AddInt32(&totalNot200,1)
		} else {
			atomic.AddInt32(&totalSucc,1)
		}
	}

}

func main() {
	var tempTotalRequest int

	flag.IntVar(&cocurrent,"c",10,"please input a cocurrent")
	flag.IntVar(&tempTotalRequest,"n",10000,"please input a cocurrent")
	flag.StringVar(&url,"url","http://localhost:8080","please input a ab test url")
	flag.Parse()

	totalRequest = int32(tempTotalRequest)

	startTime := time.Now().UnixNano()   //获取纳秒
	for i:=0;i < cocurrent;i++ {
		wg.Add(1)
		//并发执行
		go run()
	}
	wg.Wait()
	endTime := time.Now().UnixNano()
	costMs := (endTime-startTime)/1000/1000
	if costMs == 0 {
		panic("cost ms is zero")
	}
	requestPerSec := int64(totalRequest) / costMs

	fmt.Printf("total request:%d\n",totalRequest)
	fmt.Printf("total faild:%d\n",totalFaild)
	fmt.Printf("total not 200 request:%d\n",totalNot200)
	fmt.Printf("total success request:%d\n",totalSucc)
	fmt.Printf("request per ms:%d\n",requestPerSec)
}