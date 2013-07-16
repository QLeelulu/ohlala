package controllers

import (
    "bytes"
    "fmt"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/forms"
    "github.com/QLeelulu/ohlala/golink/models"
    "github.com/QLeelulu/ohlala/golink/utils"
    "html/template"
    "net/url"
    "strconv"
    "strings"
    "time"
)

var _ = goku.Controller("link").
    /**
     * 查看某评论
     */
    Get("permacoment", link_permacoment).
    /**
     * 查看一个链接的评论
     */
    Get("show", link_show).

    /**
     * 删除link
     */
    Post("ajax-del", link_ajaxDel).Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter()).

    /**
     * 提交一个链接并保存到数据库
     */
    Post("submit", link_submit).Filters(filters.NewRequireLoginFilter()).

    /**
     * 提交评论并保存到数据库
     */
    Post("ajax-comment", link_ajax_comment).Filters(filters.NewRequireLoginFilter(), filters.NewAjaxFilter()).

    /**
     * 提交评论并保存到数据库
     */
    Post("inc-click", link_incClick).Filters(filters.NewAjaxFilter()).

    /**
     * 提交链接的表单页面
     */
    Get("submit", func(ctx *goku.HttpContext) goku.ActionResulter {
    ctx.ViewData["Values"] = map[string]string{
        "title":   ctx.Get("t"),
        "context": ctx.Get("u"),
    }
    return ctx.View(nil)
}).Filters(filters.NewRequireLoginFilter()).
    /**
     * 提交链接的表单页面
     */
    Get("search", link_search).
    Get("searchmorelink", link_search_loadMore).
    Filters(filters.NewAjaxFilter())

//

/**
 * 提交评论并保存到数据库
 */
func link_ajax_comment(ctx *goku.HttpContext) goku.ActionResulter {

    f := forms.NewCommentSubmitForm()
    f.FillByRequest(ctx.Request)

    var success bool
    var errorMsgs, commentHTML string
    var commentId int64
    if ctx.RouteData.Params["id"] != f.Values()["link_id"] {
        errorMsgs = "参数错误"
    } else {
        var errors []string
        user := ctx.Data["user"].(*models.User)
        success, commentId, errors = models.Comment_SaveForm(f, user.Id)
        if errors != nil {
            errorMsgs = strings.Join(errors, "\n")
        } else {
            linkId, _ := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
            m := f.CleanValues()
            cn := models.CommentNode{}
            cn.Id = commentId
            cn.LinkId = linkId
            cn.UserId = user.Id
            cn.Status = 1
            cn.Content = m["content"].(string)
            cn.ParentId = m["parent_id"].(int64)
            cn.ChildrenCount = 0
            cn.VoteUp = 0
            cn.CreateTime = time.Now()
            cn.UserName = user.Name

            sortType := ""
            var b *bytes.Buffer = new(bytes.Buffer)
            cn.RenderSelfOnly(b, sortType)
            commentHTML = b.String()
            //models.GetPermalinkComment(linkId, commentId, "")
        }
    }
    r := map[string]interface{}{
        "success":     success,
        "errors":      errorMsgs,
        "commentHTML": commentHTML,
    }
    return ctx.Json(r)
}

// 增加链接的点击统计数
func link_incClick(ctx *goku.HttpContext) goku.ActionResulter {
    var success bool
    var errorMsgs string
    id := ctx.Get("id")
    if id == "" {
        errorMsgs = "参数错误"
    } else {
        linkId, err := strconv.ParseInt(id, 10, 64)
        if err == nil && linkId > 0 {
            _, err = models.Link_IncClickCount(linkId, 1)
            if err == nil {
                success = true
            }
        }
        if err != nil {
            goku.Logger().Error(err.Error())
            errorMsgs = err.Error()
        }
    }

    r := map[string]interface{}{
        "success": success,
        "errors":  errorMsgs,
    }
    return ctx.Json(r)
}

/**
 * 提交一个链接并保存到数据库
 */
func link_submit(ctx *goku.HttpContext) goku.ActionResulter {

    f := forms.CreateLinkSubmitForm()
    f.FillByRequest(ctx.Request)

    var resubmit bool
    if ctx.Get("resubmit") == "true" {
        resubmit = true
    }
    user := ctx.Data["user"].(*models.User)
    success, linkId, errorMsgs, _ := models.Link_SaveForm(f, user.Id, resubmit)

    if success {
        go addLinkForSearch(0, m, linkId, user.Name) //contextType:0: url, 1:文本   TODO:

        return ctx.Redirect(fmt.Sprintf("/link/%d", linkId))
    } else if linkId > 0 {
        return ctx.Redirect(fmt.Sprintf("/link/%d?already_submitted=true", linkId))
    } else {
        ctx.ViewData["Errors"] = errorMsgs
        ctx.ViewData["Values"] = f.Values()
    }
    return ctx.View(nil)

}

//添加link到es搜索; contextType:0: url, 1:文本
func addLinkForSearch(contextType int, m map[string]interface{}, linkId int64, userName string) {

    m["id"] = linkId
    m["username"] = userName
    if contextType == 0 {
        m["host"] = utils.GetUrlHost(m["context"].(string))
        m["context"] = ""
    } else {
        m["host"] = ""
    }
    ls := utils.LinkSearch{}
    ls.AddLink(m)
}

// 删除link
func link_ajaxDel(ctx *goku.HttpContext) goku.ActionResulter {
    var errs string
    var ok = false

    linkId, err := strconv.ParseInt(ctx.RouteData.Params["id"], 10, 64)
    if err == nil {
        user := ctx.Data["user"].(*models.User)
        link, err := models.Link_GetById(linkId)
        if err == nil {
            // 只可以删除自己的链接
            if link.UserId == user.Id {
                err = models.Link_DelById(linkId)
                if err == nil {
                    ok = true
                }
            } else {
                errs = "不允许的操作"
            }
        }
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

// 查看link
func link_show(ctx *goku.HttpContext) goku.ActionResulter {
    if ctx.Get("already_submitted") == "true" {
        ctx.ViewData["AlreadySubmitted"] = true
    }
    return link_showWithComments(ctx, ctx.RouteData.Params["id"], "0")
}

func link_permacoment(ctx *goku.HttpContext) goku.ActionResulter {
    return link_showWithComments(ctx,
        ctx.RouteData.Params["id"], ctx.RouteData.Params["cid"])
}

var ORDER_NAMES map[string]string = map[string]string{
    "hot":  "最佳",
    "hotc": "热议",
    "time": "最新",
    "vote": "得分",
    "ctvl": "争议",
}

func link_showWithComments(ctx *goku.HttpContext, slinkId, scommentId string) goku.ActionResulter {

    linkId, err1 := strconv.ParseInt(slinkId, 10, 64)
    commentId, err2 := strconv.ParseInt(scommentId, 10, 64)
    if err1 != nil || err2 != nil {
        ctx.ViewData["errorMsg"] = "参数错误"
        return ctx.Render("error", nil)
    }
    link, err := models.Link_GetById(linkId)
    if err != nil {
        ctx.ViewData["errorMsg"] = "服务器开小差了 >_<!!"
        return ctx.Render("error", nil)
    }

    if link == nil {
        ctx.ViewData["errorMsg"] = "内容不存在，去首页逛逛吧"
        return ctx.Render("error", nil)
    }

    if link.Deleted() {
        ctx.ViewData["errorMsg"] = "内容已被摧毁，去首页逛逛吧"
        return ctx.Render("error", nil)
    }

    if !utils.IsSpider(ctx.Request.UserAgent()) {
        // 更新链接的评论查看计数
        models.Link_IncViewCount(link.Id, 1)
    }

    vlink := models.Link_ToVLink([]models.Link{*link}, ctx)
    sortType := strings.ToLower(ctx.Get("cm_order")) //"hot":热门；"hotc":热议；"time":最新；"vote":得分；"ctvl":"争议"
    if sortType == "" {
        sortType = "hot"
    }
    var comments string
    if commentId > 0 {
        comments = models.GetPermalinkComment(linkId, commentId, sortType)
        ctx.ViewData["SubLinkUrl"] = fmt.Sprintf("permacoment/%d/%d/", linkId, commentId)
    } else {
        comments = models.GetSortComments("", "/", int64(0), linkId, sortType, "", false) //models.Comment_SortForLink(link.Id, "hot")
        ctx.ViewData["SubLinkUrl"] = linkId
    }

    ctx.ViewData["Comments"] = template.HTML(comments)
    ctx.ViewData["SortType"] = sortType
    ctx.ViewData["SortTypeName"] = ORDER_NAMES[sortType]

    return ctx.Render("/link/show", vlink[0])
}

//link搜索界面
func link_search(ctx *goku.HttpContext) goku.ActionResulter {
    ls := utils.LinkSearch{}
    searchResult, err := ls.SearchLink(ctx.Get("term"), 1, golink.PAGE_SIZE)
    if err == nil && searchResult.TimedOut == false && searchResult.HitResult.HitArray != nil && len(searchResult.HitResult.HitArray) > 0 {
        links, _ := models.Link_GetByIdList(searchResult.HitResult.HitArray)
        ctx.ViewData["Links"] = models.Link_ToVLink(links, ctx)
        ctx.ViewData["HasMoreLink"] = len(links) >= golink.PAGE_SIZE
    } else {
        ctx.ViewData["Links"] = nil
        ctx.ViewData["HasMoreLink"] = false
    }
    ctx.ViewData["Term"] = ctx.Get("term")

    return ctx.Render("/link/search", nil)
}

// 加载更多的搜索link
func link_search_loadMore(ctx *goku.HttpContext) goku.ActionResulter {
    term, _ := url.QueryUnescape(ctx.Get("term"))
    page, err := strconv.Atoi(ctx.Get("page"))
    success, hasmore := false, false
    errorMsgs, html := "", ""
    if err == nil && page > 1 {
        ls := utils.LinkSearch{}
        searchResult, err := ls.SearchLink(term, page, golink.PAGE_SIZE)
        if err == nil && searchResult.TimedOut == false && searchResult.HitResult.HitArray != nil {
            if len(searchResult.HitResult.HitArray) > 0 {
                links, _ := models.Link_GetByIdList(searchResult.HitResult.HitArray)
                if links != nil && len(links) > 0 {
                    ctx.ViewData["Links"] = models.Link_ToVLink(links, ctx)
                    vr := ctx.RenderPartial("loadmorelink", nil)
                    vr.Render(ctx, vr.Body)
                    html = vr.Body.String()
                    hasmore = len(links) >= golink.PAGE_SIZE
                }
            }
            success = true
        }
    } else {
        errorMsgs = "参数错误"
    }
    r := map[string]interface{}{
        "success": success,
        "errors":  errorMsgs,
        "html":    html,
        "hasmore": hasmore,
    }
    return ctx.Json(r)
}
