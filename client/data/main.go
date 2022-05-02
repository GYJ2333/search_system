package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var f = []string{"味道", "价格", "颜色", "口感", "品牌", "品类"}

var 味道 = []string{"酸", "甜", "苦", "辣"}
var 价格 = []string{"贵", "便宜"}
var 颜色 = []string{"红", "绿", "蓝"}
var 口感 = []string{"软糯", "嘎嘣脆", "清爽", "绵密"}
var 品牌 = []string{"旺旺", "徐福记", "良品铺子", "周黑鸭"}
var 品类 = []string{"糖", "饼干", "面包", "饮料"}
var feature = [][]string{味道, 价格, 颜色, 口感, 品牌, 品类}

var fp *os.File

func main() {
	fp, _ = os.OpenFile("./data", os.O_APPEND|os.O_CREATE, 0644)
	dataMaker(map[string]string{}, 0)
}

func dataMaker(m map[string]string, i int) {
	if i == len(f) {
		data, err := json.Marshal(m)
		if err != nil {
			fmt.Printf("oh god plz no! the map is (%v)", m)
			return
		}
		fp.Write(data)
		fp.WriteString("\n")
		return
	}
	for _, v := range feature[i] {
		m[f[i]] = v
		dataMaker(m, i+1)
	}
}
