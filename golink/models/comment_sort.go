package models

import (
    //"bytes"
    //"errors"
    "fmt"
    "github.com/QLeelulu/goku"
    //"github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink"
    //"github.com/QLeelulu/ohlala/golink/utils"
    //"html/template"
    "time"
    "strings"
    "database/sql"
    "strconv"
)


type CommentNode struct {
    Id            int64
    LinkId        int64
    UserId        int64
    Status        int // 评论状态：1代表正常、2代表删除
    Content       string
    ParentId      int64
    Deep          int
    TopParentId   int64
    ParentPath    string
    ChildrenCount int
    VoteUp        int
    VoteDown      int
    RedditScore   float64
    CreateTime    time.Time
	UserName      string

	Children      []*CommentNode

    //user *User `db:"exclude"`
}

//exceptIds:被点击的loadmore(x)所在的层级已经显示的id列表，如果为空字符代表第一次获取评论 
/* parentPath:被点击的loadmore(x)所在的层级的parent_path， 
* 根节点的parent_path="" 
* 第二级的parent_path="父节点id" 
* 第三级的parent_path="第一级父节点id|第二级父节点id" 
*/ 
// topId:评论根节点id，加他过滤缩小范围，提升速度 
func GetComments(exceptIds string, parentPath string, topId int64, linkId int64, sortField string) string { 
    arrExceptIds := strings.Split(exceptIds, ",") 
    //检查每个都是整数才能往后执行
	for _, id := range arrExceptIds { 
		_, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return ""
		}
	} 
    arrParentPath := strings.Split(strings.Trim(parentPath, "/"), "/") 
    //检查每个都是整数才能往后执行,通过arrParentPath.len知道当前loadmore第几级
	for _, id := range arrParentPath { 
		_, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return ""
		}
	}  

    level := len(arrParentPath)
    
	var db *goku.MysqlDB = GetDB()
    defer db.Close()

	where := " c.link_id=? " 
    if level == 0 { //根级别的loadmore 
		if exceptIds != "" { 
			where += fmt.Sprintf("AND c.top_parent_id NOT IN(%s)", exceptIds)
		} 
		sql := fmt.Sprintf("SELECT c.`id`,c.`link_id`,c.`user_id`,c.`parent_path`,c.`children_count`,c.`top_parent_id`,c.`parent_id`,c.`deep`,c.`status`,c.`content`,c.`create_time`,c.`vote_up`,c.`vote_down`,c.`reddit_score`,u.name AS user_name FROM comment c INNER JOIN `user` u ON %s AND c.user_id=u.id order by %s LIMIT 0,%v", where, sortField, golink.MaxCommentCount)
		rows, err := db.Query(sql, linkId) 
		if err == nil {
			link, errLink := Link_GetById(linkId)
			if errLink == nil {
				return BuildCommentTree(db, &rows, link.CommentRootCount - len(arrExceptIds), exceptIds, level)
			}
		}
    } else if level > 0 && exceptIds != "" { 
		where += fmt.Sprintf(" AND c.top_parent_id=? AND c.id NOT IN(?) AND c.parent_path like '%s%' ", parentPath)
		for _, id := range arrExceptIds { 
			where += fmt.Sprintf(" AND c.parent_path not like '%s%'", parentPath + id + "/")
		} 
		sql := fmt.Sprintf("SELECT c.`id`,c.`link_id`,c.`user_id`,c.`parent_path`,c.`children_count`,c.`top_parent_id`,c.`parent_id`,c.`deep`,c.`status`,c.`content`,c.`create_time`,c.`vote_up`,c.`vote_down`,c.`reddit_score`,u.name AS user_name FROM comment c INNER JOIN `user` u ON %s AND c.user_id=u.id order by %s LIMIT 0,%v", where, sortField, golink.MaxCommentCount)
		rows, err := db.Query(sql, linkId, topId, exceptIds) 
		if err == nil {
			commentId, _ := strconv.ParseInt(arrParentPath[len(arrParentPath)-1], 10, 64)
			pComment, errComment := Comment_GetById(commentId)
			if errComment == nil {
				return BuildCommentTree(db, &rows, pComment.ChildrenCount - len(arrExceptIds), exceptIds, level)
			}
		}
    } 
    
    return ""
}

func BuildCommentTree(db *goku.MysqlDB, rows **sql.Rows, childCount int, exceptIds string, level int) string {
	hashTable := map[int64]*CommentNode{} 
	arrRoots := make([]*CommentNode, 0, 1) //记录根节点的数组，递归他显示即可，无需再排序
	strNeedQueryIds := ""
	for rows.Next() {
		comment := ScanCommentNode(rows)//读出一行
		//oldCommentList := CommentList{comment, []CommentList} //初始化一个树节点
		hashTable[comment.Id] = comment //插入hash表中
		
		parentIds := strings.Split(strings.Trim(comment.ParentPath, "/"), "/")
		pLen := len(parentIds)
		passStep := false
		 //上个父节点的对象
		for i:=pLen-1; i>=level; i-- { //循环父节点id(如果loadmore不是处在根节点，就不需要循环到根节点)
			pid,_ := strconv.ParseInt(parentIds[i], 10, 64) //取出parentIds中的pid
			var pComment *CommentNode
			if hashTable[pid] == nil { //hash表中没父节点记录
				pComment = &CommentNode{pid,0,0,0,"",0,0,0,"",0,0,0,0.0,time.Now(),"",nil} //用pid初始化父节点，这个节点数据需要到数据取
				//pCommentList = CommentList{pComment, []CommentList}
				pComment.Children = make([]*CommentNode, 0, 1)
				hashTable[pid] = pComment //加入hash表中
				pComment.Children = append(pComment.Children, comment)
				//pCommentList.Childs.add(oldCommentList) //把子节点添加进去，这里要保证先加入的排在第一位
				strNeedQueryIds += string(pid) + "," //这个节点数据需要到数据取
			} else {
				pComment = hashTable[pid]
				pComment.Children = append(pComment.Children, comment)
				passStep = true
				break //如果有一个父节点已经被包含在hash中，就代表它的父节点的父节点已经初始化过了
			}
			comment = pComment
		}
		
		if passStep == false {
			arrRoots = append(arrRoots, comment)
		}
		
		if golink.MaxCommentCount < len(hashTable) { //如果达到最大节点就跳出了
			break
		}
	
	}
	
	//从数据库读出未填充Comment的数据
	if strNeedQueryIds != "" {
		rows , _ := db.Query("SELECT * FROM Comment where id in(?)", strNeedQueryIds)
		var id int64
		for rows.Next() {
			rows.Scan(&id)
			comment := hashTable[id] //读出一行
			CopyCommentNode(&rows, comment)
		}
	}
	//递归构建html
	return BuildHtmlString(&arrRoots, childCount, exceptIds)
}


func BuildHtmlString(arrRoots *[]*CommentNode, childCount int, exceptIds string) string {
	html := ""
	//parentPath := ""

	return html
}

func ScanCommentNode(rows **sql.Rows) *CommentNode {

	comment := CommentNode{}
	
	rows.Scan(&(comment.Id))
	rows.Scan(&(comment.LinkId))
	rows.Scan(&(comment.UserId))
	rows.Scan(&(comment.ParentPath))
	rows.Scan(&(comment.ChildrenCount))
	rows.Scan(&(comment.TopParentId))
	rows.Scan(&(comment.ParentId))
	rows.Scan(&(comment.Deep))
	rows.Scan(&(comment.Status))
	rows.Scan(&(comment.Content))
	rows.Scan(&(comment.CreateTime))
	rows.Scan(&(comment.VoteUp))
	rows.Scan(&(comment.VoteDown))
	rows.Scan(&(comment.RedditScore))
	rows.Scan(&(comment.UserName))
	
	return &comment
}

func CopyCommentNode(rows **sql.Rows, comment *CommentNode) {
	
	//rows.Scan(&comment.Id)
	rows.Scan(&comment.LinkId)
	rows.Scan(&comment.UserId)
	rows.Scan(&comment.ParentPath)
	rows.Scan(&comment.ChildrenCount)
	rows.Scan(&comment.TopParentId)
	rows.Scan(&comment.ParentId)
	rows.Scan(&comment.Deep)
	rows.Scan(&comment.Status)
	rows.Scan(&comment.Content)
	rows.Scan(&comment.CreateTime)
	rows.Scan(&comment.VoteUp)
	rows.Scan(&comment.VoteDown)
	rows.Scan(&comment.RedditScore)
	rows.Scan(&comment.UserName)
}




