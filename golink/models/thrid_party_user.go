package models

import (
    oauth2 "code.google.com/p/goauth2/oauth"
    "database/sql"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/config"
    "time"
)

type ThirdPartyUser struct {
    UserId           int64
    ThirdParty       string
    ThirdPartyUserId string
    ThirdPartyEmail  string
    AccessToken      string
    RefreshToken     string
    TokenExpireTime  time.Time
    CreateTime       time.Time
    LastActiveTime   time.Time

    user *User `db:"exclude"`
}

func (u *ThirdPartyUser) User() *User {
    if u.user == nil {
        u.user = User_GetById(u.UserId)
    }
    return u.user
}

func (u *ThirdPartyUser) Save() (sql.Result, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    m := make(map[string]interface{})
    m["user_id"] = u.UserId
    m["third_party"] = u.ThirdParty
    m["third_party_user_id"] = u.ThirdPartyUserId
    m["third_party_email"] = u.ThirdPartyEmail
    m["access_token"] = u.AccessToken
    m["refresh_token"] = u.RefreshToken
    m["token_expire_time"] = u.TokenExpireTime
    m["create_time"] = u.CreateTime
    m["last_active_time"] = u.LastActiveTime
    r, err := db.Insert("third_party_user", m)
    return r, err
}

func (u *ThirdPartyUser) Update() (sql.Result, error) {
    m := make(map[string]interface{})
    m["third_party_email"] = u.ThirdPartyEmail
    m["access_token"] = u.AccessToken
    m["refresh_token"] = u.RefreshToken
    m["token_expire_time"] = u.TokenExpireTime
    m["create_time"] = u.CreateTime
    m["last_active_time"] = u.LastActiveTime

    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    r, err := db.Update("third_party_user", m, "`user_id`=? AND `third_party`=?", u.UserId, u.ThirdParty)
    return r, err
}

func ThirdPartyUser_GetByThirdParty(thirdParty string, thirdPartyUserId string) (u *ThirdPartyUser) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    sql := "SELECT `user_id`, `third_party`, `third_party_user_id`, `third_party_email`, `access_token`, `refresh_token`, `token_expire_time`, `create_time`, `last_active_time` FROM `third_party_user` WHERE `third_party`=? AND `third_party_user_id`=? limit 1"
    thirdPartyUserRow, err := db.Query(sql, thirdParty, thirdPartyUserId)
    if err != nil {
        return
    }
    if thirdPartyUserRow == nil {
        return
    }

    if thirdPartyUserRow.Next() {
        u = &ThirdPartyUser{}
        err = thirdPartyUserRow.Scan(
            &u.UserId, &u.ThirdParty, &u.ThirdPartyUserId, &u.ThirdPartyEmail,
            &u.AccessToken, &u.RefreshToken, &u.TokenExpireTime, &u.CreateTime, &u.LastActiveTime)
    }

    if err != nil {
        u = nil
    }

    return
}

const (
    oauth2_protocol_name = "OAuth2"

    google_provider_name = "google"
    sina_provider_name   = "sina"
    qq_provider_name     = "qq"
)

type thirdPartyProvider interface {
    Protocol() string
    ProviderName() string
    GetEmail() string
}

type oauth2Provider struct {
    Config *oauth2.Config
}

func (g *oauth2Provider) Protocol() string {
    return oauth2_protocol_name
}

type googleProvider struct {
    oauth2Provider
}

func (g *googleProvider) ProviderName() string {
    return google_provider_name
}

func (g *googleProvider) GetEmail() string {
    panic("not implemented yet.")
}

var thirdPartyProviderBuilders map[string]func() thirdPartyProvider

func googleProviderBuilder() thirdPartyProvider {
    p := &googleProvider{}
    c := config.OAuthConfigs[google_provider_name]
    p.Config = &oauth2.Config{
        ClientId:     c.ClientId,
        ClientSecret: c.ClientSecret,
        Scope:        c.Scope,
        AuthURL:      c.AuthURL,
        TokenURL:     c.TokenURL,
        RedirectURL:  c.RedirectURL,
    }
    //p.Config.TokenCache

    return p
}

func init() {
    thirdPartyProviderBuilders = make(map[string]func() thirdPartyProvider)
    thirdPartyProviderBuilders[google_provider_name] = googleProviderBuilder
}
