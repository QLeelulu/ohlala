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
    "github.com/QLeelulu/ohlala/golink/utils"
    //"io/ioutil"
    "net/http"
    "net/url"
    "strconv"
    "strings"
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
    AvatarUrl        string
    Link             string

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
    m["avatar_url"] = u.AvatarUrl
    m["link"] = u.Link
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
    m["avatar_url"] = u.AvatarUrl
    m["link"] = u.Link

    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    r, err := db.Update("third_party_user", m, "`user_id`=? AND `third_party`=?", u.UserId, u.ThirdParty)
    return r, err
}

func ThirdPartyUser_GetByThirdParty(thirdParty string, thirdPartyUserId string) (u *ThirdPartyUser) {
    u = thirdPartyUser_SearchOneBy("`third_party`=? AND `third_party_user_id`=?", thirdParty, thirdPartyUserId)
    return
}

func ThirdPartyUser_GetByUserAndThirdParty(user *User, thirdParty string) (u *ThirdPartyUser) {
    u = thirdPartyUser_SearchOneBy("`user_id`=? AND `third_party`=?", user.Id, thirdParty)
    return
}

func thirdPartyUser_SearchOneBy(criteria string, values ...interface{}) (u *ThirdPartyUser) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    sql := "SELECT `user_id`, `third_party`, `third_party_user_id`, `third_party_email`, `access_token`, `refresh_token`, `token_expire_time`, `create_time`, `last_active_time`, `avatar_url`, `link` FROM `third_party_user` WHERE " + criteria + " limit 1"
    thirdPartyUserRow, err := db.Query(sql, values...)
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
            &u.AccessToken, &u.RefreshToken, &u.TokenExpireTime, &u.CreateTime, &u.LastActiveTime, &u.AvatarUrl, &u.Link)
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
    github_provider_name = "github"
    qq_provider_name     = "qq"
)

type ThirdPartyUserProfile struct {
    Id           string `json:"Id"`
    UserName     string `json:"UserName"`
    FirstName    string `json:"FirstName"`
    LastName     string `json:"LastName"`
    Email        string `json:"Email"`
    AvatarUrl    string `json:"AvatarUrl"`
    Link         string `json:"Link"`
    ProviderName string `json:"ProviderName"`
}

func (profile *ThirdPartyUserProfile) GetDisplayName() string {
    if len(profile.UserName) > 0 {
        return profile.UserName
    }

    name := strings.Trim(profile.FirstName+" "+profile.LastName, " ")
    if len(name) > 0 {
        return name
    }

    if len(profile.Email) > 0 {
        return profile.Email
    }

    return profile.ProviderName + profile.Id
}

// third party provider, potential support protocols: oauth 1.0a, oauth 2.0, openid
type ThirdPartyProvider interface {
    Protocol() string
    ProviderName() string
    GetProfile() (*ThirdPartyUserProfile, error)

    Login(ctx *goku.HttpContext) (actionResult goku.ActionResulter, err error)
}

type ThirdPartyBindError struct {
    message string
}

func (e *ThirdPartyBindError) Error() string {
    return e.message
}

func NewThirdPartyBindError(msg string) *ThirdPartyBindError {
    return &ThirdPartyBindError{
        message: msg,
    }
}

type OAuth2Provider struct {
    Config *oauth2.Config
    Token  *oauth2.Token

    getProviderNameFunc func() string
    getUserProfileFunc  func(p *OAuth2Provider) (*ThirdPartyUserProfile, error)

    exchangeTokenFunc func(provider *OAuth2Provider, code string) (*oauth2.Token, error)
}

func (p OAuth2Provider) Protocol() string {
    return oauth2_protocol_name
}

func (p OAuth2Provider) ProviderName() string {
    return p.getProviderNameFunc()
}

func (p OAuth2Provider) GetProfile() (profile *ThirdPartyUserProfile, err error) {
    profile, err = p.getUserProfileFunc(&p)

    if err != nil {
        return
    } else if profile == nil || len(profile.Id) == 0 {
        err = errors.New("failed to get third party user profie.")
    }

    profile.ProviderName = p.ProviderName()
    return
}

func (p OAuth2Provider) Login(ctx *goku.HttpContext) (actionResult goku.ActionResulter, err error) {
    url := p.Config.AuthCodeURL("")
    actionResult = ctx.Redirect(url)
    return
}

func (p *OAuth2Provider) ExchangeToken(code string) (tok *oauth2.Token, err error) {
    fmt.Sprintf("code: %v\n", code)

    if p.exchangeTokenFunc != nil {
        tok, err = p.exchangeTokenFunc(p, code)
        return
    }

    transport := &oauth2.Transport{Config: p.Config}
    tok, err = transport.Exchange(code)
    return
}

type oauth2TokenCache struct {
    user         *User
    providerName string
}

func newOAuth2TokenCache(user *User, providerName string) (cache *oauth2TokenCache) {
    if user == nil {
        panic("user can't be nil to initilize oauth 2 token cache.")
        return
    }
    if len(providerName) == 0 {
        panic("provider can't be empty to initilize oauth 2 token cache.")
        return
    }

    cache = &oauth2TokenCache{
        user:         user,
        providerName: providerName,
    }
    return
}

func (cache *oauth2TokenCache) Token() (tok *oauth2.Token, err error) {
    u := ThirdPartyUser_GetByUserAndThirdParty(cache.user, cache.providerName)
    if u == nil {
        return
    }

    tok = &oauth2.Token{
        AccessToken:  u.AccessToken,
        RefreshToken: u.RefreshToken,
        Expiry:       u.TokenExpireTime,
    }
    return
}

func (cache *oauth2TokenCache) PutToken(tok *oauth2.Token) (err error) {
    u := ThirdPartyUser_GetByUserAndThirdParty(cache.user, cache.providerName)
    if u == nil {
        return
    }

    u.AccessToken = tok.AccessToken
    u.RefreshToken = tok.RefreshToken
    u.TokenExpireTime = tok.Expiry
    u.Update()

    return
}

// all supported providers
var thirdPartyProviderBuilders = make(map[string]func(u *User) ThirdPartyProvider)
var oauth2ProviderBuilders = make(map[string]func(u *User) *OAuth2Provider)

func ThirdParty_RegisterOAuth2Provider(providerName string, providerBuilder func(u *User) *OAuth2Provider) {
    if len(providerName) == 0 {
        panic("provider can't be empty.")
    }
    if providerBuilder == nil {
        panic("providerBuilder can't be nil.")
    }

    if _, existed := thirdPartyProviderBuilders[providerName]; existed {
        panic(fmt.Sprintf("provider %v already registered.", providerName))
    }

    oauth2ProviderBuilders[providerName] = providerBuilder
    thirdPartyProviderBuilders[providerName] = func(u *User) ThirdPartyProvider {
        return providerBuilder(u)
    }
}

const (
    google_oauth2_get_userinfo_url = "https://www.googleapis.com/oauth2/v1/userinfo"
    sina_oauth2_get_userinfo_url   = "https://api.weibo.com/2/users/show.json"
    sina_oauth2_get_uid_url        = "https://api.weibo.com/2/account/get_uid.json"
    sina_oauth2_get_email_url      = "https://api.weibo.com/2/account/profile/email.json"
    github_oauth2_get_userinfo_url = "https://api.github.com/user"
)

func googleProviderBuilder(u *User) *OAuth2Provider {
    p := &OAuth2Provider{}
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
    p.getUserProfileFunc = func(provider *OAuth2Provider) (profile *ThirdPartyUserProfile, err error) {
        if provider.Token == nil {
            panic("oauth2 token not provided yet.")
        }

        transport := &oauth2.Transport{Config: provider.Config}
        transport.Token = provider.Token
        client := transport.Client()

        r, err := client.Get(google_oauth2_get_userinfo_url)
        if err != nil {
            return
        }
        defer r.Body.Close()

        var gProfile struct {
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
        json.NewDecoder(r.Body).Decode(&gProfile)

        profile = &ThirdPartyUserProfile{
            Id:        gProfile.Id,
            FirstName: gProfile.GivenName,
            LastName:  gProfile.FamilyName,
            Email:     gProfile.Email,
            AvatarUrl: gProfile.Picture,
            Link:      gProfile.Link,
        }
        return
    }

    if u != nil {
        p.Config.TokenCache = newOAuth2TokenCache(u, google_provider_name)
    }

    return p
}

func sinaProviderBuilder(u *User) *OAuth2Provider {
    p := &OAuth2Provider{}
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
    p.getUserProfileFunc = func(provider *OAuth2Provider) (profile *ThirdPartyUserProfile, err error) {
        if provider.Token == nil {
            panic("oauth2 token not provided yet.")
        }

        client := &http.Client{}
        v := url.Values{}
        v.Add("access_token", provider.Token.AccessToken)

        getUserIdFunc := func() (id string) {
            r, err := client.Get(sina_oauth2_get_uid_url + "?" + v.Encode())
            if err != nil {
                return
            }
            defer r.Body.Close()

            var idProfile struct {
                Id int `json:"uid"`
            }
            json.NewDecoder(r.Body).Decode(&idProfile)

            if idProfile.Id > 0 {
                id = strconv.Itoa(idProfile.Id)
            }

            return
        }
        getEmailFunc := func() (email string) {
            r, err := client.Get(sina_oauth2_get_email_url + "?" + v.Encode())
            if err != nil {
                return
            }
            defer r.Body.Close()

            var emailProfile struct {
                Email string `json:"email"`
            }
            json.NewDecoder(r.Body).Decode(&emailProfile)

            email = emailProfile.Email
            return
        }

        userId, email := getUserIdFunc(), getEmailFunc()

        var userName, avatarUrl, link string
        //  get sina profile
        func() {
            if userId == "" {
                return
            }

            v.Add("uid", userId)
            r, err := client.Get(sina_oauth2_get_userinfo_url + "?" + v.Encode())
            if err != nil {
                return
            }
            defer r.Body.Close()

            var sinaProfile struct {
                UserName  string `json:"screen_name"`
                Gender    string `json:"gender"`
                AvatarUrl string `json:"profile_image_url"`
                Link      string `json:"profile_url"`
            }
            json.NewDecoder(r.Body).Decode(&sinaProfile)

            userName = sinaProfile.UserName
            avatarUrl = sinaProfile.AvatarUrl
            if len(sinaProfile.Link) > 0 {
                link = "http://weibo.com/" + sinaProfile.Link
            }
        }()

        profile = &ThirdPartyUserProfile{
            Id:        userId,
            UserName:  userName,
            FirstName: "",
            LastName:  "",
            Email:     email,
            AvatarUrl: avatarUrl,
            Link:      link,
        }
        return
    }
    // sina exchange token return json string in text/plain content, need to manually decode it here.
    p.exchangeTokenFunc = func(provider *OAuth2Provider, code string) (tok *oauth2.Token, err error) {
        if provider.Config == nil {
            return nil, errors.New("no Config supplied for exchanging token")
        }

        cfg := provider.Config
        v := url.Values{}
        v.Add("client_id", cfg.ClientId)
        v.Add("client_secret", cfg.ClientSecret)
        v.Add("grant_type", "authorization_code")
        v.Add("code", code)
        v.Add("redirect_uri", cfg.RedirectURL)
        v.Add("scope", cfg.Scope)

        response, err := http.PostForm(cfg.TokenURL, v)

        if err != nil {
            return nil, err
        }
        defer response.Body.Close()

        var tkn struct {
            AccessToken string        `json:"access_token"`
            ExpiresIn   time.Duration `json:"expires_in"`
        }

        if err = json.NewDecoder(response.Body).Decode(&tkn); err != nil {
            return
        }
        tkn.ExpiresIn *= time.Second

        tok = &oauth2.Token{
            AccessToken: tkn.AccessToken,
            Expiry:      time.Now().Add(tkn.ExpiresIn),
        }

        if cfg.TokenCache != nil {
            cfg.TokenCache.PutToken(tok)
        }

        return
    }

    if u != nil {
        p.Config.TokenCache = newOAuth2TokenCache(u, sina_provider_name)
    }

    return p
}

func githubProviderBuilder(u *User) *OAuth2Provider {
    p := &OAuth2Provider{}
    c := config.OAuth2Configs[github_provider_name]
    p.Config = &oauth2.Config{
        ClientId:     c.ClientId,
        ClientSecret: c.ClientSecret,
        Scope:        c.Scope,
        AuthURL:      c.AuthURL,
        TokenURL:     c.TokenURL,
        RedirectURL:  c.RedirectURL,
    }
    p.getProviderNameFunc = func() string {
        return github_provider_name
    }
    p.getUserProfileFunc = func(provider *OAuth2Provider) (profile *ThirdPartyUserProfile, err error) {
        if provider.Token == nil {
            panic("oauth2 token not provided yet.")
        }

        transport := &oauth2.Transport{Config: provider.Config}
        transport.Token = provider.Token
        client := transport.Client()

        r, err := client.Get(github_oauth2_get_userinfo_url)
        if err != nil {
            return
        }
        defer r.Body.Close()

        var githubProfile struct {
            Id        int    `json:"id"`
            UserName  string `json:"login"`
            Name      string `json:"name"`
            Email     string `json:"email"`
            AvatarUrl string `json:"avatar_url"`
            Link      string `json:"html_url"`
        }

        json.NewDecoder(r.Body).Decode(&githubProfile)

        profileName := strings.Replace(githubProfile.Name, ",", "", -1)
        firstName, lastName := profileName, ""
        idxSpace := strings.Index(profileName, " ")
        if idxSpace > 0 {
            firstName = profileName[:idxSpace]
            lastName = profileName[idxSpace+1:]
        }

        profile = &ThirdPartyUserProfile{
            Id:        strconv.Itoa(githubProfile.Id),
            UserName:  githubProfile.UserName,
            FirstName: firstName,
            LastName:  lastName,
            Email:     githubProfile.Email,
            AvatarUrl: githubProfile.AvatarUrl,
            Link:      githubProfile.Link,
        }
        return
    }

    if u != nil {
        p.Config.TokenCache = newOAuth2TokenCache(u, github_provider_name)
    }

    return p
}

func ThirdParty_Login(ctx *goku.HttpContext, providerName string) (actionResult goku.ActionResulter, err error) {
    providerBuilder, ok := thirdPartyProviderBuilders[providerName]
    if !ok {
        err = errors.New("invalid third party provider: " + providerName)
        return
    }

    provider := providerBuilder(nil)
    actionResult, err = provider.Login(ctx)

    return
}

func ThirdParty_OAuth2Callback(providerName, code string) (u *ThirdPartyUser, token *oauth2.Token, profile *ThirdPartyUserProfile, err error) {
    var provider *OAuth2Provider
    if builder, existed := oauth2ProviderBuilders[providerName]; existed {
        provider = builder(nil)
    } else {
        err = errors.New(fmt.Sprintf("invalid OAuth2 provider `%v`", providerName))
        return
    }

    token, err = provider.ExchangeToken(code)
    if err != nil {
        return
    }

    if token == nil || len(token.AccessToken) == 0 {
        err = errors.New("failed to get access token")
        return
    }

    provider.Token = token

    fmt.Printf("\naccess token: %v\n", token.AccessToken)

    u, profile, err = thirdParty_GetExistedThirdPartyUser(provider)

    if u != nil {
        u.AccessToken = provider.Token.AccessToken
        u.RefreshToken = provider.Token.RefreshToken
        u.TokenExpireTime = provider.Token.Expiry
        u.Update()
    }

    return
}

func thirdParty_GetExistedThirdPartyUser(provider ThirdPartyProvider) (u *ThirdPartyUser, profile *ThirdPartyUserProfile, err error) {
    profile, err = provider.GetProfile()
    thirdPartyName := provider.ProviderName()

    if err != nil {
        return
    }

    fmt.Printf("third party profile -- Id: %v, email: %v\n\n", profile.Id, profile.Email)

    u = ThirdPartyUser_GetByThirdParty(thirdPartyName, profile.Id)
    if u != nil {
        u.LastActiveTime = time.Now().UTC()
        u.Update()
        return
    }

    u, err = thirdParty_AutoBindByMatchingEmail(profile)

    return
}

func thirdParty_AutoBindByMatchingEmail(profile *ThirdPartyUserProfile) (u *ThirdPartyUser, err error) {
    if len(profile.Email) == 0 {
        return
    }

    user, err := User_GetByEmail(profile.Email)
    if user == nil || err != nil {
        return
    }

    u, err = ThirdParty_BindExistedUser(user, profile)

    return
}

func ThirdParty_BindExistedUser(user *User, profile *ThirdPartyUserProfile) (u *ThirdPartyUser, err error) {
    if user == nil {
        err = NewThirdPartyBindError("待绑定的用户不能为空")
        return
    }

    existedThirdPartyUser := ThirdPartyUser_GetByUserAndThirdParty(user, profile.ProviderName)
    if existedThirdPartyUser != nil {
        err = NewThirdPartyBindError(fmt.Sprintf("亲，你已经绑定过 %v 的另外一个账户了~~~~~", profile.ProviderName))
        return
    }

    utcNow := time.Now().UTC()
    u = &ThirdPartyUser{
        UserId:           user.Id,
        ThirdParty:       profile.ProviderName,
        ThirdPartyUserId: profile.Id,
        ThirdPartyEmail:  profile.Email,
        CreateTime:       utcNow,
        LastActiveTime:   utcNow,
        AvatarUrl:        profile.AvatarUrl,
        Link:             profile.Link,
    }
    _, err = u.Save()

    if err != nil {
        u = nil
    }

    return
}

func ThirdParty_CreateAndBind(email string, name string, profile *ThirdPartyUserProfile) (u *ThirdPartyUser, err error) {
    if len(email) == 0 {
        err = NewThirdPartyBindError("邮箱不能为空")
        return
    }
    if len(name) == 0 {
        err = NewThirdPartyBindError("昵称不能为空")
        return
    }

    if User_IsEmailExist(email) {
        err = NewThirdPartyBindError("邮箱已经被注册过了，使用登录并绑定吧~~~")
        return
    }
    if User_IsUserExist(name) {
        err = NewThirdPartyBindError("哎呀，昵称已经被占用了哟，换一个试一下吧~~~")
        return
    }

    pwd, _ := utils.GenerateRandomString(5)
    pwdHash := utils.PasswordHash(pwd)

    //TODO: send notification email
    m := make(map[string]interface{})
    m["name"] = name
    m["email"] = email
    m["pwd"] = pwdHash
    m["create_time"] = time.Now()
    _, err = User_SaveMap(m)

    if err != nil {
        return
    }

    user, err := User_GetByEmail(email)
    u, err = ThirdParty_BindExistedUser(user, profile)

    return
}

func ThirdParty_SaveThirdPartyProfileToSession(
    ctx *goku.HttpContext,
    profile *ThirdPartyUserProfile) (err error) {

    providerName := profile.ProviderName
    sessionKeyBase := thirdParty_GetSessionKeyBase(providerName, profile.Id)
    profileSessionId := ThirdParty_GetThirdPartyProfileSessionId(sessionKeyBase)
    expires := time.Now().Add(time.Duration(3600) * time.Second)

    b, _ := json.Marshal(profile)
    s := string(b)
    err = SaveItemToSession(profileSessionId, s, expires)
    if err != nil {
        return
    }

    c := &http.Cookie{
        Name:     config.ThirdPartyCookieKey,
        Value:    sessionKeyBase,
        Expires:  expires,
        Path:     "/",
        HttpOnly: true,
    }
    ctx.SetCookie(c)

    return
}

func ThirdParty_SaveOAuth2TokenToSession(
    ctx *goku.HttpContext,
    profile *ThirdPartyUserProfile,
    token *oauth2.Token) (err error) {

    providerName := profile.ProviderName
    sessionKeyBase := thirdParty_GetSessionKeyBase(providerName, profile.Id)
    oauth2TokenSessionId := ThirdParty_GetOAuthTokenSessionId(sessionKeyBase)
    expires := time.Now().Add(time.Duration(3600) * time.Second)

    b, _ := json.Marshal(token)
    s := string(b)
    err = SaveItemToSession(oauth2TokenSessionId, s, expires)
    return
}

func ThirdParty_GetThirdPartyProfileFromSession(sessinId string) (u *ThirdPartyUserProfile) {
    redisClient := GetRedis()
    defer redisClient.Quit()

    p, err := redisClient.Get(sessinId)
    if err != nil {
        return nil
    }

    jsonString := p.String()
    if jsonString == "" {
        return nil
    }

    u = new(ThirdPartyUserProfile)
    json.Unmarshal([]byte(jsonString), u)

    if len(u.Id) == 0 {
        u = nil
    }

    return
}

func ThirdParty_ManualBindUserDone(u *ThirdPartyUser, ctx *goku.HttpContext) {
    thirdParty_ClearThirdPartyProfileFromSession(ctx)

    provider := thirdPartyProviderBuilders[u.ThirdParty](u.User())
    switch provider.Protocol() {
    case oauth2_protocol_name:
        thirdParty_OAuth2BindUserDone(provider, ctx)
    }
}

func thirdParty_ClearThirdPartyProfileFromSession(ctx *goku.HttpContext) {
    sessionIdBase := ctx.Data["thirdPartySessionIdBase"].(string)
    sessinId := ThirdParty_GetThirdPartyProfileSessionId(sessionIdBase)
    RemoveItemFromSession(sessinId)

    c := &http.Cookie{
        Name:    config.ThirdPartyCookieKey,
        Expires: time.Now().Add(-10 * time.Second),
        Path:    "/",
    }
    ctx.SetCookie(c)
}

func thirdParty_OAuth2BindUserDone(p ThirdPartyProvider, ctx *goku.HttpContext) {
    sessionIdBase := ctx.Data["thirdPartySessionIdBase"].(string)
    oauth2TokenSessionId := ThirdParty_GetOAuthTokenSessionId(sessionIdBase)

    token := thirdParty_GetOAuth2TokenFromSession(oauth2TokenSessionId)
    RemoveItemFromSession(oauth2TokenSessionId)

    if token == nil {
        return
    }

    provider := p.(*OAuth2Provider)
    if provider.Config.TokenCache != nil {
        provider.Config.TokenCache.PutToken(token)
    }
}

func thirdParty_GetOAuth2TokenFromSession(sessinId string) (tok *oauth2.Token) {
    redisClient := GetRedis()
    defer redisClient.Quit()

    p, err := redisClient.Get(sessinId)
    if err != nil {
        return nil
    }

    jsonString := p.String()
    if jsonString == "" {
        return nil
    }

    tok = new(oauth2.Token)
    json.Unmarshal([]byte(jsonString), tok)

    if len(tok.AccessToken) == 0 {
        tok = nil
    }

    return
}

func thirdParty_GetSessionKeyBase(providerName string, thirdPartyUserId string) string {
    return fmt.Sprintf("third-party-%v-%v", providerName, thirdPartyUserId)
}

func ThirdParty_GetThirdPartyProfileSessionId(sessionKeyBase string) string {
    return sessionKeyBase + "-profile"
}

func ThirdParty_GetOAuthTokenSessionId(sessionKeyBase string) string {
    return sessionKeyBase + "-token"
}

func init() {
    ThirdParty_RegisterOAuth2Provider(google_provider_name, googleProviderBuilder)
    ThirdParty_RegisterOAuth2Provider(sina_provider_name, sinaProviderBuilder)
    ThirdParty_RegisterOAuth2Provider(github_provider_name, githubProviderBuilder)
}
