package config

import (
    "github.com/QLeelulu/ohlala/golink"
)

type thirdPartyProviderConfig struct {
    Name        string
    ImageUrl    string
    DisplayName string
    CssClass    string
}

var ThirdPartyProviderConfigs []*thirdPartyProviderConfig

type oauth2Config struct {
    ClientId     string
    ClientSecret string
    Scope        string
    AuthURL      string
    TokenURL     string
    RedirectURL  string
}

var OAuth2Configs map[string]*oauth2Config
var ThirdPartyCookieKey = "_third_party"

func init() {
    OAuth2Configs = make(map[string]*oauth2Config)

    oauthRedirectUrl := "http://" + golink.Host_Name + "/user/oauth2callback?from="

    initGoogleConfig(oauthRedirectUrl)
    initSinaConfig(oauthRedirectUrl)
    initGithubConfig(oauthRedirectUrl)
}

func initGoogleConfig(oauthRedirectUrl string) {
    googleOAuth2Config := &oauth2Config{
        ClientId:     "1098296103309.apps.googleusercontent.com",
        ClientSecret: "g707twAeUlECzD4BIy9ShEnD",
        Scope:        "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile",
        AuthURL:      "https://accounts.google.com/o/oauth2/auth",
        TokenURL:     "https://accounts.google.com/o/oauth2/token",
        RedirectURL:  oauthRedirectUrl + "google",
    }
    OAuth2Configs["google"] = googleOAuth2Config

    googleProviderConfig := &thirdPartyProviderConfig{
        Name:        "google",
        DisplayName: "Google",
        CssClass:    "google",
    }
    ThirdPartyProviderConfigs = append(ThirdPartyProviderConfigs, googleProviderConfig)
}

func initSinaConfig(oauthRedirectUrl string) {
    sinaOAuth2Config := &oauth2Config{
        ClientId:     "4257644885",
        ClientSecret: "bf7ee19929c59e363492569a17ad98fd",
        Scope:        "email",
        AuthURL:      "https://api.weibo.com/oauth2/authorize?forcelogin=false&display=default",
        TokenURL:     "https://api.weibo.com/oauth2/access_token",
        RedirectURL:  oauthRedirectUrl + "sina",
    }
    OAuth2Configs["sina"] = sinaOAuth2Config

    sinaProviderConfig := &thirdPartyProviderConfig{
        Name:        "sina",
        DisplayName: "新浪微博",
        CssClass:    "sina",
    }
    ThirdPartyProviderConfigs = append(ThirdPartyProviderConfigs, sinaProviderConfig)
}

func initGithubConfig(oauthRedirectUrl string) {
    const github_name = "github"

    githubOAuth2Config := &oauth2Config{
        ClientId:     "ca871de84f24727910a1",
        ClientSecret: "c4b9e9851c98dced1e5f33d8972cd7bae63de2bc",
        Scope:        "user:email",
        AuthURL:      "https://github.com/login/oauth/authorize",
        TokenURL:     "https://github.com/login/oauth/access_token",
        RedirectURL:  oauthRedirectUrl + github_name,
    }
    OAuth2Configs[github_name] = githubOAuth2Config

    githubProviderConfig := &thirdPartyProviderConfig{
        Name:        github_name,
        DisplayName: github_name,
        CssClass:    github_name,
    }
    ThirdPartyProviderConfigs = append(ThirdPartyProviderConfigs, githubProviderConfig)
}
