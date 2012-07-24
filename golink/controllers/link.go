package controllers

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/models"
    "strings"
)

/**
 * form
 */
func createLinkSubmitForm() *form.Form {
    title := form.NewCharField("title", "标题", true).Min(8).Max(140).
        Error("require", "标题必须填写").
        Error("range", "标题长度必须在{0}到{1}之间").Field()

    context := form.NewRegexpField("context", "URL地址", true, `http(s)?://([\w-]+\.)+[\w-]+(/[\w- ./?%&=]*)?`).
        Error("required", "URL地址必须填写").
        Error("invalid", "URL格式不正确").Field()

    tags := form.NewCharField("tags", "TAG标签", false).Field()

    form := form.NewForm(title, context, tags)
    return form
}

var _ = goku.Controller("link").
    // 
    Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {

    return ctx.View(nil)
}).
    /**
     * 查看一个链接的评论
     */
    Get("show", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
}).

    /**
     * 提交链接的表单
     */
    Get("submit", func(ctx *goku.HttpContext) goku.ActionResulter {

    return ctx.View(nil)

}).Filters(requireLoginFilter).

    /**
     * 提交一个链接并保存到数据库
     */
    Post("submit", func(ctx *goku.HttpContext) goku.ActionResulter {

    f := createLinkSubmitForm()
    f.FillByRequest(ctx.Request)

    errorMsgs := make([]string, 0)
    if f.Valid() {
        m := f.CleanValues()
        m["tags"] = buildTags(m["tags"].(string))
        m["user_id"] = (ctx.Data["user"].(*models.User)).Id
        id := models.Link_SaveMap(m)
        if id > 0 {
            models.Tag_SaveTags(m["tags"].(string), id)
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
        return ctx.Redirect("/")
    } else {
        ctx.ViewData["Errors"] = errorMsgs
        ctx.ViewData["Values"] = f.Values()
    }
    return ctx.View(nil)

}).Filters(requireLoginFilter).

    /**
     * 添加评论
     */
    Post("comment", func(ctx *goku.HttpContext) goku.ActionResulter {

    return ctx.View(nil)

}).Filters(requireLoginFilter)

// tag可以用英文逗号或者空格分隔
// 过滤重复tag，最终返回的tag列表只用英文逗号分隔
func buildTags(tags string) string {
    if tags == "" {
        return ""
    }
    m := make(map[string]string)
    t := strings.Split(tags, ",")
    for _, tag := range t {
        tag = strings.TrimSpace(tag)
        if tag != "" {
            t2 := strings.Split(tag, " ")
            for _, tag2 := range t2 {
                tag2 = strings.TrimSpace(tag2)
                if tag2 != "" {
                    m[strings.ToLower(tag2)] = tag2
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
