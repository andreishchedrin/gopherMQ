package main

import (
	"andreishchedrin/gopherMQ/server/message"
	"bytes"
	"context"
	"google.golang.org/grpc"
	"net/http"
	"sync"
	"testing"
)

// set correct mode in env before run
func BenchmarkFiberApi(b *testing.B) {
	b.ReportAllocs()
	url := "http://127.0.0.1:8888/push"

	var jsonStr = []byte(`{"channel":"channel1","message":"payload"}`)
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
			if err != nil {
				panic(err)
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				panic("bad request")
			}
		}()
	}
	wg.Wait()
}

func BenchmarkGrpcApi(b *testing.B) {
	b.ReportAllocs()

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := grpc.Dial("127.0.0.1:8888", grpc.WithInsecure())
			if err != nil {
				panic(err)
			}
			defer conn.Close()

			client := message.NewPusherClient(conn)

			resp, err := client.Push(context.Background(), &message.PushStruct{Channel: "channel1", Message: "payload"})
			if err != nil {
				panic(err)
			}

			if resp.Code != 200 {
				panic("bad request")
			}
		}()
	}
	wg.Wait()
}
