package forms

import (
    "github.com/QLeelulu/goku/form"
)

func CreateLinkSubmitForm() *form.Form {
    title := form.NewCharField("title", "标题", true).Min(8).Max(140).
        Error("require", "标题必须填写").
        Error("range", "标题长度必须在{0}到{1}之间").Field()

    context := form.NewRegexpField("context", "URL地址", true, `http(s)?://([\w-]+\.)+[\w-]+(/[\w- ./?%&=]*)?`).
        Error("required", "URL地址必须填写").
        Error("invalid", "URL格式不正确").Field()

    topics := form.NewCharField("topics", "话题", false).Field()

    form := form.NewForm(title, context, topics)
    return form
}
