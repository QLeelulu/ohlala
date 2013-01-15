package controllers

import (
    "bytes"
    "fmt"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/forms"
    "github.com/QLeelulu/ohlala/golink/models"
    "github.com/QLeelulu/ohlala/golink/utils"
    "html/template"
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

}).Filters(filters.NewRequireLoginFilter())

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
            cn.VoteUp = 1
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

    success, linkId, errorMsgs := models.Link_SaveForm(f, (ctx.Data["user"].(*models.User)).Id)

    if success {
        return ctx.Redirect(fmt.Sprintf("/link/%d", linkId))
    } else {
        ctx.ViewData["Errors"] = errorMsgs
        ctx.ViewData["Values"] = f.Values()
    }
    return ctx.View(nil)

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
