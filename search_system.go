package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/GYJ2333/search_system/common"
	"github.com/GYJ2333/search_system/log"
	"github.com/GYJ2333/search_system/operator"
	searchPb "github.com/GYJ2333/search_system/pb/search_query_feature"
	"github.com/GYJ2333/search_system/storage"
)

var wg sync.WaitGroup

// main struct
type Searcher struct {
	queryRootPath, indexRootPath, profileRootPath string

	processer    operator.Processer
	queryStorage *storage.SProxy
	indexStorage *storage.SProxy
}

func makeOfflineRspHeader(rsp *searchPb.OfflineResponse, code uint32, err string) {
	rsp.Header = &searchPb.ResponseHeader{}
	rsp.Header.Code = code
	rsp.Header.Err = err
}

// 判断请求类型，调用相应的函数
func (s *Searcher) Set(ctx context.Context, req *searchPb.OfflineRequest, rsp *searchPb.OfflineResponse) error {
	if req.Querys == nil {
		log.FeatureLogger.Println("Empty req")
		makeOfflineRspHeader(rsp, 1, "empty req")
		return errors.New("empty req, set failed")
	}

	var err error
	switch req.Type {
	case searchPb.SetType_TYPE_ADD, searchPb.SetType_TYPE_UPDATE:
		if err = s.Write(req.Querys, rsp); err != nil {
			log.FeatureLogger.Printf("Add|Update data err(%v)", err)
			makeOfflineRspHeader(rsp, 1, fmt.Sprintf("Add|Update data err(%v)", err))
			return err
		}
	case searchPb.SetType_TYPE_DELETE:
		if err = s.Delete(req.Querys); err != nil {
			log.FeatureLogger.Printf("Delete data err(%v)", err)
			makeOfflineRspHeader(rsp, 1, fmt.Sprintf("Delete data err(%v)\n", err))
			return err
		}
	case searchPb.SetType_TYPE_UNKNOWN:
		fallthrough
	default:
		log.FeatureLogger.Println("Unkown type of set request")
		makeOfflineRspHeader(rsp, 1, err.Error())
		return errors.New("unknow set type, set failed")
	}

	makeOfflineRspHeader(rsp, 0, "")
	return nil
}

func (s *Searcher) Get(ctx context.Context, req *searchPb.OnlineRequest, rsp *searchPb.OnlineResponse) error {
	return nil
}

// 初始化各种算子
func (s *Searcher) Init(queryRootPath, indexRootPath, profileRootPath string) error {
	s.queryRootPath = queryRootPath
	s.indexRootPath = indexRootPath
	s.profileRootPath = profileRootPath
	s.indexStorage = &storage.SProxy{}
	s.queryStorage = &storage.SProxy{}
	s.indexStorage.Init(s.indexRootPath)
	s.queryStorage.Init(s.queryRootPath)
	return s.RegisterProcesser(&operator.DefaultProcesser{})
}

// 写入/更新数据
func (s *Searcher) Write(querys []*searchPb.Query, rsp *searchPb.OfflineResponse) error {
	if querys == nil {
		log.FeatureLogger.Println("empty data")
		return fmt.Errorf("empty data")
	}

	ch := make(chan pkg, len(querys))
	for i, v := range querys {
		wg.Add(1)
		go s.innerWrite(v, i, ch)
	}

	// 牺牲一些时间成本  串行调整rsp的顺序 使其与req相同  可以尝试直接使用数组加锁的方式来解决此问题  但是也需要考虑效率问题
	rsp.QueryStatus = make([]*searchPb.Status, len(querys))
	for flag := true; flag; {
		select {
		case v := <-ch:
			rsp.QueryStatus[v.position] = &v.status
		case <-time.After(time.Second * 1):
			flag = false
		}
	}

	wg.Wait()
	return nil
}

type pkg struct {
	position int
	status   searchPb.Status
}

// 对query数据的Add/Update
func (s *Searcher) innerWrite(query *searchPb.Query, position int, resChan chan<- pkg) {
	res := pkg{
		position: position,
		status:   searchPb.Status{QueryId: query.QueryId},
	}
	defer func() {
		resChan <- res
		wg.Done()
	}()
	// 读取旧数据
	storageData, _ := s.queryStorage.Read(query.QueryId)
	// if err != nil {
	// 	log.FeatureLogger.Printf("Read query(%s) err(%v)", query.QueryId, err)
	// 	res.status.Msg = err.Error()
	// 	return
	// }
	// 数据处理 得到待写入数据
	newQueryData, err := s.processer.ProcessData(query, storageData)
	if err != nil {
		log.FeatureLogger.Printf("Process data err(%v)", err)
		res.status.Msg = err.Error()
		return
	}
	// 写入待写入数据
	if err := s.queryStorage.Write(query.QueryId, newQueryData); err != nil {
		log.FeatureLogger.Printf("Write new query data err(%v)", err)
		res.status.Msg = err.Error()
		return
	}
	// 正常结束
	res.status.Ok = true
	res.status.Msg = "success"
}

// 删除数据
func (s *Searcher) Delete(querys []*searchPb.Query) error {
	return nil
}

// 注册processer
func (s *Searcher) RegisterProcesser(p operator.Processer) error {
	if err := p.Init(common.FeatureYaml, s.indexStorage); err != nil {
		log.FeatureLogger.Printf("Processer init err(%v)", err)
		return err
	}
	s.processer = p
	log.FeatureLogger.Println("Processer init success")
	return nil
}
