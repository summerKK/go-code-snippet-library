package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

func main() {
	for i := 0; i < 10240; i++ {
		handle()
	}
}

func handle() {
	conn, err := net.Dial("tcp", ":2020")
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})
	go func() {
		go func() {
			time.Sleep(10)
			done <- struct{}{}
		}()
		_, _ = io.Copy(os.Stdout, conn)
		log.Println("done")
		done <- struct{}{}
	}()

	if _, err = io.Copy(conn, bytes.NewReader([]byte(fmt.Sprintf("hello,world %d", rand.Intn(100))))); err != nil {
		panic(err)
	}
	_ = conn.Close()

	<-done
}
