package cache

type AuthCache interface {
	GetValue(key string) (string, error)
	SetValue(key string, value string, ex int) error
}