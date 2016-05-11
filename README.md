Session Middleware for [restgo](https://github.com/restgo/restgo)
[![GoDoc](https://godoc.org/github.com/restgo/session?status.svg)](https://godoc.org/github.com/restgo/session)

This package only contains cookie session store, you need implement other store if you want it.
Session store must implements `Store` interface.

## Install
```
    go get github.com/restgo/session
```

## Exampe `example/app.go`
```go
    sessionOpts := `{
        "Secret"     :"secret",
        "Secure"     :false,
        "Path"       :"/",
        "HttpOnly"   :true,
        "CookieName" :"cookie-session",
        "MaxAge"     : 84600,
        "EncyptCookie": false
    }`
    
    router.Use("/", session.NewSessionManager(session.NewCookieSessionStore(), sessionOpts))
    
    router.GET("/about", func(ctx *fasthttp.RequestCtx, next grest.Next) {
        s := ctx.UserValue("session")
        session, _ := s.(*session.Session)
        if _, ok := session.Values["time"]; ok {
            fmt.Println(session.Values["time"])
        } else {
            session.Values["time"] = time.Now().Format("2006-01-02 15:04:05")
        }
        grest.ServeTEXT(ctx, "About", 200)
    })
```