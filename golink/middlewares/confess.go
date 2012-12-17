package middlewares

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/models"
    "math/rand"
)

var (
    luid int64   = 10011
    uids []int64 = []int64{10011, 10012, 10013}
)

// 前期推广不容易啊，
// 随机用户吧，别老是同一个用户,
// 看着不好看
type ConfessMiddleware struct {
}

func (m *ConfessMiddleware) OnBeginRequest(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    return nil, nil
}

func (m *ConfessMiddleware) OnBeginMvcHandle(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    if iuser, ok := ctx.Data["user"]; ok && iuser != nil {
        user := iuser.(*models.User)
        if user.Id == luid {
            n := rand.Intn(3)
            uid := uids[n]
            if uid != user.Id {
                user := models.User_GetById(uid)
                if user != nil {
                    ctx.Data["user"] = user
                    ctx.ViewData["user"] = user
                }
            }
        }
    }
    return nil, nil
}
func (m *ConfessMiddleware) OnEndMvcHandle(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    return nil, nil
}

func (m *ConfessMiddleware) OnEndRequest(ctx *goku.HttpContext) (goku.ActionResulter, error) {
    return nil, nil
}
