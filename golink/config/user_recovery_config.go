package config

type mailSenderConfig struct {
    SmtpServer string
    From       string
    Password   string
}

type mailContentConfig struct {
    SubjectTemplate string
    ContentTemplate string
}

type userRecoveryConfig struct {
    MailSender  *mailSenderConfig
    MailContent *mailContentConfig
}

var UserRecoveryConfig = &userRecoveryConfig{
    MailSender: &mailSenderConfig{
        SmtpServer: "smtp.163.com:25",
        From:       "xxx@163.com",
        Password:   "xxx",
    },
    MailContent: &mailContentConfig{
        SubjectTemplate: "【觅链】重置密码",
        ContentTemplate: `<html>
<body>
    <p>
        尊敬的觅链用户，
        <br />
        您好！
    </p>
    <p>
        请点击<a href="$recoveryLink" target="_blank">此处</a>重置您的密码。
    </p>
    <p>
        如果上述链接无效，请手动复制如下 URL 到浏览器地址栏：
    </p>
    <p style="padding: 10px; margin: 5px 10px; border: solid 1px #bdbdbd; font: 12px/1.2 Consolas;">
        $recoveryLink
    </p>
    <p>
        © 觅链(milnk.com) 2013
    </p>
</body>
</html>`,
    },
}
