package models

import (
    oauth2 "code.google.com/p/goauth2/oauth"
    //"crypto/tls"
    "database/sql"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/ohlala/golink/config"
    //"net/http"
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
    oauth1a_protocol_name = "OAuth1.0a"
    oauth2_protocol_name  = "OAuth2"

    google_provider_name = "google"
    sina_provider_name   = "sina"
    qq_provider_name     = "qq"
)

type thirdPartyUserProfile struct {
    Id        string
    FirstName string
    LastName  string
    Email     string
}

// thrid party provider, potential support protocols: oauth 1.0a, oauth 2.0, openid
type thirdPartyProvider interface {
    Protocol() string
    ProviderName() string
    GetProfile() (*thirdPartyUserProfile, error)

    Login(ctx *goku.HttpContext) (actionResult goku.ActionResulter, err error)
}

type oauth2Provider struct {
    Config *oauth2.Config
    Token  *oauth2.Token

    getProviderNameFunc func() string
    getUserProfileFunc  func(p *oauth2Provider) (*thirdPartyUserProfile, error)
}

func (p oauth2Provider) Protocol() string {
    return oauth2_protocol_name
}

func (p oauth2Provider) ProviderName() string {
    return p.getProviderNameFunc()
}

func (p oauth2Provider) GetProfile() (profile *thirdPartyUserProfile, err error) {
    profile, err = p.getUserProfileFunc(&p)
    return
}

func (p oauth2Provider) Login(ctx *goku.HttpContext) (actionResult goku.ActionResulter, err error) {
    url := p.Config.AuthCodeURL("")
    actionResult = ctx.Redirect(url)
    return
}

// all supported providers
var thirdPartyProviderBuilders map[string]func(u *User) thirdPartyProvider

const (
    google_oauth2_get_userinfo_url = "https://www.googleapis.com/oauth2/v1/userinfo"
    sina_oauth2_get_userinfo_url   = "https://api.weibo.com/2/users/show.json"
)

type googleProfile struct {
    Id            string `json:"id"`
    Email         string `json:"email"`
    VerifiedEmail bool   `json:"verified_email"`
    Name          string `json:"name"`
    GivenName     string `json:"given_name"`
    FamilyName    string `json:"family_name"`
    Link          string `json:"link"`
    Picture       string `json:"picture"`
    Gender        string `json:"gender"`
    Birthday      string `json:"birthday"`
    Locale        string `json:"locale"`
}

func googleProviderBuilder(u *User) *oauth2Provider {
    p := &oauth2Provider{}
    c := config.OAuth2Configs[google_provider_name]
    p.Config = &oauth2.Config{
        ClientId:     c.ClientId,
        ClientSecret: c.ClientSecret,
        Scope:        c.Scope,
        AuthURL:      c.AuthURL,
        TokenURL:     c.TokenURL,
        RedirectURL:  c.RedirectURL,
    }
    p.getProviderNameFunc = func() string {
        return google_provider_name
    }
    p.getUserProfileFunc = func(provider *oauth2Provider) (profile *thirdPartyUserProfile, err error) {
        if provider.Token == nil {
            panic("oauth2 token not provided yet.")
        }

        transport := &oauth2.Transport{Config: provider.Config}
        transport.Token = provider.Token
        client := transport.Client()

        //tlsConfig := &tls.Config{InsecureSkipVerify: true}
        //tr := &http.Transport{TLSClientConfig: tlsConfig}
        //client := &http.Client{Transport: tr}

        r, err := client.Get(google_oauth2_get_userinfo_url)
        if err != nil {
            return
        }
        defer r.Body.Close()

        gProfile := &googleProfile{}
        json.NewDecoder(r.Body).Decode(gProfile)

        profile = &thirdPartyUserProfile{
            Id:        gProfile.Id,
            FirstName: gProfile.GivenName,
            LastName:  gProfile.FamilyName,
            Email:     gProfile.Email,
        }
        return
    }

    if u != nil {
        //p.Config.TokenCache
    }

    return p
}

type sinaProfile struct {
    Id         string `json:"id"`
    ScreenName string `json:"screen_name"`
    Gender     string `json:"gender"`
}

func sinaProviderBuilder(u *User) *oauth2Provider {
    p := &oauth2Provider{}
    c := config.OAuth2Configs[sina_provider_name]
    p.Config = &oauth2.Config{
        ClientId:     c.ClientId,
        ClientSecret: c.ClientSecret,
        Scope:        c.Scope,
        AuthURL:      c.AuthURL,
        TokenURL:     c.TokenURL,
        RedirectURL:  c.RedirectURL,
    }
    p.getProviderNameFunc = func() string {
        return sina_provider_name
    }
    p.getUserProfileFunc = func(provider *oauth2Provider) (profile *thirdPartyUserProfile, err error) {
        if provider.Token == nil {
            panic("oauth2 token not provided yet.")
        }

        transport := &oauth2.Transport{Config: provider.Config}
        transport.Token = provider.Token

        r, err := transport.Client().Get(sina_oauth2_get_userinfo_url)
        if err != nil {
            return
        }
        defer r.Body.Close()

        sProfile := &sinaProfile{}
        json.NewDecoder(r.Body).Decode(sProfile)

        profile = &thirdPartyUserProfile{
            Id:        sProfile.Id,
            FirstName: "",
            LastName:  "",
            Email:     "",
        }
        return
    }

    if u != nil {
        //p.Config.TokenCache
    }

    return p
}

func ThrirdParty_Login(ctx *goku.HttpContext, providerName string) (actionResult goku.ActionResulter, err error) {
    providerBuilder, ok := thirdPartyProviderBuilders[providerName]
    if !ok {
        err = errors.New("invalid third party provider: " + providerName)
        return
    }

    provider := providerBuilder(nil)
    actionResult, err = provider.Login(ctx)

    return
}

func ThrirdParty_OAuth2Callback(providerName, code string) (u *ThirdPartyUser, token *oauth2.Token, err error) {
    var provider *oauth2Provider
    switch providerName {
    case google_provider_name:
        provider = googleProviderBuilder(nil)
    case sina_provider_name:
        provider = sinaProviderBuilder(nil)
    default:
        err = errors.New(fmt.Sprintf("invalid OAuth2 provider `%v`", providerName))
        return
    }

    fmt.Sprintf("code: %v", code)
    transport := &oauth2.Transport{Config: provider.Config}
    token, err = transport.Exchange(code)
    if err != nil {
        return
    }
    provider.Token = token

    fmt.Printf("\naccess token: %v\n", token.AccessToken)

    u, err = thridParty_GetExistedThridPartyUser(provider)

    if u != nil {
        u.AccessToken = provider.Token.AccessToken
        u.RefreshToken = provider.Token.RefreshToken
        u.TokenExpireTime = provider.Token.Expiry
        u.Update()
    }

    return
}

func thridParty_GetExistedThridPartyUser(provider thirdPartyProvider) (u *ThirdPartyUser, err error) {
    profile, err := provider.GetProfile()
    thirdPartyName := provider.ProviderName()

    if err != nil {
        return
    }

    fmt.Printf("thrid party profile- Id: %v, email: %v\n", profile.Id, profile.Email)

    u = ThirdPartyUser_GetByThirdParty(thirdPartyName, profile.Id)
    if u != nil {
        u.LastActiveTime = time.Now().UTC()
        u.Update()
        return
    }

    u, err = thridParty_AutoBindByMatchingEmail(provider, profile)

    return
}

func thridParty_AutoBindByMatchingEmail(provider thirdPartyProvider, profile *thirdPartyUserProfile) (u *ThirdPartyUser, err error) {
    if len(profile.Email) == 0 {
        return
    }

    user, err := User_GetByEmail(profile.Email)
    if err != nil {
        return
    }

    u = &ThirdPartyUser{
        UserId:           user.Id,
        ThirdParty:       provider.ProviderName(),
        ThirdPartyUserId: profile.Id,
        ThirdPartyEmail:  profile.Email,
        CreateTime:       time.Now().UTC(),
        LastActiveTime:   time.Now().UTC(),
    }
    u.Save()

    return
}

func init() {
    thirdPartyProviderBuilders = make(map[string]func(u *User) thirdPartyProvider)
    thirdPartyProviderBuilders[google_provider_name] = func(u *User) thirdPartyProvider {
        return googleProviderBuilder(u)
    }
    thirdPartyProviderBuilders[sina_provider_name] = func(u *User) thirdPartyProvider {
        return sinaProviderBuilder(u)
    }
}
