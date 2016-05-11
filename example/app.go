package main

import (
)
import (
	"github.com/fxding/grest"
	"github.com/valyala/fasthttp"
	"github.com/fxding/session"
	"flag"
	"fmt"
	"time"
)

func main() {

	router := grest.NewRouter()

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


	var addr = flag.String("addr", ":8080", "TCP address to listen to")
	fasthttp.ListenAndServe(*addr, router.FastHttpHandler)

}
