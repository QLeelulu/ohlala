package main

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/utils"
    "github.com/QLeelulu/ohlala/golink/models"
    "fmt"
)


func main() {
	//ls := utils.LinkSearch{}
	//ls.SearchLink("名次", 1, 30)
	//return

    var db *goku.MysqlDB = models.GetDB()
	defer db.Close()
    rows, _ := db.Query("SELECT l.id,l.title,l.context,l.topics,u.Name,l.context_type FROM `link` l INNER JOIN `user` u ON l.user_id=u.id") 
	var linkId int64
	var title string
	var context string
	var topics string
	var userName string
	var contextType int
    for rows.Next() {
        rows.Scan(&linkId, &title, &context, &topics, &userName, &contextType)
        addLinkForSearch(contextType, linkId, title, context, topics, userName)
    }

    fmt.Println("执行完成")

}

//添加link到es搜索; contextType:0: url, 1:文本
func addLinkForSearch(contextType int, linkId int64, title string, context string, topics string, userName string) {
	m := make(map[string]interface{})
	m["id"] = linkId
	m["title"] = title
	m["context"] = context
	m["topics"] = topics
	m["username"] = userName
	if contextType == 0 {
		m["host"] = utils.GetUrlHost(m["context"].(string))
		m["context"] = ""
	} else {
		m["host"] = ""
	}
	ls := utils.LinkSearch{}
	_, err := ls.AddLink(m)
	if err != nil {
		fmt.Println("失败: ", m)
	} else {
		fmt.Println("成功: ", m)
	}

}





















