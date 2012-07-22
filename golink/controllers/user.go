package controllers

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink/models"
    "github.com/QLeelulu/ohlala/golink/utils"
    "time"
)

/**
 * form
 */
func createLoginForm() *form.Form {
    // defined the field
    email := form.NewEmailField("email", "Email", true).
        Error("invalid", "Email地址错误").
        Error("require", "Email地址必须填写").Field()

    pwd := form.NewCharField("pwd", "密码", true).Min(6).Max(30).
        Error("required", "密码必须填写").
        Error("range", "密码长度必须在{0}到{1}之间").Field()

    // add the fields to a form
    form := form.NewForm(email, pwd)
    return form
}

func createRegForm() *form.Form {
    email := form.NewEmailField("email", "Email", true).
        Error("invalid", "Email地址错误").
        Error("require", "Email地址必须填写").Field()

    pwd := form.NewCharField("pwd", "密码", true).Min(6).Max(30).
        Error("required", "密码必须填写").
        Error("range", "密码长度必须在{0}到{1}之间").Field()

    repwd := form.NewCharField("repwd", "确认密码", true).Min(6).Max(30).
        Error("required", "确认密码必须填写").
        Error("range", "密码长度必须在{0}到{1}之间").Field()

    // add the fields to a form
    form := form.NewForm(email, pwd, repwd)
    return form
}

var _ = goku.Controller("user").
    Get("login", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
}).
    // login
    Post("login", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
}).
    // reg
    Post("reg", func(ctx *goku.HttpContext) goku.ActionResulter {
    f := createRegForm()
    f.FillByRequest(ctx.Request)

    errorMsgs := make([]string, 0)
    if f.Valid() {
        m := f.CleanValues()
        if m["pwd"] == m["repwd"] {
            // 检查email地址是否已经注册
            emailNotOk := models.User_IsEmailExist(m["email"].(string))
            if !emailNotOk {
                m["pwd"] = utils.PasswordHash(m["pwd"].(string))
                delete(m, "repwd")
                m["create_at"] = time.Now()
                _, err := models.User_SaveMap(m)
                if err != nil {
                    errorMsgs = append(errorMsgs, "Database error")
                    goku.Logger().Errorln(err)
                }
            } else {
                errorMsgs = append(errorMsgs, "Email地址已经被注册，请换一个")
            }
        } else {
            errorMsgs = append(errorMsgs, "两次输入的密码不一样")
        }
    } else {
        errs := f.Errors()
        for _, v := range errs {
            errorMsgs = append(errorMsgs, v[0]+": "+v[1])
        }
    }

    if len(errorMsgs) < 1 {
        ctx.ViewData["success"] = true
    } else {
        ctx.ViewData["regErrors"] = errorMsgs
        ctx.ViewData["regValues"] = f.Values()
    }

    return ctx.Render("login", nil)
})
