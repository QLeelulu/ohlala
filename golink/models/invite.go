package models

import (
    //"fmt"
	"errors"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink"
    "github.com/QLeelulu/ohlala/golink/utils"
    "time"
    "strings"
)

type RegisterInvite struct {
	Guid           string
	UserId        int64
	ToEmail       string
	IsRegister    bool
	ExpiredDate   time.Time
	IsSend        bool
	FailCount     int
}
type EmailInvite struct {
	Guid         string
	ToEmail      string
	UserName     string
	SendSuccess  bool
}
//往邀请表添加多个邀请
func CreateRegisterInvite(userId int64, toEmails string) (bool, error) {
	var arrEmails []string
	var db *goku.MysqlDB = GetDB()
//db.Debug = true
	defer db.Close()
    if toEmails != "" {
        arrEmails = strings.Split(toEmails, ";")
		for _, email := range arrEmails {
            _, err := saveRegisterInvite(userId, email, db)
			if err != nil {
				return false, err
			}
        }
	}
	
	return true, nil
}
//保存一个邀请
func saveRegisterInvite(userId int64, toEmail string, db *goku.MysqlDB) (string, error) {
	if userId <= int64(0) || toEmail == "" {
        return "", errors.New("用户id不合法")
    }
	
	invite := new(RegisterInvite)
	invite.Guid = utils.GeneticKey()
	invite.UserId = userId
	invite.ToEmail = toEmail
	invite.IsRegister = false
	invite.ExpiredDate = time.Now().AddDate(0, 0, golink.Register_Invite_Expired_Day)
	invite.IsSend = false
	invite.FailCount = 0
	_, err := db.InsertStruct(invite)

	return invite.Guid, err
}

//获取一个邀请码(非邮件邀请方式,可能是通过qq和微薄发送)
func CreateRegisterInviteWithoutEmail(userId int64) (string, error) {
	if userId <= int64(0) {
        return "", errors.New("用户id不合法")
    }

	var db *goku.MysqlDB = GetDB()
	defer db.Close()

	invite := new(RegisterInvite)
	invite.Guid = utils.GeneticKey()
	invite.UserId = userId
	invite.ToEmail = ""
	invite.IsRegister = false
	invite.ExpiredDate = time.Now().AddDate(0, 0, golink.Register_Invite_Expired_Day)
	invite.IsSend = true
	invite.FailCount = 0
	_, err := db.InsertStruct(invite)

	return invite.Guid, err
}

//剩余的邀请码数量
func RegisterInviteRemainCount(userId int64) int {
	iCount := 0
	if userId <= int64(0) {
        return iCount
    }

	var db *goku.MysqlDB = GetDB()
	defer db.Close()
	
	sql := "SELECT COUNT(1) AS FCount FROM register_invite WHERE user_id=?"
	rows, err := db.Query(sql, userId)
	if err == nil && rows.Next() {
		rows.Scan(&iCount)
		iCount = golink.Register_Invite_Count_Max - iCount
		if iCount < 0 {
			iCount = 0
		}
	}

	return iCount
}

//验证邀请码
func VerifyInviteKey(key string) *RegisterInvite {
	if len(key) != golink.Genetic_Key_Len {
		return nil
	}

    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    ri := new(RegisterInvite)
    err := db.GetStruct(ri, "`Guid`=? AND `expired_date`>=? AND `is_register`=0", key, time.Now())
    if err != nil {
        goku.Logger().Errorln(err.Error())
        return nil
    }

	if len(ri.Guid) > 0 {
		return ri
	}
	return nil
}

//更新邀请码
func UpdateIsRegister(invite *RegisterInvite) {

    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    db.Query("UPDATE `register_invite` SET `is_register`=1 WHERE `Guid`=?", invite.Guid)
	
}

//获取需要发送的邀请email
func GetEmailForSend() ([]*EmailInvite, error) {

	var db *goku.MysqlDB = GetDB()
	defer db.Close()
	emails := make([]*EmailInvite, 0)

	strSQL := "SELECT RI.guid,RI.to_email,U.name FROM `register_invite` RI INNER JOIN `user` U ON RI.user_id=U.id AND RI.is_register=0 AND RI.is_send=0 AND RI.fail_count<? LIMIT 0,100"
	rows, err := db.Query(strSQL, golink.Register_Invite_Fail_Count_Max)
	if err == nil {
		for rows.Next() {
			email := &EmailInvite{}
			rows.Scan(&email.Guid, &email.ToEmail, &email.UserName)
			email.SendSuccess = false
		    emails = append(emails, email)
		}
	}

	return emails, err
}

//更新发送的邀请email状态
func UpdateInviteEmailStatus(emails []*EmailInvite) {

	var db *goku.MysqlDB = GetDB()
	defer db.Close()

	for _, email := range emails {
		if email.SendSuccess == true {
			db.Query("UPDATE register_invite SET is_send=1 WHERE guid=?", email.Guid)
		} else {
			db.Query("UPDATE register_invite SET fail_count=fail_count+1 WHERE guid=?", email.Guid)
		}
	}
}



