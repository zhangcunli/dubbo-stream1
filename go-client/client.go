package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol/triple/triple_protocol"
	"github.com/zhangcunli/dubbo-stream1/proto/greetmsg"
)

func main() {
	testTripleStream()
}

func testTripleStream() error {
	triAddr := "tri://127.0.0.1:20000"
	cli, err := client.NewClient(
		client.WithClientURL(triAddr),
	)
	if err != nil {
		fmt.Printf(">>>testTripleStream NewClient err:%s\n", err)
		panic(err)
	}

	key := strings.ToLower("fusion-context")
	fusionContext := make(map[string]interface{})
	value, _ := json.Marshal(newFusionContext())
	fusionContext[key] = string(value)
	ctx := context.WithValue(context.Background(), constant.AttachmentKey, fusionContext)

	interfaceName := "greetsvr.GreetService"
	greetService_ClientInfo := client.ClientInfo{
		InterfaceName: interfaceName,
		MethodNames:   []string{"GreetStream", "GreetServerStream"},
		//ConnectionInjectFunc: func(dubboCliRaw interface{}, conn *client.Connection) {
		//	//dubboCli := dubboCliRaw.(*GreetServiceImpl)
		//	dubboCli.conn = conn
		//},
	}

	conn, err := cli.DialWithInfo(interfaceName, &greetService_ClientInfo)
	if err != nil {
		fmt.Printf(">>>testTripleStream DialWithInfo err:%s\n", err)
		return err
	}

	methodName := "GreetServerStream"
	svrStream, err := conn.CallServerStream(ctx, &greetmsg.GreetServerStreamRequest{Name: "Triple.stream.client"}, methodName)
	if err != nil {
		fmt.Printf(">>>testTripleStream CallServerStream err:%s\n", err)
		return err
	}

	fmt.Printf("\n")
	stream := svrStream.(*triple_protocol.ServerStreamForClient)
	msg := new(greetmsg.GreetServerStreamResponse)
	for stream.Receive(msg) {
		fmt.Printf("TRIPLE server stream resp Greeting: %s\n", msg.Greeting)
	}
	if stream.Err() != nil {
		return err
	}
	if err := stream.Close(); err != nil {
		return err
	}
	return nil
}

type FusionContext struct {
	ClientInfo ClientInfo `json:"clientInfo"`
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ClientInfo struct {
	ClientId *string `json:"clientId,omitempty"`
	BizType  *int    `json:"bizType,omitempty"`
	AppId    *int    `json:"appId,omitempty"`
}

type RequestInfo struct {
	Lang          string             `json:"lang"`
	NeutralDomain bool               `json:"neutralDomain"`
	BizData       *map[string]string `json:"bizData,omitempty"`
}

func newFusionContext() *FusionContext {
	clientId := "clientId_0001"
	bizType := 10
	clientInfo := ClientInfo{
		ClientId: &clientId,
		BizType:  &bizType,
	}

	bizData := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	requestInfo := RequestInfo{
		Lang:          "zh",
		NeutralDomain: true,
		BizData:       &bizData,
	}
	fusionContext := FusionContext{
		ClientInfo:  clientInfo,
		RequestInfo: requestInfo,
	}

	return &fusionContext
}
