package session

import (
	"encoding/json"
	"github.com/gorilla/securecookie"
	"github.com/restgo/jsonhelper"
	"github.com/restgo/restgo"
	"github.com/valyala/fasthttp"
	"time"
)

type (
	Session struct {
		Sid    string
		store  Store
		Values map[string]interface{} // session data
	}
)

func NewSession(store Store, sid string, values map[string]interface{}) *Session {
	return &Session{
		Sid:    sid,
		store:  store,
		Values: values,
	}
}

type SessionManager struct {
	store        Store
	secureCookie *securecookie.SecureCookie
	options      jsonhelper.JsonHelper
}

// router.Use("/", NewSessionManager(newCookieStore(cookieStoreConfig), sessionManagerConfig))
// name: name for session id in cookie, default sid
func NewSessionManager(store Store, options string) restgo.HTTPHandler {
	// init store
	err := store.Init(options)
	if err != nil {
		panic("Session store init faild: " + err.Error())
	}

	manager := &SessionManager{
		store:   store,
		options: jsonhelper.NewJsonHelper([]byte(options)),
	}
	manager.initSecret()

	// use store interface to manager session
	return func(ctx *restgo.Context, next restgo.Next) {
		// 1. get session id from cookie
		sid := manager.getSidFromCookie(ctx)

		// 2. get session by id, or create one
		session, err := store.Get(sid)

		if err == nil {
			//3. set session data to requestCtx by ctx.SetUserValue
			// now you can get session anywhere. ctx.UserValue("session")
			ctx.SetUserValue("session", session)
		} else {
			//TODO
		}

		//4. save session data to store
		defer func() {
			session := ctx.UserValue("session")
			if s, ok := session.(*Session); ok {
				sid, err := store.Save(s)
				if err != nil {
					// TODO
					return
				}

				manager.setCookie(ctx, sid) // must use sid returned from store.Save()
			}
		}()

		// next handler
		next(err)
	}
}

func (this *SessionManager) initSecret() {
	// Hash keys should be at least 32 bytes long
	var hashKey = securecookie.GenerateRandomKey(32)
	// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
	// Shorter keys may weaken the encryption used.
	var blockKey = securecookie.GenerateRandomKey(32)
	this.secureCookie = securecookie.New(hashKey, blockKey)
}

func (this *SessionManager) getSidFromCookie(ctx *restgo.Context) interface{} {
	cookieName := this.options.String("CookieName", "session")
	sessionData := ctx.Request.Header.Cookie(cookieName)

	encrypted := this.options.Bool("EncyptCookie", false)
	if len(sessionData) != 0 {
		if this.store.StoreName() != "cookie" {
			var sid string = ""
			if encrypted == true {
				this.secureCookie.Decode(cookieName, string(sessionData), &sid)
			} else {
				sid = string(sessionData)
			}
			return sid
		} else {
			// for cookie store, the value of sid is session values
			var sid map[string]interface{} = make(map[string]interface{})
			if encrypted == true {
				this.secureCookie.Decode(cookieName, string(sessionData), &sid)
			} else {
				json.Unmarshal(sessionData, &sid)
			}
			return sid
		}
	}

	return ""
}

func (this *SessionManager) setCookie(ctx *restgo.Context, sid interface{}) {
	cookieName := this.options.String("CookieName", "session")
	var sessionId string
	var err error

	encrypted := this.options.Bool("EncyptCookie", false)
	if encrypted == true {
		sessionId, err = this.secureCookie.Encode(cookieName, sid)
	} else {
		if tmpId, ok := sid.(string); ok {
			sessionId = tmpId
		} else {
			var tmpId []byte
			tmpId, err = json.Marshal(sid)
			sessionId = string(tmpId)
		}
	}

	if err == nil {
		cookie := fasthttp.AcquireCookie()
		cookie.SetKey(this.options.String("CookieName", "session"))
		cookie.SetPath(this.options.String("Path", "/"))
		cookie.SetHTTPOnly(this.options.Bool("HttpOnly", true))

		exp := this.options.Int64("MaxAge", 86400) // default one day
		cookie.SetExpire(time.Now().Add(time.Duration(exp) * time.Second))

		cookie.SetValue(sessionId)

		ctx.Response.Header.SetCookie(cookie)
	}
}
