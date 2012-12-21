package main

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/models"
    "github.com/QLeelulu/ohlala/golink/utils"
    //"strings"
    "fmt"
    "time"
)

func main() {

    for {
        fmt.Println("entry")
        oneSuccess := false
        emails, err := models.GetEmailForSend()

        if err != nil {
            goku.Logger().Errorln(err.Error())
        } else {
            for _, email := range emails {
                err := sendMail(email)
                if err != nil {
                    fmt.Println(err)
                    email.SendSuccess = false
                } else {
                    fmt.Println("send", email)
                    email.SendSuccess = true
                    oneSuccess = true
                }
            }
            //更新状态
            if oneSuccess == true && len(emails) > 0 {
                models.UpdateInviteEmailStatus(emails)
                continue
            }
        }

        fmt.Println("sleep")
        time.Sleep(300 * time.Second) // 每5分钟
    }
}

func sendMail(email *models.EmailInvite) error {

    user := "xxx@163.com"
    password := "xxx"
    host := "smtp.163.com:25"
    to := email.ToEmail

    subject := fmt.Sprintf("%s 邀请你加入觅链", email.UserName)
    body := fmt.Sprintf(`
	<html>
	<body>
%s 邀请你加入觅链
<br/>
觅链(<a href="http://%s">%s</a>)是一个由大家共建的XXXXXXX社区，简介XXXXXXXXX。
<br/>
请点击以下链接完成注册：
<a href="http://%s/user/reg?key=%s">http://%s/user/reg?key=%s</a>
<br/>
© 觅链 2013
	</body>
	</html>
	`, email.UserName, golink.Host_Name, golink.Host_Name, golink.Host_Name, email.Guid, golink.Host_Name, email.Guid)

    err := utils.SendMail(user, password, host, to, subject, body, "html")
    return err
}
