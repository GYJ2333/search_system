package operator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"unsafe"

	"github.com/GYJ2333/search_system/common"
	"github.com/GYJ2333/search_system/feature"
	"github.com/GYJ2333/search_system/log"
	searchPb "github.com/GYJ2333/search_system/pb/search_query_feature"
	"github.com/GYJ2333/search_system/storage"
	"github.com/valyala/fastjson"
	"gopkg.in/yaml.v2"
)

var wg sync.WaitGroup

type Processer interface {
	// Processer初始化、配置文件的加载
	Init(filename string, sProxy *storage.SProxy) error
	// ProcessData 处理离线写入数据
	// 入参为待写入query，以及存储中的query数据
	// 返回值为待写入存储的query数据以及索引表数据
	ProcessData(query *searchPb.Query, storageQueryData []byte) ([]byte, error)
	// 配套的打分手段
}

type DefaultProcesser struct {
	featureConf  map[string][]string
	indexStorage *storage.SProxy
}

func (p *DefaultProcesser) Init(filename string, sProxy *storage.SProxy) error {
	featureConf := &feature.Conf{}
	p.featureConf = make(map[string][]string)
	p.indexStorage = sProxy

	data, err := ioutil.ReadFile(common.YamlRootPath + common.FeatureYaml)
	if err != nil {
		log.FeatureLogger.Printf("Read feature yaml file err(%v)", err)
		return err
	}
	if err := yaml.Unmarshal(data, featureConf); err != nil {
		log.FeatureLogger.Printf("Unmarshal feature yaml err(%v)", err)
		return err
	}

	for _, v := range featureConf.Kind {
		p.featureConf[v.Name] = v.Feature
	}
	log.FeatureLogger.Printf("Feature yaml(%v)", p.featureConf)

	return nil
}

func (p *DefaultProcesser) ProcessData(query *searchPb.Query, storageQueryData []byte) (newQueryData []byte, err error) {
	if _, ok := p.featureConf[query.Kind]; !ok {
		log.FeatureLogger.Printf("Search system doesn't support this kind(%s) of query", query.Kind)
		err = fmt.Errorf("search system doesn't support this kind(%s) of query", query.Kind)
		return
	}
	// 存量数据中没有查询到待写入数据，直接将待写入数据写入
	if storageQueryData == nil {
		if newQueryData, err = p.marshalQueryData(query); err != nil {
			log.FeatureLogger.Printf("Marshal query data err(%v)", err)
			err = fmt.Errorf("marshal query data err(%v)", err)
			return
		}
		return
	}
	return p.mergeData(query, storageQueryData)
}

// 将pb协议结构的数据转为存储结构
// TODO 硬编码能不能修正一下
func (p *DefaultProcesser) marshalQueryData(query *searchPb.Query) ([]byte, error) {
	resMap := make(map[string]string, 0)
	resMap["query_id"] = query.QueryId
	resMap["query_name"] = query.QueryName
	resMap["kind"] = query.Kind
	var parser fastjson.Parser
	v, err := parser.Parse(query.Feature)
	if err != nil {
		log.FeatureLogger.Printf("Parse features err(%v)", err)
		return nil, fmt.Errorf("parse features err(%v)", err)
	}
	for _, feature := range p.featureConf[query.Kind] {
		data := v.GetStringBytes(feature)
		resMap[feature] = *(*string)(unsafe.Pointer(&data))
	}

	return json.Marshal(resMap)
}

func (p *DefaultProcesser) mergeData(query *searchPb.Query, storageQueryData []byte) (newQueryData []byte, err error) {
	// 将原数据解析为map
	rowData := map[string]interface{}{}
	err = json.Unmarshal(storageQueryData, &rowData)
	if err != nil {
		log.FeatureLogger.Printf("Unmarshal rowData(%s) err(%v)", query.QueryId, err)
		return nil, fmt.Errorf("unmarshal rowData(%s) err(%v)", query.QueryId, err)
	}

	var parser fastjson.Parser
	v, err := parser.Parse(query.Feature)
	if err != nil {
		log.FeatureLogger.Printf("Parse features err(%vn", err)
		return nil, fmt.Errorf("parse features err(%v)", err)
	}
	for _, feature := range p.featureConf[query.Kind] {
		data := v.GetStringBytes(feature)
		if len(data) == 0 {
			continue
		}
		rowData[feature] = *(*string)(unsafe.Pointer(&data))
	}
	return json.Marshal(rowData)
}
