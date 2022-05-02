package main

import (
	"context"
	"fmt"
	"io/ioutil"

	searchPb "github.com/GYJ2333/search_system/pb/search_query_feature"
	"github.com/GYJ2333/search_system/tool"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

func main() {
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	service := micro.NewService(
		micro.Name("Search_System"),
		micro.Version("latest"),
		micro.Metadata(map[string]string{
			"type": "search_system",
		}),
		micro.Registry(reg),
	)

	onlineClient := searchPb.NewFeatureService("Search_System", service.Client())

	rsp, err := onlineClient.Get(context.Background(), &searchPb.OnlineRequest{
		UserId:   "gyj",
		Features: []string{"便宜","徐福记","饼干"},
		Ext:      map[string]string{},
	})

	// 如果没找到相应数据 直接结束
	if err != nil {
		fmt.Printf("rsp:(%v) err(%v)", rsp, err)
		return
	}

	for i, d := range rsp.QueryIds {
		fmt.Printf("%d. %s\n", i, getQueryData(d))
	}

	num := 0
	fmt.Printf("请输入您的选择: ")
	fmt.Scanln(&num)

	pRsp, err := onlineClient.Chose(context.Background(), &searchPb.ChoseRequest{
		UserId:  "gyj",
		QueryId: rsp.QueryIds[num],
	})

	fmt.Printf("rsp:(%v) err(%v)", pRsp, err)
}

func getQueryData(queryId string) string {
	rowData, err := ioutil.ReadFile("../../data/query/" + queryId)
	if err != nil {
		fmt.Printf("read file(%s) err(%v)", queryId, err)
		return ""
	}

	decompressedData, err := tool.Decompress(rowData)
	if err != nil {
		// log.StorageLogger.Printf("Decompress file(%s) err(%v)", queryId, err)
		// return nil, fmt.Errorf("decompress file(%s) err(%v)", queryId, err)
		return string(rowData)
	}
	return string(decompressedData)
}
