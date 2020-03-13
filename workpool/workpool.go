package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Job struct {
	id     int
	number int
}

type Result struct {
	job    *Job
	result int
}

func calc(job *Job, resultChan chan<- *Result) {
	var sum int
	n := job.number
	for n != 0 {
		temp := n % 10
		sum += temp
		n /= 10
	}

	resultChan <- &Result{
		job:    job,
		result: sum,
	}
}

func work(jobChan chan *Job, resultChan chan<- *Result, ctx context.Context, id int) {
loop:
	for {
		select {
		case job := <-jobChan:
			calc(job, resultChan)
		case <-ctx.Done():
			fmt.Printf("goroutine %d quite\n", id)
			break loop
		}
	}
}

func workPool(workNum int, jobChan chan *Job, resultChan chan<- *Result, ctx context.Context) {
	for i := 0; i < workNum; i++ {
		go work(jobChan, resultChan, ctx, i)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	jobChan := make(chan *Job, 1000)
	resultchan := make(chan *Result, 1000)
	ctx, cancelFunc := context.WithCancel(context.Background())
	go workPool(128, jobChan, resultchan, ctx)

	go func() {
		for result := range resultchan {
			fmt.Printf("n:%d,id:%d,result:%v\n", result.job.number, result.job.id, result.result)
		}
	}()

	index := 0
	after := time.After(time.Second * 10)
quit:
	for {
		select {
		case <-after:
			cancelFunc()
			break quit
		default:
			index++
			n := rand.Int()
			jobChan <- &Job{
				id:     index,
				number: n,
			}
			time.Sleep(time.Millisecond * 100)
		}
	}

	time.Sleep(time.Second)
	fmt.Println("done")
}
