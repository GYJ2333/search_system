package operator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/GYJ2333/search_system/common"
	"github.com/GYJ2333/search_system/feature"
	"github.com/GYJ2333/search_system/log"
	searchPb "github.com/GYJ2333/search_system/pb/search_query_feature"
	"github.com/GYJ2333/search_system/storage"
	"github.com/GYJ2333/search_system/tool"
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

	return p.innerProcess(query, storageQueryData)

	// // 存量数据中没有查询到待写入数据，直接将待写入数据写入
	// if storageQueryData == nil {
	// 	if newQueryData, err = p.marshalQueryData(query); err != nil {
	// 		log.FeatureLogger.Printf("Marshal query data err(%v)", err)
	// 		err = fmt.Errorf("marshal query data err(%v)", err)
	// 		return
	// 	}
	// 	// 并发写入倒排索引表中

	// 	return
	// }
	// return p.mergeData(query, storageQueryData)
}

// 1.将新旧数据解码为map
// 2.并发处理索引数据和query数据
func (p *DefaultProcesser) innerProcess(query *searchPb.Query, storageQueryData []byte) (newQueryData []byte, err error) {
	defer wg.Wait()
	newDataMap, oldDataMap, err := p.marshalQueryData(query, storageQueryData)
	if err != nil {
		log.FeatureLogger.Printf("Marshal query data to map err(%v)", err)
		return nil, fmt.Errorf("marshal query data to map err(%v)", err)
	}
	// 不论旧数据不为空，merge新旧数据
	wg.Add(1)
	go p.processIndex(newDataMap, oldDataMap)
	newQueryData, err = p.processQuery(newDataMap, oldDataMap)
	return
}

// 将pb协议结构的数据转为存储结构
// TODO 硬编码能不能修正一下
func (p *DefaultProcesser) marshalQueryData(query *searchPb.Query, storageQueryData []byte) (map[string]string, map[string]interface{}, error) {
	var parser fastjson.Parser
	oldDataMap := map[string]interface{}{}
	newDataMap := map[string]string{}
	newDataMap["query_id"] = query.QueryId
	newDataMap["query_name"] = query.QueryName
	newDataMap["kind"] = query.Kind
	// 从请求feature中获取特征值
	v, err := parser.Parse(query.Feature)
	if err != nil {
		log.FeatureLogger.Printf("Parse features err(%v)", err)
		return nil, nil, fmt.Errorf("parse features err(%v)", err)
	}
	for _, feature := range p.featureConf[query.Kind] {
		data := v.GetStringBytes(feature)
		newDataMap[feature] = *(tool.Bytes2String(data))
	}

	if storageQueryData != nil {
		// 将原数据解析为map
		err = json.Unmarshal(storageQueryData, &oldDataMap)
		if err != nil {
			log.FeatureLogger.Printf("Unmarshal rowData(%s) err(%v)", query.QueryId, err)
			return nil, nil, fmt.Errorf("unmarshal rowData(%s) err(%v)", query.QueryId, err)
		}
	}

	return newDataMap, oldDataMap, nil
}

// 处理倒排索引表
func (p *DefaultProcesser) processIndex(newDataMap map[string]string, oldDataMap map[string]interface{}) {
	defer wg.Done()
	for _, f := range p.featureConf[newDataMap["kind"]] {
		if newDataMap[f] == oldDataMap[f] {
			continue
		}
		wg.Add(1)
		go p.addIndex(newDataMap[f], newDataMap["query_id"])
		if oldDataMap[f] == nil {
			continue
		}
		wg.Add(1)
		go p.deleteIndex(oldDataMap[f].(string), oldDataMap["query_id"].(string))
	}
}

// 更新索引文件
func (p *DefaultProcesser) addIndex(key, queryId string) {
	defer wg.Done()
	data, _ := p.indexStorage.Read(key)
	// data为空，直接创建新的索引文件
	if data == nil {
		if err := p.indexStorage.Write(key, tool.String2Bytes(queryId)); err != nil {
			log.FeatureLogger.Printf("Write index file err(%v)", err)
		}
		return
	}
	// data不为空，merge新旧数据，写入merge后的数据
	data = append(data, tool.String2Bytes(","+queryId)...)
	if err := p.indexStorage.Write(key, data); err != nil {
		log.FeatureLogger.Printf("Write index file err(%v)", err)
	}
}

// 删除索引文件中的指定key
func (p *DefaultProcesser) deleteIndex(key, queryId string) {
	defer wg.Done()
	data, _ := p.indexStorage.Read(key)
	// data为空，报错日志并返回
	if data == nil {
		log.FeatureLogger.Printf("Terrible thing! Index and Query's feature don't match!!!!!")
		return
	}
	// data不为空，merge新旧数据，写入merge后的数据
	newData := strings.Replace(*tool.Bytes2String(data), queryId, "", 1)
	if err := p.indexStorage.Write(key, tool.String2Bytes(newData)); err != nil {
		log.FeatureLogger.Printf("Write index file err(%v)", err)
	}
}

// 处理query数据
func (p *DefaultProcesser) processQuery(newDataMap map[string]string, oldDataMap map[string]interface{}) (newQueryData []byte, err error) {
	resMap := map[string]string{}
	for k, v := range oldDataMap {
		resMap[k] = v.(string)
	}
	for k, v := range newDataMap {
		resMap[k] = v
	}
	return json.Marshal(resMap)
}
