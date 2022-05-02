package operator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/GYJ2333/search_system/common"
	"github.com/GYJ2333/search_system/feature"
	"github.com/GYJ2333/search_system/log"
	searchPb "github.com/GYJ2333/search_system/pb/search_query_feature"
	"github.com/GYJ2333/search_system/storage"
	"github.com/GYJ2333/search_system/tool"
	"github.com/valyala/fastjson"
	"gopkg.in/yaml.v2"
)

var offlineWG, onlineWG sync.WaitGroup

type Processer interface {
	// Processer初始化、配置文件的加载
	Init(filename string, indexStorage, profileStorage, queryStorage *storage.SProxy) error
	// ProcessData 处理离线写入数据
	// 入参为待写入query，以及存储中的query数据
	// 返回值为待写入存储的query数据以及索引表数据
	ProcessData(query *searchPb.Query, storageQueryData []byte) ([]byte, error)
	// 删除query，即删除索引表中对应的queryId
	DeleteData(query *searchPb.Query, storageQueryData []byte) error
	// 配套的打分手段
	Search(userId string, features []string) ([]string, error)
	// 配套的画像更新
	UpdateProfile(userId, queryId string) error
}

type DefaultProcesser struct {
	featureConf                                map[string][]string
	indexStorage, profileStorage, queryStorage *storage.SProxy
}

func (p *DefaultProcesser) Init(filename string, indexStorage, profileStorage, queryStorage *storage.SProxy) error {
	featureConf := &feature.Conf{}
	p.featureConf = make(map[string][]string)
	p.indexStorage = indexStorage
	p.profileStorage = profileStorage
	p.queryStorage = queryStorage

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
}

// 1.将新旧数据解码为map
// 2.并发处理索引数据和query数据
func (p *DefaultProcesser) innerProcess(query *searchPb.Query, storageQueryData []byte) (newQueryData []byte, err error) {
	defer offlineWG.Wait()
	newDataMap, oldDataMap, err := p.marshalQueryData(query, storageQueryData)
	if err != nil {
		log.FeatureLogger.Printf("Marshal query data to map err(%v)", err)
		return nil, fmt.Errorf("marshal query data to map err(%v)", err)
	}
	// 不论旧数据不为空，merge新旧数据
	offlineWG.Add(1)
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
			log.FeatureLogger.Printf("Unmarshal oldData(%s) err(%v)", query.QueryId, err)
			return nil, nil, fmt.Errorf("unmarshal oldData(%s) err(%v)", query.QueryId, err)
		}
	}

	return newDataMap, oldDataMap, nil
}

// 处理倒排索引表
func (p *DefaultProcesser) processIndex(newDataMap map[string]string, oldDataMap map[string]interface{}) {
	defer offlineWG.Done()
	for _, f := range p.featureConf[newDataMap["kind"]] {
		if newDataMap[f] == oldDataMap[f] {
			continue
		}
		offlineWG.Add(1)
		go p.addIndex(newDataMap[f], newDataMap["query_id"])
		if oldDataMap[f] == nil {
			continue
		}
		offlineWG.Add(1)
		go p.deleteIndex(oldDataMap[f].(string), oldDataMap["query_id"].(string))
	}
}

// 更新索引文件
func (p *DefaultProcesser) addIndex(key, queryId string) {
	defer offlineWG.Done()
	data, _ := p.indexStorage.Read(key)
	// data为空，直接创建新的索引文件
	if data == nil {
		if err := p.indexStorage.Write(key, tool.String2Bytes(queryId)); err != nil {
			log.FeatureLogger.Printf("Write index file err(%v)", err)
		}
		return
	}
	// 将待写入的queryID写入索引表中
	data = append(data, tool.String2Bytes(","+queryId)...)
	if err := p.indexStorage.Write(key, data); err != nil {
		log.FeatureLogger.Printf("Write index file err(%v)", err)
	}
}

// 删除索引文件中的指定key
func (p *DefaultProcesser) deleteIndex(key, queryId string) {
	defer offlineWG.Done()
	data, _ := p.indexStorage.Read(key)
	// data为空，报错日志并返回
	if data == nil || len(data) < len(queryId) {
		log.FeatureLogger.Printf("Terrible thing! Index(%s) and Query(%s)'s feature don't match!!!!!", key, queryId)
		return
	}
	// data不为空，merge新旧数据，写入merge后的数据
	if *tool.Bytes2String(data[:len(queryId)]) == queryId {
		if len(data) == len(queryId) {
			data = []byte{}
		} else {
			data = data[len(queryId)+1:]
		}
	} else {
		data = tool.String2Bytes(strings.Replace(*tool.Bytes2String(data), ","+queryId, "", 1))
	}
	if len(data) == 0 {
		p.indexStorage.Delete(key)
	} else if err := p.indexStorage.Write(key, tool.String2Bytes(string(data))); err != nil {
		log.FeatureLogger.Printf("Write index(%s) file after delete(%s) err(%v)", key, queryId, err)
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

// 删除一个query
func (p *DefaultProcesser) DeleteData(query *searchPb.Query, storageQueryData []byte) error {
	defer offlineWG.Wait()
	if storageQueryData == nil {
		return fmt.Errorf("query(%s) has been deleted", query.QueryId)
	}
	oldDataMap := map[string]interface{}{}
	if err := json.Unmarshal(storageQueryData, &oldDataMap); err != nil {
		log.FeatureLogger.Printf("Unmarshal oldData(%s) err(%v)", query.QueryId, err)
		return fmt.Errorf("unmarshal oldData(%s) err(%v)", query.QueryId, err)
	}
	for _, v := range oldDataMap {
		offlineWG.Add(1)
		go p.deleteIndex(v.(string), query.QueryId)
	}
	return nil
}

// 查询
func (p *DefaultProcesser) Search(userId string, features []string) ([]string, error) {
	userProfile := p.getProfile(userId)
	bestQuerys, err := p.getBestQuerys(features)
	if err != nil {
		return nil, err
	}
	log.FeatureLogger.Printf("Debug bestQuerys(%v)", bestQuerys)
	result, err := p.innerSearch(userProfile, bestQuerys)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 获取用户画像
func (p *DefaultProcesser) getProfile(userId string) (res map[string]map[string][]string) {
	profileData, err := p.profileStorage.Read(userId)
	if err != nil {
		log.ProfileLogger.Printf("Read user(%s)'s profile err(%v)", userId, err)
		return p.makeDefaultProfile()
	}
	res = map[string]map[string][]string{}
	if err := json.Unmarshal(profileData, &res); err != nil {
		log.ProfileLogger.Printf("Read user(%s)'s profile err(%v)", userId, err)
		return p.makeDefaultProfile()
	}
	return
}

// 制作默认用户画像
func (p *DefaultProcesser) makeDefaultProfile() (res map[string]map[string][]string) {
	defaultMap := make(map[string]map[string][]string, 0)
	for k, v := range p.featureConf {
		defaultMap[k] = map[string][]string{}
		for _, d := range v {
			defaultMap[k][d] = make([]string, 3)
		}
	}
	log.ProfileLogger.Printf("default profile:(%v)", defaultMap)
	return defaultMap
}

// 并发安全map
type safeMap struct {
	m map[string]int
	sync.RWMutex
}

// 取最优query集，即倒排索引表中出现次数最多的query
// TODO 这里应该是定数量返回优质queryID 还是设定为出现频率最高的queryID组（这里可能有问题是如果只有一个query出现的频率最高，那其他的可能就不在考虑范围内了）
func (p *DefaultProcesser) getBestQuerys(features []string) (res []string, err error) {
	// 并发安全的map，用做频率桶，统计索引表交集
	bestQueryMap := &safeMap{
		m:       map[string]int{},
		RWMutex: sync.RWMutex{},
	}
	// 并发统计交集
	for _, d := range features {
		onlineWG.Add(1)
		go p.countQuery(d, bestQueryMap)
	}
	onlineWG.Wait()
	rankMap := map[int][]string{}
	for k, v := range bestQueryMap.m {
		rankMap[v] = append(rankMap[v], k)
	}
	for i := len(features); i > 0; i-- {
		if tmp, ok := rankMap[i]; ok {
			log.FeatureLogger.Printf("Debug rankMap[%d] is (%v)", i, tmp)
			res = append(res, tmp...)
			// 已经选取出不少于100个query
			if len(res) >= 100 {
				return res, nil
			}
		}
	}
	if res == nil {
		return res, fmt.Errorf("no similar query in data base")
	}
	return res, nil
}

// 求索引表交集
func (p *DefaultProcesser) countQuery(feature string, bestQueryMap *safeMap) {
	defer onlineWG.Done()
	data, err := p.indexStorage.Read(feature)
	if err != nil {
		log.FeatureLogger.Printf("Read feature(%s) err(%v)", feature, err)
		return
	}
	querys := strings.Split(*tool.Bytes2String(data), ",")
	for _, d := range querys {
		bestQueryMap.Lock()
		bestQueryMap.m[d]++
		bestQueryMap.Unlock()
	}
}

// 记录每个query的分数
type Q struct {
	query string
	score float32
}

// 排序用
type Sorter []Q

func (q Sorter) Len() int {
	return len(q)
}
func (q Sorter) Less(i, j int) bool {
	return q[i].score < q[j].score
}
func (q Sorter) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

// 打分并排序
func (p *DefaultProcesser) innerSearch(profile map[string]map[string][]string, querys []string) ([]string, error) {
	ch := make(chan Q, len(querys))
	for _, d := range querys {
		onlineWG.Add(1)
		go p.grade(ch, d, profile)
	}

	tmp := make([]Q, len(querys))
	for i := 0; i < len(querys); i++ {
		select {
		case v := <-ch:
			tmp[i] = v
		// 这里意味着每次至少要100ms以上的时间
		case <-time.After(time.Millisecond * 100):
			i = len(querys)
		}
	}

	log.FeatureLogger.Printf("Debug querys before sort(%v)", tmp)
	onlineWG.Wait()
	sort.Sort(sort.Reverse(Sorter(tmp)))
	log.FeatureLogger.Printf("Debug querys after sort(%v)", tmp)

	res := []string{}
	for i := 0; i < 20 && i < len(tmp); i++ {
		res = append(res, tmp[i].query)
	}
	return res, nil
}

// 对一个query进行打分
func (p *DefaultProcesser) grade(ch chan Q, queryID string, profile map[string]map[string][]string) {
	defer onlineWG.Done()
	feature, err := p.queryStorage.Read(queryID)
	if err != nil {
		log.FeatureLogger.Printf("Terrible thing! Index and Query(%s)'s feature don't match!!!!!", queryID)
		return
	}
	featureMap := map[string]string{}
	if err = json.Unmarshal(feature, &featureMap); err != nil {
		log.FeatureLogger.Printf("Query(%s) unmarshal err(%v)", queryID, err)
		return
	}
	res := Q{query: queryID}
	res.score = p.innerGrade(featureMap, profile[featureMap["kind"]])
	ch <- res
}

// 打分详细过程
func (p *DefaultProcesser) innerGrade(queryFeature map[string]string, profile map[string][]string) (score float32) {
	// 新增品类时旧有的profile内并没有这部分画像  需要忽略
	if profile == nil {
		return
	}
	for k, v := range queryFeature {
		for i, d := range profile[k] {
			if d == v {
				score += float32(len(profile[k]) - i)
			}
		}
	}
	return
}

// 配套的画像更新
func (p *DefaultProcesser) UpdateProfile(userId, queryId string) error {
	profile := p.getProfile(userId)
	featureData, err := p.queryStorage.Read(queryId)
	if err != nil {
		log.ProfileLogger.Printf("Read query(%s)'s feature err(%v)", queryId, err)
		return fmt.Errorf("read query(%s)'s feature err(%v)", queryId, err)
	}
	featureMap := map[string]string{}
	if err = json.Unmarshal(featureData, &featureMap); err != nil {
		log.ProfileLogger.Printf("Query(%s) unmarshal err(%v)", queryId, err)
		return fmt.Errorf("query(%s) unmarshal err(%v)", queryId, err)
	}
	newProfile, err := p.innerUpdate(profile, featureMap)
	log.ProfileLogger.Printf("Debug new profile (%s)", newProfile)
	if err != nil {
		return err
	}
	if err = p.profileStorage.Write(userId, newProfile); err != nil {
		log.ProfileLogger.Printf("Write new profile(%s) err(%v)", userId, newProfile)
		return fmt.Errorf("write new profile(%s) err(%v)", userId, newProfile)
	}
	return nil
}

// 画像更新
func (p *DefaultProcesser) innerUpdate(profile map[string]map[string][]string, queryFeature map[string]string) ([]byte, error) {
	innerProfile, ok := profile[queryFeature["kind"]]
	// 新增品类时旧有的profile内并没有这部分画像  需要新建
	if !ok {
		newProfile := map[string][]string{}
		for _, d := range p.featureConf[queryFeature["kind"]] {
			newProfile[d] = make([]string, 3)
		}
		profile[queryFeature["kind"]] = newProfile
		innerProfile = profile[queryFeature["kind"]]
	}
	for _, d := range p.featureConf[queryFeature["kind"]] {
		if innerProfile[d] == nil {
			innerProfile[d] = make([]string, 3)
		}
		innerProfile[d] = append([]string{queryFeature[d]}, innerProfile[d][:len(innerProfile[d])-1]...)
	}
	return json.Marshal(profile)
}
