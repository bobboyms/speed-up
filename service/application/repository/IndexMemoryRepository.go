package repositories

type DataRepository interface {
	Get(key string) string
	Set(key, value string)
	Gets(key ...string) []string
}
