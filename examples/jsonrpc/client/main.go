package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/rokumoe/arpc"
	"github.com/rokumoe/arpc/examples/jsonrpc/proto"
	"github.com/rokumoe/arpc/jsonrpc"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9876")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected")
	cc := arpc.NewClient(jsonrpc.WrapClientConn(conn), jsonrpc.GetCodec())
	defer cc.Close()
	greeter := proto.NewGreeterClient(cc)
	wg := sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			log.Printf("invoke at %d", i)
			ctx := context.Background()
			reply, err := greeter.SayHello(ctx, &proto.HelloRequest{
				Name: fmt.Sprintf("req-%d", i),
			})
			if err != nil {
				log.Printf("invoke error: %v", err)
			} else {
				log.Printf("get message %s at %d", reply.Message, i)
			}
		}(i)
	}
	wg.Wait()
}
