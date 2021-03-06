package main

import (
	"context"
	"errors"
	"fmt"
	"os"
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

	processer      operator.Processer
	queryStorage   *storage.SProxy
	indexStorage   *storage.SProxy
	profileStorage *storage.SProxy
}

func makeRspHeader(rsp interface{}, code uint32, err string) {
	var header *searchPb.ResponseHeader
	switch t := rsp.(type) {
	case *searchPb.OfflineResponse:
		t.Header = &searchPb.ResponseHeader{}
		header = t.Header
	case *searchPb.OnlineResponse:
		t.Header = &searchPb.ResponseHeader{}
		header = t.Header
	case *searchPb.ChoseResponse:
		t.Header = &searchPb.ResponseHeader{}
		header = t.Header
	}
	header.Code = code
	header.Err = err
}

// 判断请求类型，调用相应的函数
func (s *Searcher) Set(ctx context.Context, req *searchPb.OfflineRequest, rsp *searchPb.OfflineResponse) error {
	if req.Querys == nil {
		log.FeatureLogger.Println("Empty req")
		makeRspHeader(rsp, 1, "empty req")
		return errors.New("empty req, set failed")
	}

	var err error
	switch req.Type {
	case searchPb.SetType_TYPE_ADD, searchPb.SetType_TYPE_UPDATE:
		if err = s.Write(req.Querys, rsp); err != nil {
			log.FeatureLogger.Printf("Add|Update data err(%v)", err)
			makeRspHeader(rsp, 1, fmt.Sprintf("Add|Update data err(%v)", err))
			return err
		}
	case searchPb.SetType_TYPE_DELETE:
		if err = s.Delete(req.Querys, rsp); err != nil {
			log.FeatureLogger.Printf("Delete data err(%v)", err)
			makeRspHeader(rsp, 1, fmt.Sprintf("Delete data err(%v)\n", err))
			return err
		}
	case searchPb.SetType_TYPE_UNKNOWN:
		fallthrough
	default:
		log.FeatureLogger.Println("Unkown type of set request")
		makeRspHeader(rsp, 1, err.Error())
		return errors.New("unknow set type, set failed")
	}

	makeRspHeader(rsp, 0, "")
	return nil
}

func (s *Searcher) Get(ctx context.Context, req *searchPb.OnlineRequest, rsp *searchPb.OnlineResponse) error {
	if req == nil || req.UserId == "" || req.Features == nil {
		log.FeatureLogger.Println("Invalid req")
		makeRspHeader(rsp, 1, "invalid req")
		return errors.New("invalid req, get failed")
	}
	res, err := s.processer.Search(req.UserId, req.Features)
	if err != nil {
		log.FeatureLogger.Printf("Search (%v) by user(%s) err(%v)", req.Features, req.UserId, err)
		makeRspHeader(rsp, 1, fmt.Sprintf("search failed(%v)", err))
		return err
	}
	rsp.QueryIds = res

	makeRspHeader(rsp, 0, "")
	return nil
}

// 初始化存储以及算子
func (s *Searcher) Init(queryRootPath, indexRootPath, profileRootPath string) error {
	os.Mkdir(queryRootPath, os.ModePerm)
	os.Mkdir(indexRootPath, os.ModePerm)
	os.Mkdir(profileRootPath, os.ModePerm)
	s.queryRootPath = queryRootPath
	s.indexRootPath = indexRootPath
	s.profileRootPath = profileRootPath
	s.indexStorage = &storage.SProxy{}
	s.queryStorage = &storage.SProxy{}
	s.profileStorage = &storage.SProxy{}
	s.indexStorage.Init(s.indexRootPath)
	s.queryStorage.Init(s.queryRootPath)
	s.profileStorage.Init(s.profileRootPath)
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
	wg.Wait()
	// 牺牲一些时间成本  串行调整rsp的顺序 使其与req相同  可以尝试直接使用数组加锁的方式 在协程中解决此问题  但是也需要考虑效率问题
	rsp.QueryStatus = make([]*searchPb.Status, len(querys))
	for flag := true; flag; {
		select {
		case v := <-ch:
			rsp.QueryStatus[v.position] = &v.status
		// 这里意味着每次至少要10ms以上的时间
		case <-time.After(time.Millisecond * 10):
			flag = false
		}
	}
	log.FeatureLogger.Print("Debug we come here")
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
	if query.QueryId == "" || query.Feature == "" || query.Kind == "" || query.QueryName == "" {
		log.FeatureLogger.Print("Invalid request query!Please check data's completeness")
		res.status.Msg = "Invalid request query!Please check data's completeness"
		return
	}
	// 读取旧数据
	storageData, _ := s.queryStorage.Read(query.QueryId)
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
func (s *Searcher) Delete(querys []*searchPb.Query, rsp *searchPb.OfflineResponse) error {
	if querys == nil {
		log.FeatureLogger.Println("empty data")
		return fmt.Errorf("empty data")
	}

	ch := make(chan pkg, len(querys))
	for i, v := range querys {
		wg.Add(1)
		go s.innerDelete(v, i, ch)
	}

	// 牺牲一些时间成本  串行调整rsp的顺序 使其与req相同  可以尝试直接使用数组加锁的方式 在协程中解决此问题  但是也需要考虑效率问题
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

// 对query数据的Delete
func (s *Searcher) innerDelete(query *searchPb.Query, position int, resChan chan<- pkg) {
	res := pkg{
		position: position,
		status:   searchPb.Status{QueryId: query.QueryId},
	}
	if query.QueryId == "" {
		log.FeatureLogger.Print("Invalid request query!Please check data's completeness")
		res.status.Msg = "Invalid request query!Please check data's completeness"
		return
	}
	defer func() {
		resChan <- res
		wg.Done()
	}()
	// 读取旧数据
	storageData, _ := s.queryStorage.Read(query.QueryId)
	// 删除对应的索引文件
	err := s.processer.DeleteData(query, storageData)
	if err != nil {
		log.FeatureLogger.Printf("Delete query(%v)'s index file err(%v)", query.QueryId, err)
		res.status.Msg = err.Error()
		return
	}
	// 删除query数据对应的文件
	if err := s.queryStorage.Delete(query.QueryId); err != nil {
		log.FeatureLogger.Printf("Delete query(%s) err(%v)", query.QueryId, err)
		res.status.Msg = err.Error()
		return
	}
	// 正常结束
	res.status.Ok = true
	res.status.Msg = "success"
}

// 注册processer
func (s *Searcher) RegisterProcesser(p operator.Processer) error {
	if err := p.Init(common.FeatureYaml, s.indexStorage, s.profileStorage, s.queryStorage); err != nil {
		log.FeatureLogger.Printf("Processer init err(%v)", err)
		return err
	}
	s.processer = p
	log.FeatureLogger.Println("Processer init success")
	return nil
}

// 画像更新
func (s *Searcher) Chose(ctx context.Context, req *searchPb.ChoseRequest, rsp *searchPb.ChoseResponse) error {
	if req == nil || req.QueryId == "" || req.UserId == "" {
		log.ProfileLogger.Println("Invalid req")
		makeRspHeader(rsp, 1, "Invalid req")
		return errors.New("invalid req")
	}
	if err := s.processer.UpdateProfile(req.UserId, req.QueryId); err != nil {
		log.ProfileLogger.Printf("Update profile err(%v)", err)
		makeRspHeader(rsp, 1, fmt.Sprintf("Update profile err(%v)", err))
		return fmt.Errorf("update profile err(%v)", err)
	}
	makeRspHeader(rsp, 0, "success")
	return nil
}
