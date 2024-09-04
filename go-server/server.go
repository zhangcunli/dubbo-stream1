package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
	"github.com/zhangcunli/dubbo-stream1/proto/greetmsg"
	"github.com/zhangcunli/dubbo-stream1/proto/greetsvr"
)

func main() {
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(20000),
		),
	)
	if err != nil {
		fmt.Printf(">>>main new server failed, err:%s\n", err)
		panic(err)
	}

	if err := greetsvr.RegisterGreetServiceHandler(srv, &GreetTripleServer{}); err != nil {
		fmt.Printf(">>>main regist server failed, err:%s\n", err)
		panic(err)
	}
	if err := srv.Serve(); err != nil {
		fmt.Printf("Server err:%s\n", err)
	}
}

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) GreetStream(ctx context.Context, stream greetsvr.GreetService_GreetStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("triple BidiStream recv error: %s", err)
		}

		msg := req.Name
		if err := stream.Send(&greetmsg.GreetStreamResponse{Greeting: msg}); err != nil {
			fmt.Printf(">>>GreetStream Send msg:%s\n", msg)
			return fmt.Errorf("triple BidiStream send error: %s", err)
		}
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
func (srv *GreetTripleServer) GreetServerStream(ctx context.Context, req *greetmsg.GreetServerStreamRequest, stream greetsvr.GreetService_GreetServerStreamServer) error {
	key := constant.AttachmentKey
	value, ok := ctx.Value(key).(map[string]interface{})
	if !ok {
		fmt.Printf(">>>GreetServerStream ctx key:%s not exist!!!\n", key)
	} else {
		//values, _ := json.Marshal(value)
		//fmt.Printf(">>>ctx key:%s, value:%s\n", key, values)

		key1 := strings.ToLower("fusion-context")
		fValue, _ := value[key1]
		fmt.Printf(">>>ctx attachement key:%s, value:%s\n", key1, fValue)
	}

	for i := 0; i < 5; i++ {
		message := fmt.Sprintf("Hello, %s!", req.Name)
		fmt.Printf(">>>GreetServerStream send: %s\n", message)
		if err := stream.Send(&greetmsg.GreetServerStreamResponse{Greeting: message}); err != nil {
			return fmt.Errorf("triple ServerStream send err: %s", err)
		}
	}
	return nil
}
