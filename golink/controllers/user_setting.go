package controllers

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/filters"
    "github.com/QLeelulu/ohlala/golink/models"
    "github.com/QLeelulu/ohlala/golink/utils"
    "strings"
)

// 更新基本信息
func createBaseInfoForm() *form.Form {
    description := form.NewTextField("description", "自我介绍", false).Max(100).
        Error("max-length", "自我介绍的字数不能多于{0}个").Field()

    name := form.NewCharField("name", "用户名", true).Min(2).Max(15).
        Error("required", "用户名必须填写").
        Error("range", "用户名长度必须在{0}到{1}之间").Field()

    // add the fields to a form
    form := form.NewForm(description, name)
    return form
}

// 更新密码
func createUpdatePwdForm() *form.Form {
    oldPwd := form.NewCharField("old-pwd", "原密码", true).Min(6).Max(30).
        Error("required", "原密码必须填写").
        Error("range", "原密码长度必须在{0}到{1}之间").Field()

    newPwd := form.NewCharField("new-pwd", "新密码", true).Min(6).Max(30).
        Error("required", "新密码必须填写").
        Error("range", "新密码长度必须在{0}到{1}之间").Field()

    newPwd2 := form.NewCharField("new-pwd2", "确认密码", true).Min(6).Max(30).
        Error("required", "确认密码必须填写").
        Error("range", "确认密码长度必须在{0}到{1}之间").Field()

    // add the fields to a form
    form := form.NewForm(oldPwd, newPwd, newPwd2)
    return form
}

/**
 * Controller: user
 */
var _ = goku.Controller("user").

    /**
     * 查看用户设置页
     */
    Get("setting", func(ctx *goku.HttpContext) goku.ActionResulter {

    user := ctx.Data["user"].(*models.User)
    return ctx.View(models.User_ToVUser(user, ctx))

}).Filters(filters.NewRequireLoginFilter()).

    /**
     * 更新用户基本信息
     */
    Post("update-base", func(ctx *goku.HttpContext) goku.ActionResulter {

    user := ctx.Data["user"].(*models.User)
    f := createBaseInfoForm()
    f.FillByRequest(ctx.Request)

    errorMsgs := make([]string, 0)
    if f.Valid() {
        m := f.CleanValues()
        // 检查用户名是否已经注册
        userExist := false
        if strings.ToLower(m["name"].(string)) != strings.ToLower(user.Name) {
            userExist = models.User_IsUserExist(m["name"].(string))
        }
        if userExist {
            errorMsgs = append(errorMsgs, "用户名已经被注册，请换一个")
        } else {
            _, err := models.User_Update(user.Id, m)
            if err != nil {
                errorMsgs = append(errorMsgs, golink.ERROR_DATABASE)
                goku.Logger().Errorln(err)
            } else {
                user = models.User_GetById(user.Id)
            }
        }
    } else {
        errs := f.Errors()
        for _, v := range errs {
            errorMsgs = append(errorMsgs, v[0]+": "+v[1])
        }
    }

    if len(errorMsgs) < 1 {
        ctx.ViewData["updateBaseSuccess"] = true
    } else {
        ctx.ViewData["updateBaseErrors"] = errorMsgs
        v := f.Values()
        user.Name = v["name"]
        user.Description = v["description"]
    }

    return ctx.Render("setting", user)

}).Filters(filters.NewRequireLoginFilter()).

    /**
     * 更新用户密码
     */
    Post("change-pwd", func(ctx *goku.HttpContext) goku.ActionResulter {

    user := ctx.Data["user"].(*models.User)
    f := createUpdatePwdForm()
    f.FillByRequest(ctx.Request)

    errorMsgs := make([]string, 0)
    if f.Valid() {
        m := f.CleanValues()
        if m["new-pwd"] == m["new-pwd"] {
            // 检查原密码是否正确
            if utils.PasswordHash(m["old-pwd"].(string)) == user.Pwd {
                saveMap := map[string]interface{}{"pwd": utils.PasswordHash(m["new-pwd"].(string))}
                _, err := models.User_Update(user.Id, saveMap)
                if err != nil {
                    errorMsgs = append(errorMsgs, golink.ERROR_DATABASE)
                    goku.Logger().Errorln(err)
                }
            } else {
                errorMsgs = append(errorMsgs, "原密码不正确，请重新输入")
            }
        } else {
            errorMsgs = append(errorMsgs, "两次输入的新密码不一致")
        }
    } else {
        errs := f.Errors()
        for _, v := range errs {
            errorMsgs = append(errorMsgs, v[0]+": "+v[1])
        }
    }

    if len(errorMsgs) < 1 {
        ctx.ViewData["updatePwdSuccess"] = true
    } else {
        ctx.ViewData["updatePwdErrors"] = errorMsgs
    }

    return ctx.Render("setting", user)

}).Filters(filters.NewRequireLoginFilter())
