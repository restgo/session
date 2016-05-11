package session


type (
	Store interface {
		Init(options string) error// give options to store, so that store can connect to db, or store options
		Get(sid interface{}) (*Session, error)
		Save(session *Session) (interface{}, error) // save session to store, update expire date, nothing need to do for cookie session
		Destroy(sid interface{}) error
		StoreName() string
	}
)
