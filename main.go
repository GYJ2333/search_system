package main

import (
	"fmt"

	"github.com/GYJ2333/search_system/common"
	"github.com/GYJ2333/search_system/log"
	searchPb "github.com/GYJ2333/search_system/pb/search_query_feature"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

func main() {
	// TODO 补齐log文件名
	if err := log.Init("./log/feature_log.log", "./log/profile_log.log", "./log/storage_log.log"); err != nil {
		fmt.Printf("logger init err(%v)\n", err)
		return
	}

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

	service.Init()

	searcher := &Searcher{}
	if err := searcher.Init(common.QueryDataPath, common.IndexDataPath, common.ProfileDataPath); err != nil {
		log.FeatureLogger.Fatalf("search system init err(%v)", err)
	}

	searchPb.RegisterFeatureHandler(service.Server(), searcher)
	if err := service.Run(); err != nil {
		fmt.Printf("server run err(%v)\n", err)
	}
}
