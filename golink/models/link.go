package models

import (
    "database/sql"
    "fmt"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/utils"
    // "net/url"
    "strconv"
    "strings"
    "time"
)

type Link struct {
    Id               int64
    UserId           int64
    Title            string
    Context          string // 如为链接，则为url地址
    ContextType      int    // 0: url, 1:文本
    Topics           string // 话题，用英文逗号分隔
    VoteUp           int64
    VoteDown         int64
    RedditScore      float64
    ViewCount        int // 评论页面查看次数
    ClickCount       int // 链接点击次数
    CommentCount     int // 评论数
    CommentRootCount int // 第一级的评论数
    Status           int // 状态， 2:删除
    CreateTime       time.Time

    user *User `db:"exclude"`
}

// 发布该link的用户
func (l Link) User() *User {
    if l.user == nil {
        l.user = User_GetById(l.UserId)
    }
    return l.user
}

// link是否已经删除掉
func (l Link) Deleted() bool {
    return l.Status == 2
}

// link是否为url
func (l Link) IsUrl() bool {
    return l.ContextType == 0
}

// link的host
func (l Link) Host() string {
    if !l.IsUrl() {
        return ""
    }

    //u, err := url.Parse(l.Context)
    //if err != nil {
    //    return ""
    //}
    //if strings.Index(u.Host, "www.") == 0 {
    //    return u.Host[4:]
    //}
    //return u.Host
    return utils.GetUrlHost(l.Context)
}

// 投票得分
func (l Link) VoteScore() int64 {
    return l.VoteUp - l.VoteDown
}

func (l Link) TopicList() []string {
    if l.Topics == "" {
        return nil
    }
    return strings.Split(l.Topics, ",")
}

func (l Link) SinceTime() string {
    return utils.SmcTimeSince(l.CreateTime)
}

// 顶的百分比。
// 20%则返回值为 20
func (l Link) VoteUpPrec() int {
    prec := float64(l.VoteUp) / float64(l.VoteUp+l.VoteDown) * 100
    return int(prec)
}

// 给view用的link数据
type VLink struct {
    Link
    VoteUped, VoteDowned bool // 是否已顶/踩
    Favorited            bool // 是否已收藏
    SharedByMe           bool // 是否由登陆的用户分享的
}

// 转换为用于view显示用的实例
func Link_ToVLink(links []Link, ctx *goku.HttpContext) []VLink {
    if links == nil || len(links) < 1 {
        return nil
    }
    var userId int64
    if user, ok := ctx.Data["user"].(*User); ok && user != nil {
        userId = user.Id
    }
    l := len(links)

    vlinks := make([]VLink, l, l)
    uids := make([]string, l, l)
    lids := make([]string, l, l)
    lindex := make(map[int64]*VLink)
    for i, link := range links {
        uids[i] = strconv.FormatInt(link.UserId, 10)
        lids[i] = strconv.FormatInt(link.Id, 10)
        vlinks[i] = VLink{Link: link}
        lindex[link.Id] = &(vlinks[i])
    }

    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    // 添加用户信息
    userIndex := make(map[int64]*User)
    qi := goku.SqlQueryInfo{}
    qi.Where = fmt.Sprintf("`id` in (%v)", strings.Join(uids, ","))
    var users []User
    err := db.GetStructs(&users, qi)
    if err != nil {
        goku.Logger().Errorln(err.Error())
    } else if users != nil {
        for i, _ := range users {
            user := &users[i]
            userIndex[user.Id] = user
        }
    }
    for i, _ := range vlinks {
        link := &vlinks[i]
        if user, ok := userIndex[link.UserId]; ok {
            link.user = user
            if user.Id == userId {
                link.SharedByMe = true
            }
        }
    }
    // 添加投票信息
    if userId > 0 {
        qi = goku.SqlQueryInfo{}
        qi.Where = fmt.Sprintf("`user_id`=%v AND `link_id` in (%v)", userId, strings.Join(lids, ","))
        var srs []LinkSupportRecord
        err = db.GetStructs(&srs, qi)
        if err != nil {
            goku.Logger().Errorln(err.Error())
        } else if srs != nil {
            for _, sr := range srs {
                if sr.Score == 1 {
                    lindex[sr.LinkId].VoteUped = true
                } else if sr.Score == -1 {
                    lindex[sr.LinkId].VoteDowned = true
                }
            }
        }

        // 添加收藏信息
        qi.Fields = "link_id"
        rows, err := db.Select("user_favorite_link", qi)
        if err != nil {
            goku.Logger().Errorln(err.Error())
        } else {
            var linkId int64
            for rows.Next() {
                err = rows.Scan(&linkId)
                if err != nil {
                    goku.Logger().Errorln(err.Error())
                    continue
                }
                lindex[linkId].Favorited = true
            }
        }
    }

    return vlinks
}

// 保存link到数据库，如果成功，则返回link的id
func Link_SaveMap(m map[string]interface{}) int64 {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    m["create_time"] = time.Now()
    //新增link默认投票1次,显示的时候默认减一
    m["vote_up"] = 0 //1
    m["reddit_score"] = utils.LinkSortAlgorithm(m["create_time"].(time.Time), int64(0), int64(0))
    m["context_md5"] = utils.MD5_16(strings.ToLower(m["context"].(string)))

    r, err := db.Insert("link", m)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return 0
    }
    var id int64
    id, err = r.LastInsertId()
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return 0
    }

    if id > 0 {
        uid := m["user_id"].(int64)
        // 更新用户的链接计数
        IncCountById(db, "user", uid, "link_count", 1)
        // 直接推送给自己，自己必须看到
        LinkForUser_Add(uid, id, LinkForUser_ByUser)

        // 存入`tui_link_for_handle` 链接处理队列表
        db.Query("INSERT ignore INTO tui_link_for_handle(link_id,create_time,user_id,insert_time,data_type) VALUES (?, ?, ?, NOW(), ?)",
            id, m["create_time"].(time.Time), uid, 1)

        redisClient := GetRedis()
        defer redisClient.Quit()
        // 加入推送队列
        // 格式: pushtype,userid,linkid,timestamp
        qv := fmt.Sprintf("%v,%v,%v,%v", LinkForUser_ByUser, uid, id, time.Now().Unix())
        _, err = redisClient.Lpush(golink.KEY_LIST_PUSH_TO_USER, qv)
        if err != nil {
            goku.Logger().Errorln(err.Error())
            // return 0
        }

    }

    return id
}

// 如果保存失败，则返回错误信息
// 返回为 success, linkId, errors.
// 如果success为false并且linkId大于0，则为提交的url已经存在.
func Link_SaveForm(f *form.Form, userId int64, resubmit bool) (bool, int64, []string, map[string]interface{}) {
    var id int64
    var m map[string]interface{}
    errorMsgs := make([]string, 0)
    if f.Valid() {
        m = f.CleanValues()
        if !resubmit {
            link, err := Link_GetByUrl(m["context"].(string))
            if err == nil && link != nil && link.Id > 0 {
                errorMsgs = append(errorMsgs, "Url已经提交过")
                return false, link.Id, errorMsgs, nil
            }
        }
        m["topics"] = buildTopics(m["topics"].(string))
        m["user_id"] = userId

        id = Link_SaveMap(m)
        if id > 0 {
            Topic_SaveTopics(m["topics"].(string), id)
        } else {
            errorMsgs = append(errorMsgs, golink.ERROR_DATABASE)
        }
    } else {
        errs := f.Errors()
        for _, v := range errs {
            errorMsgs = append(errorMsgs, v[0]+": "+v[1])
        }
    }
    if len(errorMsgs) < 1 {
        return true, id, nil, m
    }
    return false, id, errorMsgs, nil
}

// topic可以用英文逗号或者空格分隔
// 过滤重复topic，最终返回的topic列表只用英文逗号分隔
func buildTopics(topics string) string {
    if topics == "" {
        return ""
    }
    m := make(map[string]string)
    t := strings.Split(topics, ",")
    for _, topic := range t {
        topic = strings.TrimSpace(topic)
        if topic != "" {
            t2 := strings.Split(topic, " ")
            for _, topic2 := range t2 {
                topic2 = strings.TrimSpace(topic2)
                if topic2 != "" {
                    m[strings.ToLower(topic2)] = topic2
                }
            }
        }
    }
    r := ""
    for _, v := range m {
        if r != "" {
            r += ","
        }
        r += v
    }
    return r
}

func Link_GetById(id int64) (*Link, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    l := new(Link)
    err := db.GetStruct(l, "id=?", id)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil, err
    }
    if l.Id > 0 {
        return l, nil
    }
    return nil, nil
}

func Link_GetByIds(ids []int64) ([]Link, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    links := []Link{}
    qi := goku.SqlQueryInfo{}
    sids := ""
    for _, id := range ids {
        sids += "," + strconv.FormatInt(id, 10)
    }
    qi.Where = "id in (" + sids[1:] + ")"
    // qi.Params = []interface{}{ids}
    err := db.GetStructs(&links, qi)
    if err != nil {
        goku.Logger().Errorln("Link_GetByIds error:", err.Error())
        return nil, err
    }
    return links, nil
}

// url不区分大小写
func Link_GetByUrl(url string) (*Link, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    l := new(Link)
    urlMd5 := utils.MD5_16(strings.ToLower(url))
    err := db.GetStruct(l, "context_md5=? and `status`<>2 order by comment_count desc", urlMd5)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil, err
    }
    if l.Id > 0 {
        return l, nil
    }
    return nil, nil
}

// @page: 从1开始
// @return: topics, total-count, err
func Link_GetByPage(page, pagesize int, order string) ([]Link, int64, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    page, pagesize = utils.PageCheck(page, pagesize)

    qi := goku.SqlQueryInfo{}
    qi.Limit = pagesize
    qi.Offset = page * pagesize
    if order == "" {
        qi.Order = "id desc"
    } else {
        qi.Order = order
    }
    var links []Link
    err := db.GetStructs(&links, qi)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil, 0, err
    }

    total, err := db.Count("link", "")
    if err != nil {
        goku.Logger().Errorln(err.Error())
    }
    return links, total, nil
}

// 获取由用户发布的link
// @page: 从1开始
func Link_ByUser(userId int64, page, pagesize int) []Link {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    page, pagesize = utils.PageCheck(page, pagesize)

    qi := goku.SqlQueryInfo{}
    qi.Limit = pagesize
    qi.Offset = page * pagesize
    qi.Where = "`user_id`=? and `status`=0"
    qi.Params = []interface{}{userId}
    qi.Order = "id desc"
    var links []Link
    err := db.GetStructs(&links, qi)
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil
    }
    return links
}

// 获取属于某话题的link
// @page: 从1开始
func Link_ForTopic(topicId int64, page, pagesize int, sortType string, t string) ([]Link, error) {
    var db *goku.MysqlDB = GetDB()
    db.Debug = true
    defer db.Close()

    page, pagesize = utils.PageCheck(page, pagesize)

    sortField := "tl.reddit_score DESC,tl.link_id DESC"
    tableName := "tui_link_for_topic_top"
    switch {
    case sortType == golink.ORDER_TYPE_HOTC: //热议
        sortField = "l.comment_count DESC,tl.link_id DESC"
        tableName = "tui_link_for_topic_top"
    case sortType == golink.ORDER_TYPE_CTVL: //争议
        sortField = "tl.vote_abs_score ASC,tl.vote_add_score DESC,tl.link_id DESC"
        tableName = "tui_link_for_topic_hot"
    case sortType == golink.ORDER_TYPE_TIME: //最新
        sortField = "tl.link_id desc"
        tableName = "tui_link_for_topic_later"
    case sortType == golink.ORDER_TYPE_VOTE: //得分
        sortField = "tl.vote DESC, tl.link_id DESC"
        tableName = "tui_link_for_topic_vote"
    default: //热门
        sortField = "tl.reddit_score DESC,tl.link_id DESC"
        tableName = "tui_link_for_topic_top"
    }

    qi := goku.SqlQueryInfo{}
    qi.Fields = "l.id, l.user_id, l.title, l.context, l.topics, l.vote_up, l.vote_down, l.view_count, l.click_count, l.comment_count, l.create_time"
    qi.Join = " tl INNER JOIN `link` l ON tl.link_id=l.id"

    if sortType == golink.ORDER_TYPE_CTVL || sortType == golink.ORDER_TYPE_VOTE {
        qi.Where = "tl.topic_id=? AND tl.time_type=?"
        switch {
        case t == "all": //1:全部时间；2:这个小时；3:今天；4:这周；5:这个月；6:今年
            qi.Params = []interface{}{topicId, 1}
        case t == "hour":
            qi.Params = []interface{}{topicId, 2}
        case t == "day":
            qi.Params = []interface{}{topicId, 3}
        case t == "week":
            qi.Params = []interface{}{topicId, 4}
        case t == "month":
            qi.Params = []interface{}{topicId, 5}
        case t == "year":
            qi.Params = []interface{}{topicId, 6}
        default:
            qi.Params = []interface{}{topicId, 1}
        }
    } else {
        qi.Where = "tl.topic_id=?"
        qi.Params = []interface{}{topicId}
    }
    qi.Limit = pagesize
    qi.Offset = pagesize * page
    qi.Order = sortField

    rows, err := db.Select(tableName, qi)

    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil, err
    }
    defer rows.Close()

    links := make([]Link, 0)
    for rows.Next() {
        link := Link{}
        err = rows.Scan(&link.Id, &link.UserId, &link.Title, &link.Context, &link.Topics,
            &link.VoteUp, &link.VoteDown, &link.ViewCount, &link.ClickCount, &link.CommentCount, &link.CreateTime)
        if err != nil {
            goku.Logger().Errorln(err.Error())
            return nil, err
        }
        links = append(links, link)
    }
    return links, nil
}

// 获取属于某用户的link
// @page: 从1开始
// @orderType: 排序类型, hot:热门, hotc:热议, time:最新, vote:投票得分, ctvl:争议
func Link_ForUser(userId int64, orderType string, page, pagesize int) ([]Link, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    page, pagesize = utils.PageCheck(page, pagesize)

    qi := goku.SqlQueryInfo{}
    qi.Fields = "l.id, l.user_id, l.title, l.context, l.topics, l.vote_up, l.vote_down, l.view_count, l.click_count, l.comment_count, l.create_time"
    qi.Join = " ul INNER JOIN `link` l ON ul.link_id=l.id"
    qi.Where = "ul.user_id=?"
    qi.Params = []interface{}{userId}
    qi.Limit = pagesize
    qi.Offset = pagesize * page
    switch orderType {
    case golink.ORDER_TYPE_TIME: // 最新
        qi.Order = "l.id desc"
    case golink.ORDER_TYPE_HOTC: // 热议
        qi.Order = "l.comment_count desc, id desc"
    case golink.ORDER_TYPE_CTVL: // 争议
        qi.Order = "l.dispute_score desc, id desc"
    case golink.ORDER_TYPE_VOTE: // 得分
        qi.Order = "l.vote_up desc, id desc"
    default:
        qi.Order = "l.reddit_score desc, id desc"
    }

    rows, err := db.Select(LinkForUser_TableName(userId), qi)

    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil, err
    }
    defer rows.Close()

    links := make([]Link, 0)
    for rows.Next() {
        link := Link{}
        err = rows.Scan(&link.Id, &link.UserId, &link.Title, &link.Context, &link.Topics,
            &link.VoteUp, &link.VoteDown, &link.ViewCount, &link.ClickCount, &link.CommentCount, &link.CreateTime)
        if err != nil {
            goku.Logger().Errorln(err.Error())
            return nil, err
        }
        links = append(links, link)
    }
    return links, nil
}

// 更新链接的评论查看计数
func Link_IncViewCount(linkId int64, inc int) (sql.Result, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    return IncCountById(db, "link", linkId, "view_count", 1)
}

// 更新链接的点击数
func Link_IncClickCount(linkId int64, inc int) (sql.Result, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    return IncCountById(db, "link", linkId, "click_count", 1)
}

// 根据id列表获取link
func Link_GetByIdList(searchItems []utils.SearchHitItem) ([]Link, error) {
    hashTable := map[int64]*Link{}
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    var strLinkIdList string
    for _, item := range searchItems {
        strLinkIdList += item.Id + ","
    }
    strLinkIdList += "0"

    qi := goku.SqlQueryInfo{}
    qi.Fields = "id, user_id, title, context, topics, vote_up, vote_down, view_count, comment_count, create_time, status"
    qi.Where = "id IN(" + strLinkIdList + ")"
    rows, err := db.Select("link", qi)

    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil, err
    }
    defer rows.Close()

    links := make([]Link, 0)
    for rows.Next() {
        link := &Link{}
        err = rows.Scan(&link.Id, &link.UserId, &link.Title, &link.Context, &link.Topics,
            &link.VoteUp, &link.VoteDown, &link.ViewCount, &link.CommentCount, &link.CreateTime, &link.Status)
        if err != nil {
            goku.Logger().Errorln(err.Error())
            return nil, err
        }
        hashTable[link.Id] = link
    }
    for _, item := range searchItems {
        linkId, err := strconv.ParseInt(item.Id, 10, 64)
        link := hashTable[linkId]
        if err == nil && link != nil {
            links = append(links, *link)
        }
    }

    return links, nil
}
