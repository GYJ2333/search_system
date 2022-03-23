package main

import (
	"context"
	"fmt"

	searchPb "github.com/GYJ2333/search_system/pb/search_query_feature"
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

	offlineClient := searchPb.NewFeatureService("Search_System", service.Client())

	rsp, err := offlineClient.Set(context.Background(), &searchPb.OfflineRequest{
		// // 添加一个query
		// Type: searchPb.SetType_TYPE_ADD,
		// Querys: []*searchPb.Query{
		// 	{
		// 		QueryId:   "0926",
		// 		QueryName: "xx",
		// 		Kind:      "食物",
		// 		Feature:   `{"味道":"甜", "价格":"奢侈", "颜色":"白", "口感":"软","品牌":"小象","品类":"糖"}`,
		// 	},
		// },
		// // 再添加一个query  有相同特征
		// Type: searchPb.SetType_TYPE_UPDATE,
		// Querys: []*searchPb.Query{
		// 	{
		// 		QueryId:   "1030",
		// 		QueryName: "gyj",
		// 		Kind:      "食物",
		// 		Feature:   `{"味道":"甜", "价格":"奢侈", "颜色":"黄", "口感":"硬","品牌":"小象","品类":"糖"}`,
		// 	},
		// },
		// // 修改相同特征为不同特征
		// Type: searchPb.SetType_TYPE_UPDATE,
		// Querys: []*searchPb.Query{
		// 	{
		// 		QueryId:   "1030",
		// 		QueryName: "gyj",
		// 		Kind:      "食物",
		// 		Feature:   `{"味道":"甜", "价格":"奢侈", "颜色":"黄", "口感":"硬","品牌":"大象","品类":"糖"}`,
		// 	},
		// },
		// // 删除第二个query  先删第二个query是为了验证 索引表中靠后的queryID能够正常删除
		// Type: searchPb.SetType_TYPE_DELETE,
		// Querys: []*searchPb.Query{
		// 	{
		// 		QueryId: "1030",
		// 	},
		// },
		// // 删除第一个query
		// Type: searchPb.SetType_TYPE_DELETE,
		// Querys: []*searchPb.Query{
		// 	{
		// 		QueryId: "0926",
		// 	},
		// },
	})

	fmt.Printf("rsp:(%v) err(%v)", rsp, err)
}
