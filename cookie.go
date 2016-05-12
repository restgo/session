package session

type CookieStore struct{}

func NewCookieSessionStore() *CookieStore {
	return &CookieStore{}
}

/// implement Store interface ///

// do nothing
func (this *CookieStore) Init(options string) error {
	return nil
}

// store name, "cookie"
func (this *CookieStore) StoreName() string {
	return "cookie"
}

// for cookie store, sid will be session value
func (this *CookieStore) Get(sid interface{}) (*Session, error) {
	session := NewSession(this, "", make(map[string]interface{})) // session id is empty for cookie store

	if values, ok := sid.(map[string]interface{}); ok {
		session.Values = values
	}

	return session, nil
}

// Save Session, do nothing, return sessions as sid
func (this *CookieStore) Save(session *Session) (interface{}, error) {

	return session.Values, nil
}

// Destroy session by id
func (this *CookieStore) Destroy(sid interface{}) error {
	return nil
}
