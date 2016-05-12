package main

import (
	"fmt"
	"github.com/restgo/restgo"
	"github.com/restgo/session"
	"time"
)

func main() {

	app := restgo.App()

	sessionOpts := `{
		"Secret"     :"secret",
		"Secure"     :false,
		"Path"       :"/",
		"HttpOnly"   :true,
		"CookieName" :"cookie-session",
		"MaxAge"     : 84600,
		"EncyptCookie": false
	}`

	app.Use(session.NewSessionManager(session.NewCookieSessionStore(), sessionOpts))

	app.GET("/about", func(ctx *restgo.Context, next restgo.Next) {
		s := ctx.UserValue("session")
		session, _ := s.(*session.Session)
		if _, ok := session.Values["time"]; ok {
			fmt.Println(session.Values["time"])
		} else {
			session.Values["time"] = time.Now().Format("2006-01-02 15:04:05")
		}

		ctx.ServeText(200, "About")
	})

	app.Run()
}
