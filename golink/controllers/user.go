package controllers

import (
    "crypto/md5"
    "fmt"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/goku/form"
    "github.com/QLeelulu/ohlala/golink/models"
    "github.com/QLeelulu/ohlala/golink/utils"
    "net/http"
    "strings"
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

    remeber_me := form.NewCharField("remeber_me", "记住我", false)

    // add the fields to a form
    form := form.NewForm(email, pwd, remeber_me)
    return form
}

func createRegForm() *form.Form {
    email := form.NewEmailField("email", "Email", true).
        Error("invalid", "Email地址错误").
        Error("require", "Email地址必须填写").Field()

    pwd := form.NewCharField("pwd", "密码", true).Min(6).Max(30).
        Error("required", "密码必须填写").
        Error("range", "密码长度必须在{0}到{1}之间").Field()

    name := form.NewCharField("name", "昵称", true).Min(2).Max(15).
        Error("required", "昵称必须填写").
        Error("range", "昵称长度必须在{0}到{1}之间").Field()

    repwd := form.NewCharField("repwd", "确认密码", true).Min(6).Max(30).
        Error("required", "确认密码必须填写").
        Error("range", "密码长度必须在{0}到{1}之间").Field()

    // add the fields to a form
    form := form.NewForm(email, name, pwd, repwd)
    return form
}

var _ = goku.Controller("user").
    // login view
    Get("login", func(ctx *goku.HttpContext) goku.ActionResulter {
    if u, ok := ctx.Data["user"]; ok && u != nil {
        return ctx.Redirect("/")
    }
    return ctx.View(nil)
}).
    // reg view
    Get("reg", func(ctx *goku.HttpContext) goku.ActionResulter {
    if u, ok := ctx.Data["user"]; ok && u != nil {
        return ctx.Redirect("/")
    }
    return ctx.Render("login", nil)
}).
    // logout
    Get("logout", func(ctx *goku.HttpContext) goku.ActionResulter {

    redisClient := models.GetRedis()
    defer redisClient.Quit()
    redisClient.Del("_glut")
    c := &http.Cookie{
        Name:    "_glut",
        Expires: time.Now().Add(-10 * time.Second),
        Path:    "/",
    }
    ctx.SetCookie(c)
    return ctx.Redirect("/")
}).
    // login
    Post("login", func(ctx *goku.HttpContext) goku.ActionResulter {

    f := createLoginForm()
    f.FillByRequest(ctx.Request)

    errorMsgs := make([]string, 0)
    if f.Valid() {
        m := f.CleanValues()
        email := strings.ToLower(m["email"].(string))
        // 检查密码是否正确
        userId := models.User_CheckPwd(email, m["pwd"].(string))
        if userId > 0 {
            now := time.Now()
            h := md5.New()
            h.Write([]byte(fmt.Sprintf("%v-%v", email, now.Unix())))
            ticket := fmt.Sprintf("%x_%v", h.Sum(nil), now.Unix())
            var expires time.Time
            if m["remeber_me"] == "1" {
                expires = now.Add(365 * 24 * time.Hour)
            } else {
                expires = now.Add(48 * time.Hour)
            }
            redisClient := models.GetRedis()
            defer redisClient.Quit()
            err := redisClient.Set(ticket, userId)
            if err != nil {
                goku.Logger().Errorln(err.Error())
                errorMsgs = append(errorMsgs, "登陆票据服务器出错，请重试")
            } else {
                _, err = redisClient.Expireat(ticket, expires.Unix())
                if err != nil {
                    goku.Logger().Errorln(err.Error())
                    errorMsgs = append(errorMsgs, "登陆票据服务器出错，请重试")
                }
                c := &http.Cookie{
                    Name:    "_glut",
                    Value:   ticket,
                    Expires: expires,
                    //Domain:     ".godev.local", // edit (or omit)
                    Path:     "/", // ^ ditto
                    HttpOnly: true,
                }
                ctx.SetCookie(c)
            }
        } else {
            errorMsgs = append(errorMsgs, "账号密码不正确，请改正")
        }
    } else {
        errs := f.Errors()
        for _, v := range errs {
            errorMsgs = append(errorMsgs, v[0]+": "+v[1])
        }
    }

    if len(errorMsgs) < 1 {
        // ctx.ViewData["loginSuccess"] = true
        returnUrl := ctx.Request.FormValue("rurl")
        if returnUrl == "" {
            returnUrl = "/"
        }
        return ctx.Redirect(returnUrl)
    } else {
        ctx.ViewData["loginErrors"] = errorMsgs
        ctx.ViewData["loginValues"] = f.Values()
    }

    return ctx.Render("login", nil)
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
        ctx.ViewData["regSuccess"] = true
    } else {
        ctx.ViewData["regErrors"] = errorMsgs
        ctx.ViewData["regValues"] = f.Values()
    }

    return ctx.Render("login", nil)
})
