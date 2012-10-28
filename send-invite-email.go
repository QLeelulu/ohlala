package main

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/utils"
    "github.com/QLeelulu/ohlala/golink/models"
    //"strconv"
    //"strings"
	"fmt"
    "time"
)

func main() {
	
	for {
//todo: ping通网络再发送
fmt.Println("entry")
		emails, err := models.GetEmailForSend()

		if err != nil {
			goku.Logger().Errorln(err.Error())
		} else {
			for _, email := range emails {
				err := sendMail(email)
				if err != nil {
					fmt.Println(err)
					
				} else {
					fmt.Println("send", email)
				}
			}
break
			if len(emails) > 0 {
				continue
			}
		}

fmt.Println("sleep")
		time.Sleep(300 * time.Second) // 每5分钟
	}
}

func sendMail(email *models.EmailInvite) error {

	user := "zengshmin@163.com"
	password := "zeng@7839658537"
	host := "smtp.163.com:25"
	to := email.ToEmail

	subject := fmt.Sprintf("%s 邀请你加入享链", email.UserName)
	body := fmt.Sprintf(`
	<html>
	<body>
%s 邀请你加入享链
<br/>
享链(<a href="%s">%s</a>)是一个由大家共建的XXXXXXX社区，简介XXXXXXXXX。
<br/>
请点击以下链接完成注册：
<a href="http://XXX.com/register?key=%s">http://XXX.com/register?key=%s</a>
<br/>
© 享链 2012
	</body>
	</html>
	`, email.UserName, "XXX.com", "XXX.com", email.Guid, email.Guid)

	err := utils.SendMail(user, password, host, to, subject, body, "html")
	return err
}
