package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"exp/kitex_gen/item"

	"github.com/bytedance/sonic"

	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/pkg/remote/codec"
	"github.com/cloudwego/kitex/pkg/remote/codec/thrift"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/serviceinfo"
)

func mockHandler(ctx context.Context, handler, args, result interface{}) error {
	a := args.(*item.ItemGetItemsArgs)
	r := result.(*item.ItemGetItemsResult)
	reply, err := handler.(item.Item).GetItems(ctx, a.Req)
	if err != nil {
		return err
	}
	r.Success = reply
	return nil
}
func NewMockArgs() interface{} {
	return &item.ItemGetItemsArgs{}
}

func NewMockResult() interface{} {
	return &item.GetItemsRequest{}
}

func main() {
	// 以下代码应该可以通过绕过 kitex/codec 而删除掉
	remote.PutPayloadCode(serviceinfo.Thrift, thrift.NewThriftCodec())
	mockSvrRPCInfo := rpcinfo.NewRPCInfo(rpcinfo.EmptyEndpointInfo(),
		rpcinfo.FromBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "mockServiceName"}),
		rpcinfo.NewServerInvocation(),
		rpcinfo.NewRPCConfig(), rpcinfo.NewRPCStats())

	svcInfo := &serviceinfo.ServiceInfo{
		ServiceName: "MockServiceName",
		Methods: map[string]serviceinfo.MethodInfo{
			"GetItems": serviceinfo.NewMethodInfo(mockHandler, NewMockArgs, NewMockResult, false),
		},
		Extra: map[string]interface{}{
			"PackageName": "mock",
		},
	}
	_ = svcInfo
	_ = mockSvrRPCInfo
	// 以上代码仅为了复用 kitex/codec 而存在

	buff, err := os.ReadFile("results.json")
	if err != nil {
		panic(err)
	}
	var results = map[string]interface{}{}
	err = json.Unmarshal(buff, &results)
	if err != nil {
		panic(err)
	}
	for _, v := range results["results"].([]interface{}) {
		node, err := sonic.GetFromString(v.(string), "request")
		if err != nil {
			panic(err)
		}
		x, err := node.Interface()
		if err != nil {
			panic(err)
		}

		b, err := base64.StdEncoding.DecodeString(x.(string))
		if err != nil {
			panic(err)
		}

		cc := codec.NewDefaultCodec()
		in := remote.NewReaderBuffer(b)
		sg := remote.NewMessageWithNewer(svcInfo, mockSvrRPCInfo, remote.Call, remote.Server)
		err = cc.Decode(context.Background(), sg, in)
		if err != nil {
			panic(err)
		}
		log.Print(sg.Data())
	}

}
