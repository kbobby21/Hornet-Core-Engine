package factory

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

var (
	UserEmail ContextKey = "user_email"
)
