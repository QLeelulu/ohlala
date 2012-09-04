package models

import (
    "bytes"
    //"errors"
    "fmt"
    "github.com/QLeelulu/goku"
    //"github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/utils"
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

func (c CommentNode) SinceTime() string {
    return utils.SmcTimeSince(c.CreateTime)
}

func (cl CommentNode) renderItemBegin(b *bytes.Buffer, sortType string) {

    b.WriteString(fmt.Sprintf(`<div class="cm" data-id="%v">
<div class="vt">
 <a class="icon-thumbs-up up" href="javascript:"></a>
 <a class="icon-thumbs-down down" href="javascript:"></a>
</div>
<div class="ct">
 <div class="uif">
   <a class="ep" href="javascript:">[–]</a>
   <a href="/user/%v">%v</a>
   <i class="v" title="↑%v ↓%v">%v分</i> <i class="t">%v</i>
 </div>
 <div class="tx">%v</div>
 <div class="ed">
   <a href="/link/permacoment/%v/%v/%s/" class="rp">查看</a>
   <a href="javascript:" class="rp">回复</a>
 </div>`, cl.Id,
        cl.UserId, cl.UserName,
        cl.VoteUp, cl.VoteDown,
        cl.VoteUp-cl.VoteDown,
        cl.SinceTime(), cl.Content, cl.LinkId, cl.Id, sortType))
}
func (cl CommentNode) renderItemEnd(b *bytes.Buffer) {

    b.WriteString(`</div></div>`)
}

func (comment *CommentNode) Copy(temp *CommentNode) {
	comment.Id = temp.Id
	comment.LinkId = temp.LinkId
	comment.UserId = temp.UserId
	comment.ParentPath = temp.ParentPath
	comment.ChildrenCount = temp.ChildrenCount
	comment.TopParentId = temp.TopParentId
	comment.ParentId = temp.ParentId
	comment.Deep = temp.Deep
	comment.Status = temp.Status
	comment.Content = temp.Content
	comment.CreateTime = temp.CreateTime
	comment.VoteUp = temp.VoteUp
	comment.VoteDown = temp.VoteDown
	comment.RedditScore = temp.RedditScore
	comment.UserName = temp.UserName
}

func GetPermalinkComment(linkId int64, commentId int64, sortType string) string {
	strFilter := ""
	topId := int64(0)
	comment, err := Comment_GetById(commentId)
	if err == nil && comment != nil {
		topId = comment.TopParentId
		if topId == 0 { //根节点
			topId = commentId
			strFilter = fmt.Sprintf("AND (c.top_parent_id=%d OR c.Id=%v)", topId, commentId)
		} else {
			strFilter = fmt.Sprintf("AND (c.top_parent_id=%d AND c.parent_path like '%s%d/%s' OR c.Id=%v)", topId, comment.ParentPath, commentId, "%", commentId)
		}
		return GetSortComments("", comment.ParentPath, 0, linkId, sortType, strFilter, false)
	}
	return ""
}

//exceptIds:被点击的loadmore(x)所在的层级已经显示的id列表，如果为空字符代表第一次获取评论 
/* parentPath:被点击的loadmore(x)所在的层级的parent_path， 
* 根节点的parent_path="" 
* 第二级的parent_path="父节点id" 
* 第三级的parent_path="第一级父节点id|第二级父节点id" 
*/ 
// topId:评论根节点id，加他过滤缩小范围，提升速度 
// sortType:"top":热门；"hot":热议；"later":最新；"vote":得分
func GetSortComments(exceptIds string, parentPath string, topId int64, linkId int64, sortType string, permaFilter string, isLoadMore bool) string { 
	var arrExceptIds []string
	if exceptIds != "" {
		arrExceptIds = strings.Split(exceptIds, ",") 
		//检查每个都是整数才能往后执行
		for _, id := range arrExceptIds { 
			_, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return ""
			}
		} 
	}
	
	pId := int64(0)
	var arrParentPath []string
	if parentPath != "/" {
		arrParentPath = strings.Split(strings.Trim(parentPath, "/"), "/") 
		//检查每个都是整数才能往后执行,通过arrParentPath.len知道当前loadmore第几级
		for _, id := range arrParentPath { 
			id, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return ""
			}
			pId = id
		}
	}

	sortField := "c.reddit_score DESC,c.id DESC"
	switch {
		case sortType == "top": //热门
		    sortField = "c.reddit_score DESC,c.id DESC"
		case sortType == "hot": //热议
		    sortField = "ABS(c.vote_up-c.vote_down) ASC,(c.vote_up+c.vote_down) DESC,c.id DESC"
		case sortType == "later": //最新
			sortField = "c.id DESC"
		case sortType == "vote": //得分
			sortField = "(c.vote_up-c.vote_down) DESC, c.id DESC"
    }
	
    level := len(arrParentPath)

	var db *goku.MysqlDB = GetDB()
db.Debug = true
    defer db.Close()

	where := " c.link_id=? " 
	if permaFilter != "" { //显示某个评论
		sql := fmt.Sprintf("SELECT c.`id`,c.`link_id`,c.`user_id`,c.`parent_path`,c.`children_count`,c.`top_parent_id`,c.`parent_id`,c.`deep`,c.`status`,c.`content`,c.`create_time`,c.`vote_up`,c.`vote_down`,c.`reddit_score`,u.name AS user_name FROM comment c INNER JOIN `user` u ON %s %s AND c.user_id=u.id order by %s LIMIT 0,%v", where, permaFilter, sortField, golink.MaxCommentCount)

		rows, err := db.Query(sql, linkId) 
		if err == nil {
			return BuildCommentTree(db, &rows, 1, exceptIds, level, parentPath, pId, sortType, isLoadMore)
		}

	} else {
		if level == 0 { //根级别的loadmore 
			if exceptIds != "" { 
				where += fmt.Sprintf("AND c.top_parent_id NOT IN(%s) AND c.Id NOT IN(%s)", exceptIds, exceptIds)
			} 
			sql := fmt.Sprintf("SELECT c.`id`,c.`link_id`,c.`user_id`,c.`parent_path`,c.`children_count`,c.`top_parent_id`,c.`parent_id`,c.`deep`,c.`status`,c.`content`,c.`create_time`,c.`vote_up`,c.`vote_down`,c.`reddit_score`,u.name AS user_name FROM comment c INNER JOIN `user` u ON %s AND c.user_id=u.id order by %s LIMIT 0,%v", where, sortField, golink.MaxCommentCount)
			rows, err := db.Query(sql, linkId) 
			if err == nil {
				link, errLink := Link_GetById(linkId)
				if errLink == nil {
					return BuildCommentTree(db, &rows, link.CommentRootCount - len(arrExceptIds), exceptIds, level, parentPath, pId, sortType, isLoadMore)
				}
			}
		} else if level > 0 {
			if exceptIds != "" {
				where += fmt.Sprintf(" AND c.top_parent_id=? AND c.id NOT IN(%s) AND c.parent_path like '%s%s' ", exceptIds, parentPath, "%")
			} else {
				where += fmt.Sprintf(" AND c.top_parent_id=? AND c.parent_path like '%s%s' ", parentPath, "%")
			}
			for _, id := range arrExceptIds { 
				where += fmt.Sprintf(" AND c.parent_path not like '%s%s/%s'", parentPath, id, "%")
			} 
			sql := fmt.Sprintf("SELECT c.`id`,c.`link_id`,c.`user_id`,c.`parent_path`,c.`children_count`,c.`top_parent_id`,c.`parent_id`,c.`deep`,c.`status`,c.`content`,c.`create_time`,c.`vote_up`,c.`vote_down`,c.`reddit_score`,u.name AS user_name FROM comment c INNER JOIN `user` u ON %s AND c.user_id=u.id order by %s LIMIT 0,%v", where, sortField, golink.MaxCommentCount)

			rows, err := db.Query(sql, linkId, topId) 
			if err == nil {
				commentId, _ := strconv.ParseInt(arrParentPath[level-1], 10, 64)
				pComment, errComment := Comment_GetById(commentId)
				if errComment == nil {
					return BuildCommentTree(db, &rows, pComment.ChildrenCount - len(arrExceptIds), exceptIds, level, parentPath, pId, sortType, isLoadMore)
				}
			}
		} 
    }

    return ""
}

func BuildCommentTree(db *goku.MysqlDB, rows **sql.Rows, childCount int, exceptIds string, level int, parentPath string, 
					pId int64, sortType string, isLoadMore bool) string {
	hashTable := map[int64]*CommentNode{}
	hashAppend := map[int64]*CommentNode{}

	var arrRows []int64
	hashRows := map[int64]*CommentNode{}

	hashRoot := map[int64]*CommentNode{}
	arrRoots := make([]*CommentNode, 0) //记录根节点的数组，递归他显示即可，无需再排序 
	strNeedQueryIds := ""

	for rows.Next() {
		comment := ScanCommentNode(rows)//读出一行
		hashRows[comment.Id] = comment
		arrRows = append(arrRows, comment.Id)
	}
		
	for _, item := range arrRows {
		comment := hashRows[item]//读出一行
		hashTable[comment.Id] = comment //插入hash表中

		if hashRoot[comment.Id] == nil && comment.ParentPath == parentPath {
			hashRoot[comment.Id] = comment
			arrRoots = append(arrRoots, comment)
		}
		var parentIds []string
		if comment.ParentPath != "/" {
			parentIds = strings.Split(strings.Trim(comment.ParentPath, "/"), "/")
		}
		pLen := len(parentIds)
		//passStep := false
		//hasFor := false
		 //上个父节点的对象
		for i:=pLen-1; i>=level; i-- { //循环父节点id(如果loadmore不是处在根节点，就不需要循环到根节点)
			//hasFor = true
			pid,_ := strconv.ParseInt(parentIds[i], 10, 64) //取出parentIds中的pid
			var pComment *CommentNode
			pComment = hashTable[pid]
			if pComment == nil { //hash表中没父节点记录
				pComment = hashRows[pid]
				if pComment == nil {
					pComment = &CommentNode{}//&CommentNode{pid,0,0,0,"",0,0,0,"",0,0,0,0.0,time.Now(),"",nil} //用pid初始化父节点，这个节点数据需要到数据取
					pComment.Id = pid
					strNeedQueryIds += fmt.Sprintf("%d,", pid) //这个节点数据需要到数据取
				}
				pComment.Children = make([]*CommentNode, 0)
				pComment.Children = append(pComment.Children, comment)
				hashTable[pid] = pComment //加入hash表中
			} else {
				if hashAppend[comment.Id] == nil {
					pComment.Children = append(pComment.Children, comment)
					hashAppend[comment.Id] = comment
				}
				//passStep = true
				break //如果有一个父节点已经被包含在hash中，就代表它的父节点的父节点已经初始化过了
			}
			comment = pComment
				
			if hashRoot[pComment.Id] == nil && i == level { //hasFor == true && passStep == false
				hashRoot[pComment.Id] = pComment
				arrRoots = append(arrRoots, pComment)
			}
		}
		
		
		if golink.MaxCommentCount < len(hashTable) { //如果达到最大节点就跳出了
			break
		}
	
	}

	//从数据库读出未填充Comment的数据
	if strNeedQueryIds != "" {
		strNeedQueryIds = strings.TrimRight(strNeedQueryIds, ",")
		rows , _ := db.Query(fmt.Sprintf("SELECT c.`id`,c.`link_id`,c.`user_id`,c.`parent_path`,c.`children_count`,c.`top_parent_id`,c.`parent_id`,c.`deep`,c.`status`,c.`content`,c.`create_time`,c.`vote_up`,c.`vote_down`,c.`reddit_score`,u.name AS user_name FROM comment c INNER JOIN `user` u ON c.user_id=u.id AND c.Id IN(%s)", strNeedQueryIds))

		for rows.Next() {
			temp := ScanCommentNode(&rows)
			comment := hashTable[temp.Id] //读出一行

			comment.Copy(temp)
		}
	}

    var b bytes.Buffer
	BuildHtmlString(&arrRoots, childCount, exceptIds, &b, pId, false, sortType, isLoadMore)
	return b.String()
}


func BuildHtmlString(arrRoots *[]*CommentNode, childCount int, exceptIds string, b *bytes.Buffer, pId int64, 
					loadLine bool, sortType string, isLoadMore bool) {
	
    if arrRoots == nil || len(*arrRoots) == 0 {
        return
    }

	parentPath := ""
	topId := int64(0)
	linkId := int64(0)
	
	if isLoadMore == false {
		if loadLine {
			b.WriteString(fmt.Sprintf(`<div class="cd" pid="pid%d">`, pId))
		} else {
			b.WriteString(fmt.Sprintf(`<div pid="pid%d">`, pId))
		}
	}

    for _, item := range *arrRoots {

		item.renderItemBegin(b, sortType)
		BuildHtmlString(&item.Children, item.ChildrenCount, "", b, item.Id, true, sortType, isLoadMore)
		if item.ChildrenCount > 0 && len(item.Children) == 0 {
			topId = item.TopParentId
			if topId == 0 {
				topId = item.Id
			}
			b.WriteString(fmt.Sprintf("<div class='cm' pid='pid%d'><div class='fucklulu' lmid='lm%d' ><a href='javascript:' pId='%d' exIds='' pp='%s%d/' tId='%d' lId='%d' srt='%s'>追载(%d)</a></div></div>", 
			item.Id, item.Id, item.Id, item.ParentPath, item.Id, topId, item.LinkId, sortType, item.ChildrenCount))
		}
		item.renderItemEnd(b)
		

		exceptIds += fmt.Sprintf("%v,", item.Id)
		parentPath = item.ParentPath
		topId = item.TopParentId
		linkId = item.LinkId
    }

	//构建loadmore标签，exceptIds是下次点击loadmore是返回给服务器告诉它已经显示过这些，需要排除它们
	rLen := len(*arrRoots)
	if childCount > rLen {
		b.WriteString(fmt.Sprintf("<div class='fucklulu' lmid='lm%d' ><a href='javascript:' pId='%d' exIds='%s' pp='%s' tId='%d' lId='%d' srt='%s'>追载(%d)</a></div>", 
			pId, pId, strings.TrimRight(exceptIds, ","), parentPath, topId, linkId, sortType, childCount - rLen))
	}

	if isLoadMore == false {
    	b.WriteString(`</div>`)
	}
}

func ScanCommentNode(rows **sql.Rows) *CommentNode {

	comment := CommentNode{}
	
	rows.Scan(&comment.Id, &comment.LinkId, &comment.UserId, &comment.ParentPath, &comment.ChildrenCount, &comment.TopParentId, 
		&comment.ParentId, &comment.Deep, &comment.Status, &comment.Content, &comment.CreateTime, &comment.VoteUp, &comment.VoteDown, 
		&comment.RedditScore, &comment.UserName)
	
	return &comment
}




