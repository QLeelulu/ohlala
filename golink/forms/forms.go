package forms

import (
    "github.com/QLeelulu/goku/form"
)

/**
 * 链接提交表单
 */
func CreateLinkSubmitForm() *form.Form {
    title := form.NewCharField("title", "标题", true).Min(8).Max(200).
        Error("required", "标题必须填写").
        Error("range", "标题长度必须在{0}到{1}之间").Field()

    context := form.NewRegexpField("context", "URL地址", true, `http(s)?://([\w-]+\.)+[\w-]+(/[\w- ./?%&=]*)?`).
        Error("required", "URL地址必须填写").
        Error("invalid", "URL格式不正确").Field()

    topics := form.NewCharField("topics", "话题", false).Field()

    form := form.NewForm(title, context, topics)
    return form
}

/**
 * 评论提交表单
 */
func NewCommentSubmitForm() *form.Form {
    content := form.NewTextField("content", "内容", true).Min(8).
        Error("required", "内容必须填写").
        Error("min", "内容长度必须大于{0}").Field()

    linkId := form.NewIntegerField("link_id", "链接id", true).Min(1).
        Error("required", "链接ID必须填写").
        Error("invalid", "链接ID不正确").Field()

    parentId := form.NewIntegerField("parent_id", "父评论id", false).
        Error("invalid", "父评论ID不正确").Field()

    form := form.NewForm(linkId, parentId, content)
    return form
}
