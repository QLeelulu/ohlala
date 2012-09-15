package controllers

import (
    "fmt"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/models"
    "io"
    "os"
    "path"
    "regexp"
    "strconv"
    "time"
)

/**
 * Controller: topic
 */
var _ = goku.Controller("topic").

    /**
     * 话题列表页
     */
    Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {

    topics, _ := models.Topic_GetTops(1, 30)
    ctx.ViewData["TopTab"] = "topic"
    return ctx.View(models.Topic_ToVTopics(topics, ctx))

}).

    /**
     * 查看话题信息页
     */
    Get("show", func(ctx *goku.HttpContext) goku.ActionResulter {

    ctx.ViewData["TopTab"] = "topic"
    topicName, _ := ctx.RouteData.Params["name"]
    topic, _ := models.Topic_GetByName(topicName)

    if topic == nil {
        ctx.ViewData["errorMsg"] = "话题不存在"
        return ctx.Render("error", nil)
    }

	sort := ctx.Get("srt") //排序方式
	t := ctx.Get("t") //时间范围
	
	ctx.ViewData["Order"] = "top"
	if sort == "top" || sort == "hot" || sort == "later" || sort == "vote" {
		ctx.ViewData["Order"] = sort
	}

    links, _ := models.Link_ForTopic(topic.Id, 1, 20, sort, t)
    followers, _ := models.Topic_GetFollowers(topic.Id, 1, 12)

    ctx.ViewData["Links"] = models.Link_ToVLink(links, ctx)
    ctx.ViewData["Followers"] = followers
    return ctx.View(models.Topic_ToVTopic(topic, ctx))

}).Filters(filters.NewRequireLoginFilter()).

    /**
     * 关注话题
     */
    Post("follow", func(ctx *goku.HttpContext) goku.ActionResulter {

    topicId, _ := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    ok, err := models.Topic_Follow(ctx.Data["user"].(*models.User).Id, topicId)
    var errs string
    if err != nil {
        errs = err.Error()
    }
    r := map[string]interface{}{
        "success": ok,
        "errors":  errs,
    }
    return ctx.Json(r)

}).Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter()).

    /**
     * 取消关注话题
     */
    Post("unfollow", func(ctx *goku.HttpContext) goku.ActionResulter {

    topicId, _ := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    ok, err := models.Topic_UnFollow(ctx.Data["user"].(*models.User).Id, topicId)
    var errs string
    if err != nil {
        errs = err.Error()
    }
    r := map[string]interface{}{
        "success": ok,
        "errors":  errs,
    }
    return ctx.Json(r)

}).Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter()).

    /**
     * 上传话题图片
     */
    Post("upimg", actionUpimg).
    Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter()).

    /**
     * 获取用户信息
     * 用于浮动层
     */
    Get("pbox-info", actionPopupBoxInfo).
    Filters(filters.NewAjaxFilter())

var acceptFileTypes = regexp.MustCompile(`gif|jpeg|jpg|png`)

/**
 * 上传话题图片
 */
func actionUpimg(ctx *goku.HttpContext) goku.ActionResulter {
    var ok = false
    var errs string
    topicId, err := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    if err == nil && topicId > 0 {
        imgFile, header, err2 := ctx.Request.FormFile("topic-image")
        err = err2
        defer func() {
            if imgFile != nil {
                imgFile.Close()
            }
        }()

        if err == nil {
            ext := path.Ext(header.Filename)
            if acceptFileTypes.MatchString(ext[1:]) == false {
                errs = "错误的文件类型"
            } else {
                sid := strconv.FormatInt(topicId, 10)
                saveDir := path.Join(ctx.RootDir(), golink.PATH_IMAGE_AVATAR, "topic", sid[len(sid)-2:])
                err = os.MkdirAll(saveDir, 0755)
                if err == nil {
                    saveName := fmt.Sprintf("%v_%v%v",
                        strconv.FormatInt(topicId, 36),
                        strconv.FormatInt(time.Now().UnixNano(), 36),
                        ext)
                    savePath := path.Join(saveDir, saveName)
                    var f *os.File
                    f, err = os.Create(savePath)
                    defer f.Close()
                    if err == nil {
                        _, err = io.Copy(f, imgFile)
                        if err == nil {
                            // update to db
                            _, err2 := models.Topic_UpdatePic(topicId, path.Join(sid[len(sid)-2:], saveName))
                            err = err2
                            if err == nil {
                                ok = true
                            }
                        }
                    }
                }
            }
        }
    } else if topicId < 1 {
        errs = "参数错误"
    }

    if err != nil {
        errs = err.Error()
    }
    r := map[string]interface{}{
        "success": ok,
        "errors":  errs,
    }
    return ctx.Json(r)
}

/**
 * 获取用户信息
 * 用于浮动层
 */
func actionPopupBoxInfo(ctx *goku.HttpContext) goku.ActionResulter {

    topicName := ctx.Get("t")
    topic, _ := models.Topic_GetByName(topicName)

    if topic != nil {
        return ctx.RenderPartial("pop-info", models.Topic_ToVTopic(topic, ctx))
    }
    return ctx.Html("")
}
