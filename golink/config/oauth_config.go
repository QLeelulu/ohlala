package config

import (
    "github.com/QLeelulu/ohlala/golink"
)

type thridPartyProviderConfig struct {
    Name        string
    ImageUrl    string
    DisplayName string
}

var ThridPartyProviderConfigs []*thridPartyProviderConfig

type oauth2Config struct {
    ClientId     string
    ClientSecret string
    Scope        string
    AuthURL      string
    TokenURL     string
    RedirectURL  string
}

var OAuth2Configs map[string]*oauth2Config

func init() {
    OAuth2Configs = make(map[string]*oauth2Config)

    oauthRedirectUrl := "http://" + golink.Host_Name + "/user/oauth2callback?from="

    googleOAuth2Config := &oauth2Config{
        ClientId:     "1098296103309.apps.googleusercontent.com",
        ClientSecret: "g707twAeUlECzD4BIy9ShEnD",
        Scope:        "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile",
        AuthURL:      "https://accounts.google.com/o/oauth2/auth",
        TokenURL:     "https://accounts.google.com/o/oauth2/token",
        RedirectURL:  oauthRedirectUrl + "google",
    }
    OAuth2Configs["google"] = googleOAuth2Config

    googleProviderConfig := &thridPartyProviderConfig{
        Name:        "google",
        DisplayName: "使用Google账户登录",
    }
    ThridPartyProviderConfigs = append(ThridPartyProviderConfigs, googleProviderConfig)

    sinaOAuth2Config := &oauth2Config{
        ClientId:     "4257644885",
        ClientSecret: "bf7ee19929c59e363492569a17ad98fd",
        Scope:        "",
        AuthURL:      "https://api.weibo.com/oauth2/authorize?forcelogin=false&display=default",
        TokenURL:     "https://api.weibo.com/oauth2/access_token",
        RedirectURL:  oauthRedirectUrl + "sina",
    }
    OAuth2Configs["sina"] = sinaOAuth2Config

    sinaProviderConfig := &thridPartyProviderConfig{
        Name:        "sina",
        DisplayName: "使用新浪微博账户登录",
    }
    ThridPartyProviderConfigs = append(ThridPartyProviderConfigs, sinaProviderConfig)
}
