package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
)

var (
	LinkSearchUrl = "http://localhost:9200/milnk_index/milnk_type"
)

type LinkSearch struct {
	Url string
}

//添加link到es搜索
func (ls *LinkSearch) AddLink(link map[string]interface{}) (*http.Response, error) {
	if ls.Url == "" {
		ls.Url = LinkSearchUrl
	}
	data := fmt.Sprintf(`{"title":"%s","context":"%s","topics":"%s","user":"%s","host":"%s"}`, link["title"], link["context"], 
		strings.Replace(link["topics"].(string), ",", " ", -1), link["username"], strings.Replace(link["host"].(string), ".", " ", -1))
	resp, err := http.DefaultClient.Post(fmt.Sprintf("%s/%d", ls.Url, link["id"]), "application/json", bytes.NewBuffer([]byte(data)))
	//resp:{"ok":true,"_index":"milnk_index","_type":"milnk_type","_id":"2","_version":1}
	//fmt.Println("resp: ", resp)
	if err != nil {
		fmt.Println("AddLinkSearchError: ", err)
	}
	return resp, err
}


type SearchResult struct {
	Took        int                     `json:"took"`
	TimedOut    bool                    `json:"timed_out"`
	HitResult   SearchHitCollection     `json:"hits"`
}

type SearchHitCollection struct {
	Total     int                   `json:"total"`
	MaxScore  float64               `json:"max_score"`
	HitArray  []SearchHitItem       `json:"hits"`
}
type SearchHitItem struct {
	Index     string                  `json:"_index"`
	Type      string                  `json:"_type"`
	Id        string                  `json:"_id"`
	Score     float64                 `json:"_score"`
	//Source    SearchLinkItem          `json:"_source"`
}
type SearchLinkItem struct {
	Title       string                  `json:"title"`
	Context     string                  `json:"context"`
	Topics      string                  `json:"topics"`
	User        string                  `json:"user"`
	Host        string                  `json:"host"`
}

/*查询返回的结果json格式
*{
*    "took": 97,
*    "timed_out": false,
*    "_shards": {
*        "total": 5,
*        "successful": 5,
*        "failed": 0
*    },
*    "hits": {
*        "total": 5,
*        "max_score": 0.05006615,
*        "hits": [
*            {
*                "_index": "milnk_index",
*                "_type": "milnk_type",
*                "_id": "10004",
*                "_score": 0.05006615,
*                "_source": {
*                    "title": "新增新增新增新增",
*                    "context": "",
*                    "topics": "",
*                    "user": "zengshmin",
*                    "host": "127.0.0.1:8080"
*                }
*            },
*            {
*                ....
*            },
*            {
*                ....
*            }
*		]
*}
*}
*/
//搜索link
func (ls *LinkSearch) SearchLink(term string, page int, pagesize int) (*SearchResult, error) {
	if ls.Url == "" {
		ls.Url = LinkSearchUrl
	}
	page, pagesize = PageCheck(page, pagesize)
	data := fmt.Sprintf(`{
						  "query": {
							"multi_match": {
							  "use_dis_max": false,
							  "query": "%s",
							  "fields": [
								"title",
								"context",
								"topics",
								"user",
								"host"
							  ]
							}
						  },
						  "from": %d,
						  "size": %d,
						  "sort": [
							{
							  "_score": "desc"
							}
						  ]
						}`, term, page * pagesize, pagesize)
	resp, err := http.DefaultClient.Post(fmt.Sprintf("%s/_search", ls.Url), "application/json", bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, err
	}

	jData := json.NewDecoder(resp.Body)
	sResult := &SearchResult{}
	err = jData.Decode(sResult)
	if err != nil {
		return nil, err
	}

	return sResult, nil
}






