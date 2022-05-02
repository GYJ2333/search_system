package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

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

	offlineClient := searchPb.NewFeatureService("Search_System", service.Client())

	// // 自测时将下述代码注释
	// fp, _ := os.OpenFile("./log", os.O_APPEND|os.O_CREATE, 0777)
	// reqs := makeReq()
	// if reqs == nil {
	// 	return
	// }
	// for i, d := range reqs {
	// 	rsp, err := offlineClient.Set(context.Background(), &d)
	// 	fp.WriteString(fmt.Sprintf("query(%d) rsp:(%v) err(%v)\n", i, rsp, err))
	// }

	// 自测时使用
	rsp, err := offlineClient.Set(context.Background(), &searchPb.OfflineRequest{
		Type: searchPb.SetType_TYPE_UPDATE,
		Querys: []*searchPb.Query{
			{
				QueryId:   "兜底",
				QueryName: "兜底",
				Kind:      "食物",
				Feature:   `{"价格":"奢侈"}`,
			},
		},
	})
	fmt.Printf("rsp:(%v) err(%v)", rsp, err)
}

func makeReq() []searchPb.OfflineRequest {
	rowData, err := ioutil.ReadFile("../data/data")
	if err != nil {
		fmt.Printf("open data file err(%v)", err)
		return nil
	}
	data := tool.Bytes2String(rowData)
	features := strings.Split(*data, "\n")
	res := make([]searchPb.OfflineRequest, len(features))
	for i, d := range features {
		res[i].Type = searchPb.SetType_TYPE_ADD
		res[i].Querys = append(res[i].Querys, &searchPb.Query{
			QueryId:   fmt.Sprint(i),
			QueryName: fmt.Sprint(i),
			Kind:      "食物",
			Feature:   d,
		})
	}
	return res
}
